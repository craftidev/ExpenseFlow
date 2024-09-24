package tests

import (
    "testing"
    "time"
    "github.com/craftidev/expenseflow/internal/db"
	"github.com/craftidev/expenseflow/config"
)


// TODO Validates that the expense's amount is positive
// TODO Checks if the expense's receipt URL is valid when it's an empty string
// TODO Ensures that the expense's receipt URL is valid when it's a relative path
// TODO Verifies that the expense's receipt URL is valid when it's a valid URL
// TODO Tests the expense's Valid function with an amount having zero value
// TODO Checks the expense's Valid function with a session ID set to a negative value
// TODO Validates the expense's CheckReceipt function with a non-image file
// TODO Tests the expense's CheckReceipt function with a file that doesn't exist
// TODO Verifies the expense's CheckReceipt function with a valid image file
// TODO Checks the expense's CheckReceipt function with a URL that doesn't exist

func TestExpenseValid(t *testing.T) {
    validExpense := db.Expense{
        ID:        1,
        SessionID: 1,
        Amount:    db.Amount{Value: 100, Currency: "USD"},
        DateTime:  time.Now(),
        ReceiptURL: "/assets/receipt_test.png",
    }

    err := validExpense.Valid()
    if err != nil {
        t.Errorf("expected valid expense, got error: %v", err)
    }

    var invalidExpenses [5]db.Expense

    invalidExpenses[0] = db.Expense{
        ID:        0,               // 0: Invalid ID
        SessionID: 1,
        Amount:    db.Amount{Value: 100, Currency: "USD"},
        DateTime:  time.Now(),
        ReceiptURL: "/assets/receipt_test.png",
    }
    invalidExpenses[1] = db.Expense{
        ID:        1,
        SessionID: 0,               // 1: Invalid SessionID
        Amount:    db.Amount{Value: 100, Currency: "USD"},
        DateTime:  time.Now(),
        ReceiptURL: "/assets/receipt_test.png",
    }
    invalidExpenses[2] = db.Expense{
        ID:        1,
        SessionID: 1,
        Amount:    db.Amount{},     // 2: Invalid Amount
        DateTime:  time.Now(),
        ReceiptURL: "/assets/receipt_test.png",
    }
    invalidExpenses[3] = db.Expense{
        ID:        1,
        SessionID: 1,
        Amount:    db.Amount{Value: 100, Currency: "USD"},
        DateTime:  time.Time{},     // 3: Invalid DateTime
        ReceiptURL: "/assets/receipt_test.png",
    }
    invalidExpenses[4] = db.Expense{
        ID:        1,
        SessionID: 1,
        Amount:    db.Amount{Value: 100, Currency: "USD"},
        DateTime:  time.Now(),
        ReceiptURL: "",             // 4: Invalid ReceiptURL
    }

    for i, invalidExpense := range invalidExpenses {
        err = invalidExpense.Valid()
        if err == nil {
            t.Errorf("expected error, got valid expense on invalidExpense test number: %d", i)
        }
    }
}

func TestCheckReceipt(t *testing.T) {
    // Test case with default URL
    expenseWithDefaultURL := db.Expense{
        ReceiptURL: config.DefaultReceiptURL,
    }
    err := expenseWithDefaultURL.CheckReceipt()
    if err == nil {
        t.Error("expected CheckReceipt to return an error for default URL")
    }

    // Test case with non-existent file
    expenseWithNonExistentFile := db.Expense{
        ReceiptURL: "/assets/receipts/nonexistent.png",
    }
    err = expenseWithNonExistentFile.CheckReceipt()
    if err == nil {
        t.Error("expected CheckReceipt to return an error for non-existent file")
    }

    // Test case with existing file
    expenseWithRealReceipt := db.Expense{
        ReceiptURL: "/assets/receipts/receipt_test.png",
    }
    err = expenseWithRealReceipt.CheckReceipt()
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
}
