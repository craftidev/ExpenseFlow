package tests

import (
	"os"
	"testing"
	"time"

	"github.com/craftidev/expenseflow/config"
	"github.com/craftidev/expenseflow/internal/db"
)


func TestExpenseValid(t *testing.T) {
    validExpense := db.Expense{
        ID:        1,
        SessionID: 1,
        Amount:    db.Amount{Value: 100, Currency: "USD"},
        DateTime:  time.Now(),
        ReceiptURL: "/tests/assets/valid_receipt_test.png",
    }

    err := validExpense.Valid()
    if err != nil {
        t.Errorf("expected valid expense, got error: %v", err)
    }

    var invalidExpenses [11]db.Expense

    invalidExpenses[0] = db.Expense{
        ID:        0,           // 0: Invalid ID (0)
        SessionID: 1,
        Amount:    db.Amount{Value: 100, Currency: "USD"},
        DateTime:  time.Now(),
        ReceiptURL: "/tests/assets/valid_receipt_test.png",
    }
    invalidExpenses[1] = db.Expense{
        ID:        1,
        SessionID: 0,           // 01: Invalid SessionID (0)
        Amount:    db.Amount{Value: 100, Currency: "USD"},
        DateTime:  time.Now(),
        ReceiptURL: "/tests/assets/valid_receipt_test.png",
    }
    invalidExpenses[2] = db.Expense{
        ID:        1,
        SessionID: 1,
        Amount:    db.Amount{}, // 02: Invalid Amount
        DateTime:  time.Now(),
        ReceiptURL: "/tests/assets/valid_receipt_test.png",
    }
    invalidExpenses[3] = db.Expense{
        ID:        1,
        SessionID: 1,
        Amount:    db.Amount{Value: 100, Currency: "USD"},
        DateTime:  time.Time{}, // 03: Invalid DateTime
        ReceiptURL: "/tests/assets/valid_receipt_test.png",
    }
    invalidExpenses[4] = db.Expense{
        ID:        1,
        SessionID: 1,
        Amount:    db.Amount{Value: 100, Currency: "USD"},
        DateTime:  time.Now(),
        ReceiptURL: "",         // 04: Invalid ReceiptURL
    }
    invalidExpenses[5] = db.Expense{
        ID:        1,
        SessionID: 1,
        Amount:    db.Amount{Value: -100, Currency: "USD"},
                                // 05: Invalid neg Amount
        DateTime:  time.Now(),
        ReceiptURL: "/tests/assets/valid_receipt_test.png",
    }
    invalidExpenses[6] = db.Expense{
        ID:        1,
        SessionID: 1,
        Amount:    db.Amount{Value: 100},
                                // 06: Invalid currency
        DateTime:  time.Now(),
        ReceiptURL: "/tests/assets/valid_receipt_test.png",
    }
    invalidExpenses[7] = db.Expense{
        ID:        1,
        SessionID: 1,
        Amount:    db.Amount{Currency: "USD"},
                                // 07: Invalid value (empty)
        DateTime:  time.Now(),
        ReceiptURL: "/tests/assets/valid_receipt_test.png",
    }
    invalidExpenses[8] = db.Expense{
        ID:        1,
        SessionID: 1,
        Amount:    db.Amount{Value: 0, Currency: "USD"},
                                // 08: Invalid value (0)
        DateTime:  time.Now(),
        ReceiptURL: "/tests/assets/valid_receipt_test.png",
    }
    invalidExpenses[9] = db.Expense{
        ID:        -1,          // 09: Invalid ID (negative)
        SessionID: 1,
        Amount:    db.Amount{Value: 100, Currency: "USD"},
        DateTime:  time.Now(),
        ReceiptURL: "/tests/assets/valid_receipt_test.png",
    }
    invalidExpenses[10] = db.Expense{
        ID:        1,
        SessionID: -1,           // 10: Invalid SessionID (negative)
        Amount:    db.Amount{Value: 100, Currency: "USD"},
        DateTime:  time.Now(),
        ReceiptURL: "/tests/assets/valid_receipt_test.png",
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
        ReceiptURL: "/tests/assets/receipts/nonexistent.png",
    }
    err = expenseWithNonExistentFile.CheckReceipt()
    if err == nil {
        t.Error("expected CheckReceipt to return an error for non-existent file")
    }

    // Test case with existing file and correct image non-empty
    expenseWithRealReceipt := db.Expense{
        ReceiptURL: "/tests/assets/receipts/valid_receipt_test.png",
    }
    err = expenseWithRealReceipt.CheckReceipt()
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }

    // Test case with existing non-image file
    expenseWithBadExistingReceipt := db.Expense{
        ID:        1,
        SessionID: 1,
        Amount:    db.Amount{Value: 100, Currency: "USD"},
        DateTime:  time.Now(),
        ReceiptURL: "/tests/assets/receipts/invalid_receipt_test",
    }
    err = expenseWithBadExistingReceipt.CheckReceipt()
    if err == nil {
        t.Errorf("expected CheckReceipt to return an error for non-compatible file")
    }

    // Test case with existing image file but wrong image type
    // TODO: Turns out when the headers are good everythings good.
    // Check later with flutter and reverse the test to not raise err
    //
    // expenseWithBadExistingReceiptImageType := db.Expense{
    //     ID:        1,
    //     SessionID: 1,
    //     Amount:    db.Amount{Value: 100, Currency: "USD"},
    //     DateTime:  time.Now(),
    //     ReceiptURL: "/tests/assets/receipts/valid_receipt_wrong_extension_test.wrong",
    // }
    // err = expenseWithBadExistingReceiptImageType.CheckReceipt()
    // if err == nil {
    //     t.Errorf("expected CheckReceipt to return an error for non-compatible file")
    // }
}

func TestCheckReceiptWithProtectedFile(t *testing.T) {
    filePath := "/tests/assets/receipts/protected_receipt_test.png"

    // Set permissions to 000
    err := os.Chmod(config.Path + filePath, 0000)
    if err != nil {
        t.Fatalf("failed to set file permissions before testing: %v", err)
    }

    // Run the test
    expense := db.Expense{
        ReceiptURL: filePath,
    }
    err = expense.CheckReceipt()
    if err == nil {
        t.Error("expected error due to permission denied")
    }

    // Restore permissions after the test
    err = os.Chmod(config.Path +filePath, 0644)
    if err != nil {
        t.Fatalf("failed to restore file permissions: %v", err)
    }
}
