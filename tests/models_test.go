package tests

import (
	"fmt"
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
        Amount:    db.Amount{Value: 100.20, Currency: "USD"},
        DateTime:  time.Now(),
        ReceiptURL: "/tests/assets/valid_receipt_test.png",
    }
    err := validExpense.Valid()
    if err != nil {
        t.Errorf("expected valid expense, got error: %v", err)
    }

    var invalidExpenses [11]db.Expense
    for i := 0; i < len(invalidExpenses); i++ {
        invalidExpenses[i] = validExpense
    }

    invalidExpenses[0].ID           = 0
    invalidExpenses[1].SessionID    = 0
    invalidExpenses[2].Amount       = db.Amount{}
    invalidExpenses[3].DateTime     = time.Time{}
    invalidExpenses[4].ReceiptURL   = ""
    invalidExpenses[5].Amount.Value = -100
    invalidExpenses[6].Amount       = db.Amount{Value: 100.20}
    invalidExpenses[7].Amount       = db.Amount{Currency: "USD"}
    invalidExpenses[8].Amount.Value = 0
    invalidExpenses[9].ID           = -1
    invalidExpenses[10].SessionID   = -1

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

    // Test case with corrupted image file
    expenseWithReceiptNonImage := db.Expense{
        ID:        1,
        SessionID: 1,
        Amount:    db.Amount{Value: 100.20, Currency: "USD"},
        DateTime:  time.Now(),
        ReceiptURL: "/tests/assets/receipts/invalid_receipt_test.txt",
    }
    err = expenseWithReceiptNonImage.CheckReceipt()
    if err == nil {
        t.Errorf("expected CheckReceipt to return an error for corrupted file")
    }

    // Test case with corrupted image file
    expenseWithBadExistingReceipt := db.Expense{
        ID:        16,
        SessionID: 1,
        Amount:    db.Amount{Value: 100.20, Currency: "USD"},
        DateTime:  time.Now(),
        ReceiptURL: "/tests/assets/receipts/corrupted_receipt_test.png",
    }
    err = expenseWithBadExistingReceipt.CheckReceipt()
    if err == nil {
        t.Errorf("expected CheckReceipt to return an error for corrupted file")
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

// TODO test Sum / Add
func TestAmountAdd(t *testing.T) {
    // Valid Add
    amountA := db.Amount{Value: 100.20, Currency: "USD"}
    amountB := db.Amount{Value: 200.20, Currency: "USD"}

    err := amountA.Add(amountB)
    if err != nil {
        t.Errorf("expected no error, got: %v", err)
    }

    if amountA.Value != 300.40 {
        t.Errorf("expected sum to be 300, got: %f", amountA.Value)
    }

    // Invalid Add (different Currency)
    amountC := db.Amount{Value: 100.20, Currency: "EUR"}
    err = amountA.Add(amountC)
    if err == nil {
        t.Error("expected error (different Currency), got no error")
    }

    // Invalid Add (>= custom max float (define in /config/config.go))
    valueProvokeMaxFloat64 := config.MaxFloat - amountA.Value + 0.01
    fmt.Println(valueProvokeMaxFloat64)
    amountG := db.Amount{Value: valueProvokeMaxFloat64, Currency: "USD"}
    err = amountA.Add(amountG)
    if err == nil {
        t.Errorf("expected error (>= custom max float (define in /config/config.go)), got no error with: %v", amountA.Value)
    }
}

func TestAmountSum(t *testing.T) {
    // Invalid Sum (>= custom max float (define in /config/config.go))

}
