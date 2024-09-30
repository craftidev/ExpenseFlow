package config

import (
	"path/filepath"
	"runtime"
)


func getAppPath() string {
    // Get Dir of this particular file and go up the tree
    _, b, _, _ := runtime.Caller(0)
    basePath := filepath.Dir(filepath.Join(filepath.Dir(b)))
    return basePath
}

var (
    Path = getAppPath()
    DBPath = filepath.Join(Path, "internal", "db", "expenseflow.db")
    ReceiptsDir = filepath.Join(Path, "assets", "receipts")
    ReceiptsDirTest = filepath.Join(Path, "tests", "assets", "receipts")
    MigrationsDirPath = filepath.Join(Path, "internal", "db", "migrations")
)

const (
    MaxFloat = 1_000_000_000.0
)
