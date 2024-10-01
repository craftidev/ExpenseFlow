package models_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/craftidev/expenseflow/config"
	"github.com/craftidev/expenseflow/internal/db"
	"github.com/craftidev/expenseflow/tests"
)

// Return the most restricted valid Expense
func getValidExpense() db.Expense {
	return db.Expense{
		ID:             1,
		SessionID:      1,
		TypeID:         1,
		Currency:       "USD",
		ReceiptRelPath: "valid_receipt_test.png",
		Notes:          "Valid notes here.",
		DateTime:       time.Now(),
	}
}

func TestExpensePreInsertValid(t *testing.T) {
	validExpense := getValidExpense()
	err := validExpense.PreInsertValid()
	if err != nil {
		t.Errorf("expected valid expense, got error: %v", err)
	}

	// Nullable  SessionID
	validExpense.SessionID = 0
	err = validExpense.PreInsertValid()
	if err != nil {
		t.Errorf("expected valid expense with zero-valued SessionID, got error: %v", err)
	}

	// Nullable ReceiptRelPath
	validExpense.ReceiptRelPath = ""
	err = validExpense.PreInsertValid()
	if err != nil {
		t.Errorf("expected valid expense with zero-valued ReceiptRelPath, got error: %v", err)
	}

	// Nullable Notes
	validExpense.Notes = ""
	err = validExpense.PreInsertValid()
	if err != nil {
		t.Errorf("expected valid expense with zero-valued Notes, got error: %v", err)
	}

	invalidExpenses := tests.InitializeSliceOfValidAny(8, validExpense)
	invalidExpenses[0].SessionID = -1
	invalidExpenses[1].TypeID = 0
	invalidExpenses[2].TypeID = -1
	invalidExpenses[3].Currency = ""
	invalidExpenses[4].Currency = "a" + string(make([]rune, 10))
	invalidExpenses[5].ReceiptRelPath = "a" + string(make([]rune, 50))
	invalidExpenses[6].Notes = "a" + string(make([]rune, 150))
	invalidExpenses[7].DateTime = time.Time{}
	tests.ValidateEntities(t, invalidExpenses, true, func(e db.Expense) error {
		return e.PreInsertValid()
	})
}

// Don't re-test what's already tested in PreInsertValid
func TestValid(t *testing.T) {
	// Valid Expense
	validExpense := getValidExpense()

	invalidExpenses := tests.InitializeSliceOfValidAny(2, validExpense)
	invalidExpenses[0].ID = -1
	invalidExpenses[1].ID = 0
	tests.ValidateEntities(t, invalidExpenses, true, func(e db.Expense) error {
		return e.Valid()
	})
}

// Don't re-test what's already tested in PreInsertValid or Valid
func TestPreReportValid(t *testing.T) {
	// Valid Expense
	validExpense := getValidExpense()

	invalidExpenses := tests.InitializeSliceOfValidAny(4, validExpense)
	invalidExpenses[0].ReceiptRelPath = "non_exitent_file.png"
	invalidExpenses[1].ReceiptRelPath = "invalid_receipt_test.txt"
	invalidExpenses[2].ReceiptRelPath = "corrupted_receipt_test.png"
	invalidExpenses[3].ReceiptRelPath = "protected_receipt_test.png"

	// Set permissions to 000 to create a protected file
	protectedFilePath := filepath.Join(config.ReceiptsDirTest, invalidExpenses[3].ReceiptRelPath)
	err := os.Chmod(protectedFilePath, 0000)
	if err != nil {
		t.Fatalf("failed to set file permissions before testing: %v", err)
	}

	tests.ValidateEntities(t, invalidExpenses, true, func(e db.Expense) error {
		return e.PreReportValid()
	})

	// Restore permissions after the test
	err = os.Chmod(protectedFilePath, 0644)
	if err != nil {
		t.Fatalf("failed to restore file permissions: %v", err)
	}
}
