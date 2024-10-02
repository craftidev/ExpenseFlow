package models_tests

import (
	"testing"

	"github.com/craftidev/expenseflow/internal/db"
)


func TestClientValidation(t *testing.T) {
    // Test case with valid client
    client := db.Client{
        ID:   1,
        Name: "John Doe",
    }
    if err := client.Valid(); err != nil {
        t.Errorf("expected valid client, got error: %v", err)
    }

    // Test case with invalid client (zero ID)
    client.ID = 0
    if err := client.Valid(); err == nil {
        t.Error("expected error for invalid client with zero ID")
    }

    // Test case with valid client for PreInsert
    if err := client.PreInsertValid(); err != nil {
        t.Error("expected valid client for insertion with zero ID")
    }

    // Test case with invalid client (zero name)
    client.ID = 1
    client.Name = ""
    if err := client.Valid(); err == nil {
        t.Error("expected error for invalid client with zero name")
    }

    // Test case with invalid client (name lenght > 100 runes)
    client.ID = 1
    client.Name = "a" + string(make([]rune, 100))
    if err := client.Valid(); err == nil {
        t.Error("expected error for invalid client with name lenght > 100 runes")
    }
}
