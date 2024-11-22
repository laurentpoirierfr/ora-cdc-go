package model

import "time"

// LogMinerConfig contient les paramètres pour la configuration
type LogMinerConfig struct {
	DBConnectionString string
	PollFrequency      time.Duration
	Callback           func(row LogMinerRow)
}

// Struct de la view LogMiner Oracle
type LogMinerRow struct {
	SCN       int64     `json:"scn"`
	Timestamp time.Time `json:"timestamp"`
	Operation string    `json:"operation"`
	SegOwner  string    `json:"seg_owner"`
	TableName string    `json:"table_name"`
	SQLRedo   string    `json:"sql_redo"`
	SQLUndo   string    `json:"sql_undo"`
	RowID     string    `json:"row_id"`
	Username  string    `json:"username"`
	Session   int       `json:"session"`
	Rollback  string    `json:"rollback"`
}
