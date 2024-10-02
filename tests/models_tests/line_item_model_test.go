package models_tests

import (
	"testing"

	"github.com/craftidev/expenseflow/config"
	"github.com/craftidev/expenseflow/internal/db"
	"github.com/craftidev/expenseflow/tests"
)

func getValidLineItem() db.LineItem {
    return db.LineItem{
        ID:   1,
        ExpenseID: 1,
        TaxeRate: 5.5,
        Total: 15.43,
    }
}

func TestLineItemPreInsertValid(t *testing.T) {
    validLineItem := getValidLineItem()

    validLineItems := tests.InitializeSliceOfValidAny(3, validLineItem)
    validLineItems[1].ID = 0
    validLineItems[2].TaxeRate = 0
    tests.ValidateEntities(t, validLineItems, false, func(li db.LineItem) error {
        return li.PreInsertValid()
    })

    invalidLineItems := tests.InitializeSliceOfValidAny(8, validLineItem)
    invalidLineItems[0].ExpenseID = -1
    invalidLineItems[1].ExpenseID = 0
    invalidLineItems[2].TaxeRate = -1
    invalidLineItems[3].TaxeRate = 60.01
    invalidLineItems[4].TaxeRate = 0.01 + config.MaxFloat
    invalidLineItems[5].Total = -1
    invalidLineItems[6].Total = 0
    invalidLineItems[7].Total = 0.01 + config.MaxFloat
    tests.ValidateEntities(t, invalidLineItems, true, func(li db.LineItem) error {
        return li.PreInsertValid()
    })
}

func TestLineItemValid(t *testing.T) {
    validLineItem := getValidLineItem()

    invalidLineItems := tests.InitializeSliceOfValidAny(2, validLineItem)
    invalidLineItems[0].ID = -1
    invalidLineItems[1].ID = 0
    tests.ValidateEntities(t, invalidLineItems, true, func(li db.LineItem) error {
        return li.Valid()
    })
}
