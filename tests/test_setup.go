package tests

import (
	"database/sql"
	"log"
	"runtime"
	"testing"

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

func InitializeSliceOfValidAny[T any](size int, valid T) []T {
	slice := make([]T, size)
	for i := range slice {
		slice[i] = valid
	}
	return slice
}

// Provoke vscode to get confused
// when a test fail because of entities validated in that function,
// the auto-opening feature try to opent this file
// with the path of the test file and open an inexistant file
// workaround -> return real calling location in t.Error
func ValidateEntities[T any](
	t *testing.T,
	entities []T,
	shouldReturnError bool,
	validator func(T) error,
) {
	for i, entity := range entities {
		err := validator(entity)
        _, file, line, _ := runtime.Caller(1)

        switch {
        case shouldReturnError && err == nil:
			t.Errorf(
				"expected error, got valid entity on test number: %d @ %v:%v",
				i, file, line,
			)
        case shouldReturnError && err != nil:
            continue
        case !shouldReturnError && err != nil:
			t.Errorf(
                "did not expect error, got error: %v on test number: %d @ %v:%v",
                err, i, file, line,
            )
        case !shouldReturnError && err == nil:
            continue
        default:
            log.Fatalf(
                "error on ValidateEntities function on test number: %d. Caller: %v:%v",
                i, file, line,
            )
        }
    }
}
