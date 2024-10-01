package models_test

import (
	"testing"

	"github.com/craftidev/expenseflow/internal/db"
	"github.com/craftidev/expenseflow/tests"
)

func getValidExpenseType() db.ExpenseType {
    return db.ExpenseType{
        ID:   1,
        Name: "Transportation",
    }
}

func TestExpenseTypePreInsertValid(t *testing.T) {
    validExpenseType := getValidExpenseType()

    validExpenseTypes := tests.InitializeSliceOfValidAny(2, validExpenseType)
    validExpenseTypes[1].ID = 0
    tests.ValidateEntities(t, validExpenseTypes, false, func(et db.ExpenseType) error {
        return et.PreInsertValid()
    })

    invalidExpenseTypes := tests.InitializeSliceOfValidAny(2, validExpenseType)
    invalidExpenseTypes[0].Name = ""
    invalidExpenseTypes[1].Name = "a" + string(make([]rune, 50))
    tests.ValidateEntities(t, invalidExpenseTypes, true, func(et db.ExpenseType) error {
        return et.PreInsertValid()
    })
}

func TestExpenseTypeValid(t *testing.T) {
    validExpenseType := getValidExpenseType()

    invalidExpenseTypes := tests.InitializeSliceOfValidAny(2, validExpenseType)
    invalidExpenseTypes[0].ID = -1
    invalidExpenseTypes[1].ID = 0
    tests.ValidateEntities(t, invalidExpenseTypes, true, func(et db.ExpenseType) error {
        return et.Valid()
    })
}
