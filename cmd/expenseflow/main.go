package main

import (
	"log"
	"os"

	"github.com/craftidev/expenseflow/config"
	"github.com/craftidev/expenseflow/internal/db"
)


func setupLogging() {
    logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatalf("[fatal] Failed to open log file: %v", err)
    }

    log.SetOutput(logFile)

    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
    setupLogging()

    database, err := db.ConnectDB(config.DBPath)
    if err != nil {
        log.Fatalf("[fatal] Failed to connect to database: %v", err)
    }

    defer func() {
        if err := db.CloseDB(database); err != nil {
            log.Fatalf("[fatal] Failed to close database: %v", err)
        }
    }()


    if err := db.InitDB(config.DBPath, database); err != nil {
        log.Fatalf("[fatal] Failed to initialize database: %v", err)
    }

    // TODO: Implement CRUD
    log.Println("[info] ExpenseFlow DB connection established")
}
