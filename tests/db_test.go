package tests

import (
    "database/sql"
    "testing"

    "github.com/craftidev/expenseflow/internal/db"
)

func setupInMemoryDB() (*sql.DB, error) {
    // Use SQLite's in-memory DB for testing
    database, err := db.ConnectDB(":memory:")
    if err != nil {
        return nil, err
    }

    if err := db.InitDB(":memory:", database); err != nil {
        return nil, err
    }

    return database, nil
}

func TestConnectDB(t *testing.T) {
    database, err := setupInMemoryDB()
    if err != nil {
        t.Fatalf("Failed to connect to in-memory database: %v", err)
    }

    defer db.CloseDB(database)
}

func TestInitDB(t *testing.T) {
    database, err := setupInMemoryDB()
    if err != nil {
        t.Fatalf("Failed to initialize in-memory database: %v", err)
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

    db.CloseDB(database)
}

func TestCloseDB(t *testing.T) {
    database, err := setupInMemoryDB()
    if err != nil {
        t.Fatalf("Failed to connect to in-memory database: %v", err)
    }

    if err := db.CloseDB(database); err != nil {
        t.Errorf("Failed to close database: %v", err)
    }
}
