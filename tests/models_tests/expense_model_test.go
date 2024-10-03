package models_tests

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/craftidev/expenseflow/config"
	"github.com/craftidev/expenseflow/internal/db"
	"github.com/craftidev/expenseflow/tests"
)

// Return the most restricted valid Expense
func GetValidExpense() db.Expense {
	return db.Expense{
		ID:             1,
		SessionID:      sql.NullInt64{Int64: 1, Valid: true},
		TypeID:         1,
		Currency:       "USD",
		ReceiptRelPath: sql.NullString{
            String: "valid_receipt_test.png",
            Valid: true,
        },
		Notes:          sql.NullString{
            String: "Valid notes here.",
            Valid: true,
        },
		DateTime:       time.Now(),
	}
}

func TestExpensePreInsertValid(t *testing.T) {
	validExpense := GetValidExpense()
	err := validExpense.PreInsertValid()
	if err != nil {
		t.Errorf("expected valid expense, got error: %v", err)
	}

	// Nullable  SessionID
	validExpense.SessionID.Valid = false
	err = validExpense.PreInsertValid()
	if err != nil {
		t.Errorf("expected valid expense with zero-valued SessionID, got error: %v", err)
	}

	// Nullable ReceiptRelPath
	validExpense.ReceiptRelPath.Valid = false
	err = validExpense.PreInsertValid()
	if err != nil {
		t.Errorf("expected valid expense with zero-valued ReceiptRelPath, got error: %v", err)
	}

	// Nullable Notes
	validExpense.Notes.Valid = false
	err = validExpense.PreInsertValid()
	if err != nil {
		t.Errorf("expected valid expense with zero-valued Notes, got error: %v", err)
	}

    validExpense = GetValidExpense()
	invalidExpenses := tests.InitializeSliceOfValidAny(11, validExpense)
	invalidExpenses[0].SessionID.Int64 = -1
	invalidExpenses[1].SessionID.Int64 = 0
	invalidExpenses[2].TypeID = 0
	invalidExpenses[3].TypeID = -1
	invalidExpenses[4].Currency = ""
	invalidExpenses[5].Currency = "a" + string(make([]rune, 10))
	invalidExpenses[6].ReceiptRelPath.String = ""
	invalidExpenses[7].ReceiptRelPath.String = "a" + string(make([]rune, 50))
	invalidExpenses[8].Notes.String = ""
	invalidExpenses[9].Notes.String = "a" + string(make([]rune, 150))
	invalidExpenses[10].DateTime = time.Time{}
	tests.ValidateEntities(t, invalidExpenses, true, func(e db.Expense) error {
		return e.PreInsertValid()
	})
}

// Don't re-test what's already tested in PreInsertValid
func TestValid(t *testing.T) {
	// Valid Expense
	validExpense := GetValidExpense()

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
	validExpense := GetValidExpense()

	invalidExpenses := tests.InitializeSliceOfValidAny(4, validExpense)
	invalidExpenses[0].ReceiptRelPath.String = "non_exitent_file.png"
	invalidExpenses[1].ReceiptRelPath.String = "invalid_receipt_test.txt"
	invalidExpenses[2].ReceiptRelPath.String = "corrupted_receipt_test.png"
	invalidExpenses[3].ReceiptRelPath.String = "protected_receipt_test.png"

	// Set permissions to 000 to create a protected file
	protectedFilePath := filepath.Join(config.ReceiptsDirTest, invalidExpenses[3].ReceiptRelPath.String)
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
