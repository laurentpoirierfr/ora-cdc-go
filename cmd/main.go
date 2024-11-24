package main

import (
	"fmt"
	"log"
	"time"

	"github.com/laurentpoirierfr/ora-cdc-go/internal/models"
	"github.com/laurentpoirierfr/ora-cdc-go/internal/worker"
)

// Exemple de fonction callback
func handleLogMinerRow(row models.LogMinerRow) {
	fmt.Printf("LogMiner Event:\n"+
		"  SCN: %d\n"+
		"  Timestamp: %s\n"+
		"  Operation: %s\n"+
		"  Schema: %s\n"+
		"  Table: %s\n"+
		"  SQL_REDO: %s\n"+
		"  SQL_UNDO: %s\n"+
		"  RowID: %s\n"+
		"  Username: %s\n"+
		"  Session: %d\n"+
		"  Rollback: %s\n\n",
		row.SCN, row.Timestamp, row.Operation, row.SegOwner, row.TableName,
		row.SQLRedo, row.SQLUndo, row.RowID, row.Username, row.Session, row.Rollback)
}

type Account struct {
	User     string
	Password string
}

func main() {

	accounts := []Account{
		{
			User:     "system",
			Password: "password",
		},
		{
			User:     "demo",
			Password: "demo",
		},
		{
			User:     "cdc_user",
			Password: "password",
		},
	}

	dbServices := []string{"FREE", "FREEPDB1"}

	dsn := fmt.Sprintf("oracle://%s:%s@%s:%s/%s", accounts[2].User, accounts[2].Password, "localhost", "1521", dbServices[0])

	config := models.LogMinerConfig{
		DBConnectionString: dsn,
		PollFrequency:      5 * time.Second,
		Callback:           handleLogMinerRow,
	}

	worker, err := worker.NewLogMinerWorker(config)
	if err != nil {
		log.Fatalf("Failed to create LogMinerWorker: %v", err)
	}
	defer worker.Close()

	// Préparation de LogMiner
	err = worker.PrepareLogMiner()
	if err != nil {
		log.Fatalf("Erreur lors de la préparation de LogMiner : %v", err)
	}

	// Initialisation de la table LOGMINER_STATE
	err = worker.InitLogMinerStateTable()
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation de la table LOGMINER_STATE : %v", err)
	}

	// Charger le dernier SCN avant de démarrer
	if err := worker.LoadLastProcessedSCN(); err != nil {
		log.Fatalf("Failed to load last processed SCN: %v", err)
	}

	worker.Start()
	defer worker.Stop()

	// Simule une exécution prolongée
	time.Sleep(3600 * time.Second)

}
