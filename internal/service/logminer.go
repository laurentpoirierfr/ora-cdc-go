package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/laurentpoirierfr/ora-cdc-go/internal/model"

	_ "github.com/sijms/go-ora/v2" // Driver Oracle
)

// LogMinerWorker gère la logique de polling pour LogMiner
type LogMinerWorker struct {
	db               *sql.DB
	config           model.LogMinerConfig
	stopChan         chan struct{}
	running          bool
	mu               sync.Mutex
	wg               sync.WaitGroup
	lastProcessedSCN int64 // SCN pour savoir où reprendre
}

// NewLogMinerWorker crée une nouvelle instance de LogMinerWorker
func NewLogMinerWorker(config model.LogMinerConfig) (*LogMinerWorker, error) {
	db, err := sql.Open("oracle", config.DBConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return &LogMinerWorker{
		db:       db,
		config:   config,
		stopChan: make(chan struct{}),
	}, nil
}

// Start démarre la goroutine de LogMiner
func (lm *LogMinerWorker) Start() {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	if lm.running {
		log.Println("LogMinerWorker is already running")
		return
	}

	lm.running = true
	lm.wg.Add(1)
	go lm.run()
	log.Println("LogMinerWorker started")
}

// Stop arrête la goroutine de LogMiner
func (lm *LogMinerWorker) Stop() {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	if !lm.running {
		log.Println("LogMinerWorker is not running")
		return
	}

	close(lm.stopChan)
	lm.wg.Wait()
	lm.stopChan = make(chan struct{})
	lm.running = false
	log.Println("LogMinerWorker stopped")
}

// Restart redémarre la goroutine de LogMiner
func (lm *LogMinerWorker) Restart() {
	lm.Stop()
	lm.Start()
}

// Close ferme la connexion à la base de données
func (lm *LogMinerWorker) Close() error {
	return lm.db.Close()
}

// run contient la logique principale de la goroutine
func (lm *LogMinerWorker) run() {
	defer lm.wg.Done()

	ticker := time.NewTicker(lm.config.PollFrequency)
	defer ticker.Stop()

	for {
		select {
		case <-lm.stopChan:
			return
		case <-ticker.C:
			lm.pollLogMiner()
		}
	}
}

// pollLogMiner interroge LogMiner pour obtenir les événements
func (lm *LogMinerWorker) pollLogMiner() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
		SELECT 
			SCN, TIMESTAMP, OPERATION, SEG_OWNER, TABLE_NAME, 
			SQL_REDO, SQL_UNDO, ROW_ID, USERNAME, SESSION#, ROLLBACK
		FROM 
			V$LOGMNR_CONTENTS
		WHERE 
			OPERATION IN ('INSERT', 'UPDATE', 'DELETE') AND SCN > :1
		ORDER BY SCN
	`

	rows, err := lm.db.QueryContext(ctx, query, lm.lastProcessedSCN)
	if err != nil {
		log.Printf("Failed to query LogMiner: %v\n", err)
		return
	}
	defer rows.Close()

	var maxSCN int64

	for rows.Next() {
		var row model.LogMinerRow
		var timestamp sql.NullTime
		if err := rows.Scan(
			&row.SCN,
			&timestamp,
			&row.Operation,
			&row.SegOwner,
			&row.TableName,
			&row.SQLRedo,
			&row.SQLUndo,
			&row.RowID,
			&row.Username,
			&row.Session,
			&row.Rollback,
		); err != nil {
			log.Printf("Failed to scan row: %v\n", err)
			continue
		}

		// Gérer les valeurs NULL pour le champ TIMESTAMP
		if timestamp.Valid {
			row.Timestamp = timestamp.Time
		}

		// Appelle la fonction callback pour chaque ligne
		lm.config.Callback(row)

		// Met à jour le SCN maximum traité
		if row.SCN > maxSCN {
			maxSCN = row.SCN
		}
	}

	// Si un SCN maximum a été traité, sauvegarder
	if maxSCN > lm.lastProcessedSCN {
		if err := lm.saveLastProcessedSCN(maxSCN); err != nil {
			log.Printf("Failed to save last processed SCN: %v\n", err)
		}
		lm.lastProcessedSCN = maxSCN
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating through rows: %v\n", err)
	}
}

// Start initialise le LogMinerWorker
func (lm *LogMinerWorker) LoadLastProcessedSCN() error {
	query := "SELECT LAST_SCN FROM LOGMINER_STATE ORDER BY LAST_TIMESTAMP DESC FETCH FIRST 1 ROW ONLY"
	err := lm.db.QueryRow(query).Scan(&lm.lastProcessedSCN)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to load last processed SCN: %w", err)
	}

	// Si aucune ligne n'existe, démarrer à partir de SCN 0
	if err == sql.ErrNoRows {
		lm.lastProcessedSCN = 0
	}
	return nil
}

// saveLastProcessedSCN sauvegarde le dernier SCN traité
func (lm *LogMinerWorker) saveLastProcessedSCN(scn int64) error {
	query := "INSERT INTO LOGMINER_STATE (LAST_SCN) VALUES (:1)"
	_, err := lm.db.Exec(query, scn)
	if err != nil {
		return fmt.Errorf("failed to save last processed SCN: %w", err)
	}
	return nil
}

// InitLogMinerStateTable vérifie l'existence de la table LOGMINER_STATE et la crée si elle n'existe pas
// InitLogMinerStateTable vérifie l'existence de la table LOGMINER_STATE et la crée si nécessaire
func (w *LogMinerWorker) InitLogMinerStateTable() error {
	// Vérifie si la table LOGMINER_STATE existe
	query := `
	    SELECT COUNT(*)
	    FROM all_tables
	    WHERE owner = UPPER(USER) AND table_name = 'LOGMINER_STATE'
	`
	var count int
	err := w.db.QueryRow(query).Scan(&count)
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification de l'existence de la table LOGMINER_STATE : %v", err)
	}

	// Si la table n'existe pas, la créer
	if count == 0 {
		createTableQuery := `
            CREATE TABLE LOGMINER_STATE (
                ID              NUMBER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
                LAST_SCN        NUMBER NOT NULL,
                LAST_TIMESTAMP  TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
            )
        `
		_, err = w.db.Exec(createTableQuery)
		if err != nil {
			return fmt.Errorf("erreur lors de la création de la table LOGMINER_STATE : %v", err)
		}
		log.Println("Table LOGMINER_STATE créée avec succès.")
	} else {
		log.Println("La table LOGMINER_STATE existe déjà.")
	}

	return nil
}

// PrepareLogMiner prépare LogMiner dans Oracle (ajoute les redo logs et démarre LogMiner)
func (w *LogMinerWorker) PrepareLogMiner() error {
	// Étape 1 : Ajouter les fichiers redo logs
	addRedoLogsQuery := `
        BEGIN
            FOR redo_file IN (
                SELECT member FROM v$logfile
            ) LOOP
                DBMS_LOGMNR.ADD_LOGFILE(
                    logfilename => redo_file.member,
                    options      => DBMS_LOGMNR.NEW
                );
            END LOOP;
        END;
    `
	_, err := w.db.Exec(addRedoLogsQuery)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ajout des redo logs : %v", err)
	}
	log.Println("Fichiers redo logs ajoutés à LogMiner.")

	// Étape 2 : Démarrer LogMiner
	startLogMinerQuery := `
        BEGIN
            DBMS_LOGMNR.START_LOGMNR( options => DBMS_LOGMNR.DICT_FROM_ONLINE_CATALOG );
        END;
    `
	_, err = w.db.Exec(startLogMinerQuery)
	if err != nil {
		return fmt.Errorf("erreur lors du démarrage de LogMiner : %v", err)
	}
	log.Println("LogMiner démarré avec succès.")

	return nil
}
