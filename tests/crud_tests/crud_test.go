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

func TestUpdateModels(t *testing.T) {
    client := validClient
    session := validSession
    carTrip := validCarTrip
    expenseType := validExpenseType
    expense := validExpense
    lineItem := validLineItem

    clientChange := "updated"
    sessionChange := "updated"
    carTripChange := "1954-10-03"
    expenseTypeChange := "updated"
    expenseChange := sql.NullString{String: "updated", Valid: true}
    lineItemChange := 111.1

    client.Name = clientChange
    session.Location = sessionChange
    carTrip.DateOnly = carTripChange
    expenseType.Name = expenseTypeChange
    expense.Notes = expenseChange
    lineItem.Total = lineItemChange

    errClient := crud.UpdateClient(DatabaseTest, client)
    errSession := crud.UpdateSession(DatabaseTest, session)
    errCarTrip := crud.UpdateCarTrip(DatabaseTest, carTrip)
    errExpenseType := crud.UpdateExpenseType(DatabaseTest, expenseType)
    errExpense := crud.UpdateExpense(DatabaseTest, expense)
    errLineItem := crud.UpdateLineItem(DatabaseTest, lineItem)

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
            t.Errorf("expected no error on UPDATE for %s, got: %v", e.name, e.err)
        }
    }

    if  clientChange != client.Name ||
        sessionChange != session.Location ||
        carTripChange != carTrip.DateOnly ||
        expenseTypeChange != expenseType.Name ||
        expenseChange != expense.Notes ||
        lineItemChange != lineItem.Total {
            t.Error("data UPDATEd doesn't match data changed")
    }

    // UNIQUE Invalid re-update
    carTripCreation := validCarTrip // Avoid FK shenanigans later
    carTripCreation.SessionID = sql.NullInt64{Valid: false}
    _, errClientCreation := crud.CreateClient(DatabaseTest, validClient)
    _, errCarTripCreation := crud.CreateCarTrip(DatabaseTest, carTripCreation)
    _, errExpenseTypeCreation := crud.CreateExpenseType(DatabaseTest, validExpenseType)
    if  errClientCreation != nil ||
        errCarTripCreation != nil ||
        errExpenseTypeCreation != nil {
            t.Error("Unexpected error on additional creation for UNIQUE tests in UPDATE tests")
    }

    client.Name = validClient.Name
    carTrip.DateOnly = validCarTrip.DateOnly
    expenseType.Name = validExpenseType.Name

    errClient = crud.UpdateClient(DatabaseTest, client)
    errCarTrip = crud.UpdateCarTrip(DatabaseTest, carTrip)
    errExpenseType = crud.UpdateExpenseType(DatabaseTest, expenseType)
    if errClient == nil || errCarTrip == nil || errExpenseType == nil {
        t.Error("expected errors on data UPDATEd with same UNIQUE entry")
    }
}

func TestDeleteModels(t *testing.T) {
    // Invalid delete because of FK
    // ExpenseType: FK in expenses
    // Expense: FK in line_items
    // Session: FK in car_trips, expenses
    // Client: FK in sessions
    expectedErrClient := crud.DeleteClientByID(DatabaseTest, 1)
    expectedErrSession := crud.DeleteSessionByID(DatabaseTest, 1)
    expectedErrExpenseType := crud.DeleteExpenseTypeByID(DatabaseTest, 1)
    expectedErrExpense := crud.DeleteExpenseByID(DatabaseTest, 1)

    errsInvalid := []struct {
        name string
        err  error
    }{
        {"Client", expectedErrClient},
        {"Session", expectedErrSession},
        {"ExpenseType", expectedErrExpenseType},
        {"Expense", expectedErrExpense},
    }
    for _, e := range errsInvalid {
        if e.err == nil {
            t.Errorf(
                "expected error on DELETE because of FK for %s, got: %v",
                e.name, e.err,
            )
        }
    }

    // Valid DELETEs
    errLineItem := crud.DeleteLineItemByID(DatabaseTest, 1)
    errCarTrip := crud.DeleteCarTripByID(DatabaseTest, 1)
    errExpense := crud.DeleteExpenseByID(DatabaseTest, 1)
    errExpenseType := crud.DeleteExpenseTypeByID(DatabaseTest, 1)
    errSession := crud.DeleteSessionByID(DatabaseTest, 1)
    errClient := crud.DeleteClientByID(DatabaseTest, 1)

    errsValid := []struct {
        name string
        err  error
    }{
        {"LineItem", errLineItem},
        {"CarTrip", errCarTrip},
        {"Expense", errExpense},
        {"ExpenseType", errExpenseType},
        {"Session", errSession},
        {"Client", errClient},
    }
    for _, e := range errsValid {
        if e.err != nil {
            t.Fatalf("expected no error on DELETE for %s, got: %v", e.name, e.err)
        }
    }

    _, errLineItemFetch := crud.GetLineItemByID(DatabaseTest, 1)
    _, errCarTripFetch := crud.GetCarTripByID(DatabaseTest, 1)
    _, errExpenseFetch := crud.GetExpenseByID(DatabaseTest, 1)
    _, errExpenseTypeFetch := crud.GetExpenseTypeByID(DatabaseTest, 1)
    _, errSessionFetch := crud.GetSessionByID(DatabaseTest, 1)
    _, errClientFetch := crud.GetClientByID(DatabaseTest, 1)

    errsInvalid = []struct {
        name string
        err  error
    }{
        {"LineItem", errLineItemFetch},
        {"CarTrip", errCarTripFetch},
        {"Expense", errExpenseFetch},
        {"ExpenseType", errExpenseTypeFetch},
        {"Session", errSessionFetch},
        {"Client", errClientFetch},
    }
    for _, e := range errsInvalid {
        if e.err == nil {
            t.Errorf(
                "expected error on fetching DELETEd data for %s, got: %v",
            e.name, e.err,
        )
        }
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
