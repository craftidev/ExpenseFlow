package crud_tests

import (
	"database/sql"
	"os"
	"reflect"
	"testing"

	"github.com/craftidev/expenseflow/internal/db/crud"
	"github.com/craftidev/expenseflow/tests"
)


var DatabaseTest *sql.DB

func TestMain(m *testing.M) {
	DatabaseTest = tests.SetupTestDatabase()

	exitCode := m.Run()

	tests.TeardownTestDatabase()

	os.Exit(exitCode)
}

var validClient = tests.GetValidClient()
var validSession = tests.GetValidSession()
var validCarTrip = tests.GetValidCarTrip()
var validExpenseType = tests.GetValidExpenseType()
var validExpense = tests.GetValidExpense()
var validLineItem = tests.GetValidLineItem()

func TestCreateModels(t *testing.T) {
    idClient, errClient := crud.CreateClient(DatabaseTest, validClient)
    idSession, errSession := crud.CreateSession(DatabaseTest, validSession)
    idCarTrip, errCarTrip := crud.CreateCarTrip(DatabaseTest, validCarTrip)
    idExpenseType, errExpenseType := crud.CreateExpenseType(DatabaseTest, validExpenseType)
    idExpense, errExpense := crud.CreateExpense(DatabaseTest, validExpense)
    idLineItem, errLineItem := crud.CreateLineItem(DatabaseTest, validLineItem)

    errsValid := []struct {
        name string
        id int64
        err  error
    }{
        {"Client", idClient, errClient},
        {"Session", idSession, errSession},
        {"CarTrip", idCarTrip, errCarTrip},
        {"ExpenseType", idExpenseType, errExpenseType},
        {"Expense", idExpense, errExpense},
        {"LineItem", idLineItem, errLineItem},
    }
    for _, e := range errsValid {
        if e.err != nil {
            t.Errorf("expected no error on INSERT for %s, got: %v", e.name, e.err)
        }
        if e.id <= 0 {
            t.Errorf("expected non-zero and positive ID for %s", e.name)
        }
    }


    //InvalidActions
    // UNIQUE
    _, errClient = crud.CreateClient(DatabaseTest, validClient)
    _, errCarTrip = crud.CreateCarTrip(DatabaseTest, validCarTrip)
    _, errExpenseType = crud.CreateExpenseType(DatabaseTest, validExpenseType)

    errsInvalid := []struct {
        name string
        err  error
    }{
        {"Client", errClient},
        {"CarTrip", errCarTrip},
        {"ExpenseType", errExpenseType},
    }
    for _, e := range errsInvalid {
        if e.err == nil {
            t.Errorf(
                "expected error on INSERT because of UNIQUE field for %s, " +
                "got no error", e.name,
            )
        }
    }

}

// Can't run this test alone, need INSERTS
func TestGetModelsByID(t *testing.T) {
    client, errClient := crud.GetClientByID(DatabaseTest, 1)
    session, errSession := crud.GetSessionByID(DatabaseTest, 1)
    carTrip, errCarTrip := crud.GetCarTripByID(DatabaseTest, 1)
    expenseType, errExpenseType := crud.GetExpenseTypeByID(DatabaseTest, 1)
    expense, errExpense := crud.GetExpenseByID(DatabaseTest, 1)
    lineItem, errLineItem := crud.GetLineItemByID(DatabaseTest, 1)

    errsValid := []struct {
        name string
        err  error
    }{
        {"Client", errClient},
        {"Session", errSession},
        {"CarTrip", errCarTrip},
        {"ExpenseType", errExpenseType},
        {"Expense", errExpense},
        {"LineItem", errLineItem},
    }
    for _, e := range errsValid {
        if e.err != nil {
            t.Errorf("expected no error on SELECT by ID for %s, got: %v", e.name, e.err)
        }
    }

    compareEntities(t, "Client", validClient, *client)
    compareEntities(t, "CarTrip", validCarTrip, *carTrip)
    compareEntities(t, "ExpenseType", validExpenseType, *expenseType)
    compareEntities(t, "LineItem", validLineItem, *lineItem)
    // Custom manual checks because of time.Time badly handled with deepEqual
    if  session.ID != validSession.ID ||
        session.ClientID != validSession.ClientID ||
        session.Location != validSession.Location ||
        session.TripStartLocation != validSession.TripStartLocation ||
        session.TripEndLocation != validSession.TripEndLocation ||
        !session.StartAtDateTime.Equal(validSession.StartAtDateTime) ||
        !session.EndAtDateTime.Equal(validSession.EndAtDateTime) {
            t.Errorf("data fetched for Session (%v) doesn't match data inserted (%v)",
            session, validSession,
        )
    }
    if  expense.ID != validExpense.ID ||
        expense.SessionID != validExpense.SessionID ||
        expense.TypeID != validExpense.TypeID ||
        expense.Currency != validExpense.Currency ||
        expense.ReceiptRelPath != validExpense.ReceiptRelPath ||
        expense.Notes != validExpense.Notes ||
        !expense.DateTime.Equal(validExpense.DateTime) {
            t.Errorf("data fetched for Expense (%v) doesn't match data inserted (%v)",
            expense, validExpense,
        )
    }
}

func compareEntities(t *testing.T, name string, expected, actual interface{}) {
    if !reflect.DeepEqual(expected, actual) {
        t.Errorf(
            "Data fetched for %s (%v) doesn't match data inserted (%v)",
            name, actual, expected,
        )
    }
}

// FK
// Client: FK in sessions
// Session: FK in car_trips, expenses
// ExpenseType: FK in expenses
// Expense: FK in line_items
