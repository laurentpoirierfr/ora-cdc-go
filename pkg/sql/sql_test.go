package sql_test

import (
	"strings"
	"testing"

	"github.com/laurentpoirierfr/ora-cdc-go/pkg/sql"
)

func TestConvertOracleToPostgres(t *testing.T) {
	// Cas de test : entrées Oracle et sorties attendues en PostgreSQL
	tests := []struct {
		name        string
		oracleSQL   string
		expectedSQL string
	}{
		{
			name:        "Replace SYSDATE with CURRENT_TIMESTAMP",
			oracleSQL:   "SELECT SYSDATE FROM MY_TABLE",
			expectedSQL: "SELECT CURRENT_TIMESTAMP FROM MY_TABLE",
		},
		{
			name:        "Replace NVL with COALESCE",
			oracleSQL:   "SELECT NVL(col1, 'default') FROM MY_TABLE",
			expectedSQL: "SELECT COALESCE(col1, 'default') FROM MY_TABLE",
		},
		{
			name:        "Replace TO_CHAR with TO_DATE",
			oracleSQL:   "SELECT TO_CHAR(SYSDATE, 'YYYY-MM-DD') FROM MY_TABLE",
			expectedSQL: "SELECT CAST(CURRENT_TIMESTAMP, 'YYYY-MM-DD') FROM MY_TABLE",
		},
		{
			name:        "No replacement needed",
			oracleSQL:   "SELECT col1, col2 FROM MY_TABLE",
			expectedSQL: "SELECT col1, col2 FROM MY_TABLE",
		},
	}

	// Boucle sur chaque cas de test
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotSQL := strings.ToLower(sql.ConvertOracleToPostgres(tc.oracleSQL))
			expectedSQL := strings.ToLower(tc.expectedSQL)
			if gotSQL != expectedSQL {
				t.Errorf("ConvertOracleToPostgres(%q) = %q; want %q", tc.oracleSQL, gotSQL, expectedSQL)
			}
		})
	}
}
