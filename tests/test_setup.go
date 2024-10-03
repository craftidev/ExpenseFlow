package tests

import (
	"database/sql"
	"log"
	"runtime"
	"testing"
	"time"

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
				"expected error, got valid entity on test number: %d\n@ %v:%v",
				i, file, line,
			)
        case shouldReturnError && err != nil:
            continue
        case !shouldReturnError && err != nil:
			t.Errorf(
                "did not expect error on test number: %d\ngot error: %v\n@ %v:%v",
                i, err, file, line,
            )
        case !shouldReturnError && err == nil:
            continue
        default:
            log.Fatalf(
                "error on ValidateEntities function on test number: %d.\nCaller: %v:%v",
                i, file, line,
            )
        }
    }
}

func GetValidCarTrip() db.CarTrip {
    return db.CarTrip{
        ID:         1,
        SessionID:  sql.NullInt64{Int64: 1, Valid: true},
        DistanceKM: 50.5,
        DateOnly:  "2022-01-01",
    }
}

func GetValidClient() db.Client {
    return db.Client{
        ID:   1,
        Name: "John Doe",
    }
}

func GetValidExpense() db.Expense {
	return db.Expense{
		ID:             1,
		SessionID:      sql.NullInt64{Int64: 1, Valid: true},
		TypeID:         1,
		Currency:       "USD",
		ReceiptRelPath: sql.NullString{
            String: "valid_receipt_test.png",
            Valid: true,
        },
		Notes:          sql.NullString{
            String: "Valid notes here.",
            Valid: true,
        },
		DateTime:       time.Now(),
	}
}

func GetValidExpenseType() db.ExpenseType {
    return db.ExpenseType{
        ID:   1,
        Name: "Transportation",
    }
}

func GetValidLineItem() db.LineItem {
    return db.LineItem{
        ID:   1,
        ExpenseID: 1,
        TaxeRate: 5.5,
        Total: 15.43,
    }
}

func GetValidSession() db.Session {
	return db.Session{
		ID:                1,
		ClientID:          1,
		Location:          "New York City (NY)",
		TripStartLocation: sql.NullString{
            String: "Philadelphia (PA)",
            Valid: true,
        },
		TripEndLocation: sql.NullString{
            String: "Boston (MA)",
            Valid: true,
        },
		StartAtDateTime:   db.NullableTime{
            Time: time.Now(),
            Valid: true,
        },
		EndAtDateTime:   db.NullableTime{
            Time: time.Now().Add(24 * time.Hour),
            Valid: true,
        },
    }
}
