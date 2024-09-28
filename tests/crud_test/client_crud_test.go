package crud_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/craftidev/expenseflow/internal/db"
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

func TestCreateClient(t *testing.T) {
    // Test case with valid client
    client := db.Client{
        Name: "John Doe",
    }
    id, err := crud.CreateClient(DatabaseTest, client)
    if err != nil {
        t.Errorf("expected no error, got: %v", err)
    }
    if id == 0 {
        t.Error("expected non-zero ID, got 0")
    }

    // Test case with invalid client (zero name)
    client.Name = ""
    _, err = crud.CreateClient(DatabaseTest, client)
    if err == nil {
        t.Error("expected error for invalid client with zero name")
    }

    // Test case with existing name
    client.Name = "John Doe"
    _, err = crud.CreateClient(DatabaseTest, client)
    if err == nil {
        t.Error("expected error for existing client name")
    }

    // Test case with name length > 100
    client.ID = 1
    client.Name = "a" + string(make([]rune, 100))
    _, err = crud.CreateClient(DatabaseTest, client)
    if err == nil {
        t.Error("expected error for name length > 100")
    }
}

// For GetClientByID()
    // Test case with fetching valid client
    // fetchedClient, err := crud.GetClientByID(DatabaseTest, id)
    // if err!= nil {
    //     t.Errorf("expected no error, got: %v", err)
    // }
    // if fetchedClient == nil {
    //     t.Error("expected valid client, got nil")
    // }
    // if fetchedClient.ID!= client.ID || fetchedClient.Name!= client.Name {
    //     t.Errorf("expected client %+v, got %+v", client, fetchedClient)
    // }
