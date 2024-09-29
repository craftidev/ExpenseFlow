package config

import (
	"log"
	"os"
	"path/filepath"
)

func GetAppPath() string {
    path, err := os.Getwd()
    if err != nil {
        log.Fatalf("[fatal] couldn't find working directory: %v", err)
    }
    return path
}

var (
    DBPath = filepath.Join(GetAppPath(), "internal", "db", "expenseflow.db")
    DefaultReceiptURL = filepath.Join(GetAppPath(), "assets", "receipts", "default.jpg") // TODO get rid of it
    MigrationsDirPath = filepath.Join(GetAppPath(),"internal", "db", "migrations")
)
const (
    MaxFloat = 1_000_000_000.0
)
