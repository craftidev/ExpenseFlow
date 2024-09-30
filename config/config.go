package config

import (
	"log"
	"os"
	"path/filepath"
)

func getAppPath() string {
    path, err := os.Getwd()
    if err != nil {
        log.Fatalf("[fatal] couldn't find working directory: %v", err)
    }
    return path
}

var (
    Path = getAppPath()
    DBPath = filepath.Join(Path, "internal", "db", "expenseflow.db")
    ReceiptsDir = filepath.Join(Path, "assets", "receipts")
    MigrationsDirPath = filepath.Join(Path,"internal", "db", "migrations")
)
const (
    MaxFloat = 1_000_000_000.0
)
