package tests

import (
	"database/sql"
	"log"

	"github.com/craftidev/expenseflow/internal/db"
)


var SingletonDatabaseTest *sql.DB

func SetupTestDatabase() *sql.DB {
	if SingletonDatabaseTest != nil {
		return SingletonDatabaseTest
	}

	var err error
	SingletonDatabaseTest, err = db.ConnectDB(":memory:")
	if err != nil {
		log.Fatalf("Failed to connect to in-memory database: %v", err)
	}

	if err := db.InitDB(":memory:", SingletonDatabaseTest); err != nil {
		log.Fatalf("Failed to initialize in-memory database: %v", err)
	}

	return SingletonDatabaseTest
}

func TeardownTestDatabase() {
	if SingletonDatabaseTest != nil {
		if err := db.CloseDB(SingletonDatabaseTest); err != nil {
			log.Fatalf("Failed to close database: %v", err)
		}
		SingletonDatabaseTest = nil
	}
}
