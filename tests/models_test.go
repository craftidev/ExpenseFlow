package tests

import (
    "testing"
    "time"
    "github.com/craftidev/expenseflow/internal/db"
	"github.com/craftidev/expenseflow/config"
)


func TestExpenseValid(t *testing.T) {
    validExpense := db.Expense{
        ID:        1,
        SessionID: 1,
        Amount:    db.Amount{Value: 100, Currency: "USD"},
        DateTime:  time.Now(),
        ReceiptURL: config.Path + "/assets/receipt_test.png",
    }

    err := validExpense.Valid()
    if err != nil {
        t.Errorf("expected valid expense, got error: %v", err)
    }

    invalidExpenses := make([]db.Expense, 5)

    invalidExpenses = append(invalidExpenses, db.Expense{
        ID:        0,  // Invalid ID
        SessionID: 1,
        Amount:    db.Amount{Value: 100, Currency: "USD"},
        DateTime:  time.Now(),
        ReceiptURL: config.Path + "/assets/receipt_test.png",
    })
    invalidExpenses = append(invalidExpenses, db.Expense{
        ID:        1,
        SessionID: 0,  // Invalid SessionID
        Amount:    db.Amount{Value: 100, Currency: "USD"},
        DateTime:  time.Now(),
        ReceiptURL: config.Path + "/assets/receipt_test.png",
    })
    invalidExpenses = append(invalidExpenses, db.Expense{
        ID:        1,
        SessionID: 0,
        Amount:    db.Amount{},  // Invalid Amount
        DateTime:  time.Now(),
        ReceiptURL: config.Path + "/assets/receipt_test.png",
    })
    invalidExpenses = append(invalidExpenses, db.Expense{
        ID:        1,
        SessionID: 0,
        Amount:    db.Amount{Value: 100, Currency: "USD"},
        DateTime:  time.Now(),
        ReceiptURL: config.Path + "",  // Invalid ReceiptURL
    })
    invalidDateTimeExpense := db.Expense{}
    invalidDateTimeExpense.ID = 1
    invalidDateTimeExpense.SessionID = 1
    invalidDateTimeExpense.Amount = db.Amount{Value: 100, Currency: "USD"}
    invalidDateTimeExpense.ReceiptURL = config.Path + "/assets/receipt_test.png"
    invalidExpenses = append(invalidExpenses, invalidDateTimeExpense)

    for _, invalidExpense := range invalidExpenses {
        err = invalidExpense.Valid()
        if err == nil {
            t.Error("expected error, got valid expense")
        }
    }
}

func TestHasReceipt(t *testing.T) {
    // Test case with default URL
    expenseWithDefaultURL := db.Expense{
        ReceiptURL: config.DefaultReceiptURL,
    }
    err := expenseWithDefaultURL.HasReceipt()
    if err != nil {
        t.Error("expected HasReceipt to return an error for default URL")
    }

    // Test case with non-existent file
    expenseWithNonExistentFile := db.Expense{
        ReceiptURL: config.Path + "/assets/receipts/nonexistent.png",
    }
    err = expenseWithNonExistentFile.HasReceipt()
    if err == nil {
        t.Error("expected HasReceipt to return an error for non-existent file")
    }

    // Test case with existing file
    expenseWithRealReceipt := db.Expense{
        ReceiptURL: config.Path + "/assets/receipts/receipt_test.png",
    }
    err = expenseWithRealReceipt.HasReceipt()
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
}
