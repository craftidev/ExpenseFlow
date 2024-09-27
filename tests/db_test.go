package tests

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/craftidev/expenseflow/internal/db"
)


var database *sql.DB
func TestMain(m *testing.M) {
    // Use SQLite's in-memory DB for testing
    var err error
    database, err = db.ConnectDB(":memory:")
    if err != nil {
        log.Fatalf("Failed to connect to in-memory database: %v", err)
    }

    if err := db.InitDB(":memory:", database); err != nil {
        log.Fatalf("Failed to initialize in-memory database: %v", err)
    }

    exitCode := m.Run()

    if err := db.CloseDB(database); err != nil {
        log.Fatalf("Failed to close database: %v", err)
    }
    os.Exit(exitCode)
}

func TestInitDB(t *testing.T) {
    if database == nil {
        t.Fatal("Expected the in-memory database to be set up, but it was nil")
    }

    // Test some initial state
    rows, err := database.Query("SELECT name FROM sqlite_master WHERE type='table'")
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
