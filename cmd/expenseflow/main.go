package main

import (
    "log"
    "os"
    "github.com/craftidev/expenseflow/internal/db"
)


func setupLogging() {
    logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatalf("Failed to open log file: %v", err)
    }

    log.SetOutput(logFile)

    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
    setupLogging()

    database := db.ConnectDB()
    defer db.CloseDB(database)

    db.InitDB(database)

    // TODO: Implement CRUD
    log.Println("ExpenseFlow DB connection established")
}
