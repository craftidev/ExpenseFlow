package db

import (
	"database/sql"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/mattn/go-sqlite3"

	"github.com/craftidev/expenseflow/config"
)


func ConnectDB() *sql.DB {
    db, err := sql.Open("sqlite3", config.DBPath)
    if err != nil {
        log.Fatalf("failed to open database: %v", err)
    }
    return db
}

func InitDB(db *sql.DB) {
    if _, err := os.Stat(config.DBPath); !os.IsNotExist(err) {
        log.Println("Database already exists. Skipping initialization.")
        return
    }

    filesPath := filepath.Join(config.MigrationsDirPath)
    schemaDirectory := os.DirFS(filesPath)
    schemaFiles, err := fs.Glob(schemaDirectory, "*.sql")
    if err != nil {
        log.Fatal(err)
    }
    sort.Strings(schemaFiles)

    for i, schemaFile := range schemaFiles {
        schema, err := os.ReadFile(filepath.Join(filesPath, schemaFile))
        if err != nil {
            log.Fatalf("failed to read schema file: %v", err)
        }

        _, err = db.Exec(string(schema))
        if err!= nil {
            log.Fatalf("failed to apply migration %03d: %v", i + 1, err)
        }
    }

    log.Println("Database created.")
}

func CloseDB(db *sql.DB) {
    if err := db.Close(); err != nil {
        log.Fatalf("failed to close database: %v", err)
    }
    log.Println("Database closed.")
}
