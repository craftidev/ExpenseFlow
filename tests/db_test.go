package tests

import (
	"database/sql"
	"os"
	"testing"
)


var DatabaseTest *sql.DB

func TestMain(m *testing.M) {
    DatabaseTest = SetupTestDatabase()

    exitCode := m.Run()

    TeardownTestDatabase()

    os.Exit(exitCode)
}

// TODO test standard model migration
func TestInitDB(t *testing.T) {
	if SingletonDatabaseTest == nil {
		t.Fatal("Expected the in-memory database to be set up, but it was nil")
	}

	// Test some initial state
	rows, err := SingletonDatabaseTest.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		t.Fatalf("Failed to query table names: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		tables = append(tables, name)
	}

	if len(tables) == 0 {
		t.Error("Expected some tables to be created, but none were found")
	}
}
