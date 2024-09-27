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
	"github.com/craftidev/expenseflow/internal/utils"
)


func ConnectDB(DBPath string, ) (*sql.DB, error) {
    db, err := sql.Open("sqlite3", DBPath)
    if err != nil {
        return nil, utils.LogError("failed to open database: %v", err)
    }
    return db, nil
}

func InitDB(DBPath string, db *sql.DB) error {
    if _, err := os.Stat(DBPath); !os.IsNotExist(err) {
        log.Println("Database already exists. Skipping initialization.")
        return nil
    }

    filesPath := filepath.Join(config.MigrationsDirPath)
    schemaDirectory := os.DirFS(filesPath)
    schemaFiles, err := fs.Glob(schemaDirectory, "*.sql")
    if err != nil {
        return utils.LogError("failed to fetch schema files: %v", err)
    }
    sort.Strings(schemaFiles)

    for i, schemaFile := range schemaFiles {
        schema, err := os.ReadFile(filepath.Join(filesPath, schemaFile))
        if err != nil {
            return utils.LogError("failed to read schema file: %v", err)
        }

        _, err = db.Exec(string(schema))
        if err!= nil {
            return utils.LogError("failed to apply migration %03d: %v", i + 1, err)
        }
    }

    log.Println("Database created.")
    return nil
}

func CloseDB(db *sql.DB) error {
    if err := db.Close(); err != nil {
        return utils.LogError("failed to close database: %v", err)
    }
    log.Println("Database closed.")
    return nil
}
