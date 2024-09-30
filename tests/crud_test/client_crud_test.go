package crud_test

import (
	"database/sql"
	"os"
	"strconv"
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

func TestGetClientByID(t *testing.T) {
    // Can't run this test alone because use TestCreateClient valid entry
    // Valid fetch
    fetchedClient, err := crud.GetClientByID(DatabaseTest, 1)
    if err != nil {
        t.Errorf("expected no error, got: %v", err)
    }
    if fetchedClient.Valid() != nil {
        t.Error("expected valid client, got invalid")
    }

    // Invalid fetch
    _, err = crud.GetClientByID(DatabaseTest, 9999)
    if err == nil {
        t.Error("expected error with none existant id, got no error")
    }
}

func TestUpdateClient(t *testing.T) {
    // Can't run this test alone because use TestCreateClient valid entry

    // Valid update
    client := db.Client{
        ID:   1,
        Name: "John Doe Updated",
    }
    err := crud.UpdateClient(DatabaseTest, client)
    if err!= nil {
        t.Errorf("expected no error, got: %v", err)
    }

    // Invalid update (zero ID)
    client.ID = 0
    err = crud.UpdateClient(DatabaseTest, client)
    if err == nil {
        t.Error("expected error for invalid client with zero ID")
    }

	var invalidClients [4]db.Client
	for i := 0; i < len(invalidClients); i++ {
		invalidClients[i] = client
        invalidClients[i].Name += strconv.Itoa(i)
	}

	invalidClients[0].ID = -1
	invalidClients[1].ID = 0
	invalidClients[2].ID = 9999
	invalidClients[3].Name = "John Doe Updated"

	for i, invalidClient := range invalidClients {
		err = crud.UpdateClient(DatabaseTest, invalidClient)
		if err == nil {
			t.Errorf(
                "expected error, got valid client on UpdateClient method " +
                "test number: %d", i,
            )
		}
	}
}

func TestDeleteClientByID(t *testing.T) {
    // Can't run this test alone because use TestCreateClient valid entry

    err := crud.DeleteClientByID(DatabaseTest, 1)
    if err != nil {
        t.Errorf("expected no error, got: %v", err)
    }

    err = crud.DeleteClientByID(DatabaseTest, 1)
    if err == nil {
        t.Errorf("expected an error, client doesn't exist, got: %v", err)
    }
    err = crud.DeleteClientByID(DatabaseTest, 0)
    if err == nil {
        t.Errorf("expected an error, id zero-value, got: %v", err)
    }
    err = crud.DeleteClientByID(DatabaseTest, -1)
    if err == nil {
        t.Errorf("expected an error, id negative, got: %v", err)
    }
}
