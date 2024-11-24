package sql

import (
	"fmt"
	"regexp"
	"strings"
)

// ConvertOracleToPostgres transforme une requête SQL Oracle en SQL PostgreSQL.
func ConvertOracleToPostgres(oracleSQL string) string {
	// Étape 1 : Remplacer ROWNUM par LIMIT
	rownumRegex := regexp.MustCompile(`(?i)\bROWNUM\s*<=\s*(\d+)`)
	oracleSQL = rownumRegex.ReplaceAllString(oracleSQL, "LIMIT $1")

	// Étape 2 : Remplacer SYSDATE par CURRENT_TIMESTAMP
	oracleSQL = strings.ReplaceAll(strings.ToUpper(oracleSQL), "SYSDATE", "CURRENT_TIMESTAMP")

	// Étape 3 : Remplacer les double quotes par des guillemets simples pour PostgreSQL
	oracleSQL = strings.ReplaceAll(oracleSQL, "\"", "'")

	// Étape 4 : Remplacer les fonctions spécifiques Oracle
	replacements := map[string]string{
		"NVL":       "COALESCE",
		"TO_DATE":   "TO_TIMESTAMP",
		"TO_NUMBER": "CAST",
		"TO_CHAR":   "CAST",
		"TRUNC":     "DATE_TRUNC",
		"LENGTH":    "CHAR_LENGTH",
		"SUBSTR":    "SUBSTRING",
		"SYSDATE":   "CURRENT_TIMESTAMP",
		"DUAL":      "", // PostgreSQL n'a pas besoin de la table DUAL
	}

	for oracleFunc, postgresFunc := range replacements {
		regex := regexp.MustCompile(fmt.Sprintf(`(?i)\b%s\b`, oracleFunc))
		oracleSQL = regex.ReplaceAllString(oracleSQL, postgresFunc)
	}

	// Étape 5 : Gérer les types de données spécifiques Oracle -> PostgreSQL
	typeMappings := map[string]string{
		"NUMBER":   "NUMERIC",
		"VARCHAR2": "VARCHAR",
		"DATE":     "TIMESTAMP",
		"CLOB":     "TEXT",
		"BLOB":     "BYTEA",
	}
	for oracleType, postgresType := range typeMappings {
		regex := regexp.MustCompile(fmt.Sprintf(`(?i)\b%s\b`, oracleType))
		oracleSQL = regex.ReplaceAllString(oracleSQL, postgresType)
	}

	// Étape 6 : Gérer les séquences Oracle -> PostgreSQL
	sequenceRegex := regexp.MustCompile(`(?i)(\w+)\.NEXTVAL`)
	oracleSQL = sequenceRegex.ReplaceAllString(oracleSQL, "nextval('$1')")

	// Étape 7 : Nettoyer les espaces en trop
	oracleSQL = strings.TrimSpace(oracleSQL)

	return oracleSQL
}
