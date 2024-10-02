package crud

import (
	"database/sql"
	"log"

	"github.com/craftidev/expenseflow/internal/db"
	"github.com/craftidev/expenseflow/internal/utils"
)

func CreateLineItem(database *sql.DB, lineItem db.LineItem) (int64, error) {
	if err := lineItem.PreInsertValid(); err != nil {
		return 0, err
	}

	sqlQuery := `INSERT INTO clients(
                    expense_id,
                    taxe_rate,
                    total
                ) VALUES (?, ?, ?)`
	stmt, err := database.Prepare(sqlQuery)
	if err != nil {
		return 0, utils.LogError(
			"rejected querry: %v, error: %v", sqlQuery, err,
		)
	}
	defer stmt.Close()

	res, err := stmt.Exec(
        lineItem.ExpenseID,
        lineItem.TaxeRate,
        lineItem.Total,
    )
	if err != nil {
		return 0, utils.LogError(
			"unable to create line item: %v, error: %v",
			lineItem, err,
		)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, utils.LogError(
			"new line item created, but failed to get last inserted ID: %v, error: %v",
			lineItem, err,
		)
	}

	log.Printf("[info] new line item (ID: %v) created", id)
	return id, nil
}

func GetLineItemByID(database *sql.DB, id int64) (*db.LineItem, error) {
	sqlQuery := "SELECT id, name FROM line_items WHERE id = ?"
	stmt, err := database.Prepare(sqlQuery)
	if err != nil {
		return nil, utils.LogError(
			"rejected querry: %v, error: %v", sqlQuery, err,
		)
	}
	defer stmt.Close()

	var lineItem db.LineItem
	err = stmt.QueryRow(id).Scan(
        &lineItem.ID,
        &lineItem.ExpenseID,
        &lineItem.TaxeRate,
        &lineItem.Total,
    )
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.LogError("line item not found (ID: %d)", id)
		}
		return nil, utils.LogError("failed to fetch line item by ID: %v", err)
	}

	if err := lineItem.Valid(); err != nil {
		return nil, err // Integrity of data is breached
	}
	return &lineItem, nil
}

func UpdateLineItem(database *sql.DB, lineItem db.LineItem) error {
	if err := lineItem.Valid(); err != nil {
		return err
	}

	sqlQuery := "UPDATE line_items SET " + `
                    expense_id = ?,
                    taxe_rate = ?,
                    total = ?` +
                " WHERE id = ?"
	stmt, err := database.Prepare(sqlQuery)
	if err != nil {
		return utils.LogError("rejected querry: %v, error: %v", sqlQuery, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(
        lineItem.ExpenseID,
        lineItem.TaxeRate,
        lineItem.Total,
        lineItem.ID,
    )
	if err != nil {
		return utils.LogError(
            "unable to update line item: %v, error: %v", lineItem, err,
        )
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return utils.LogError("failed to check affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return utils.LogError("no line item found with ID: %d", lineItem.ID)
	}

	log.Printf("[info] line item (ID: %v) updated", lineItem.ID)
	return nil
}

func DeleteLineItemByID(database *sql.DB, id int64) error {
	if id <= 0 {
		return utils.LogError("line item ID must be positive and non-zero")
	}

	sqlQuery := "DELETE FROM line_items WHERE id = ?"
	stmt, err := database.Prepare(sqlQuery)
	if err != nil {
		return utils.LogError("rejected querry: %v, error: %v", sqlQuery, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return utils.LogError(
            "unable to delete line item with ID: %v, error: %v", id, err,
        )
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return utils.LogError("failed to check affected rows: %v", err)
	}
	if rowsAffected == 0 {
		return utils.LogError("no line item found with ID: %d", id)
	}

	log.Printf("[info] line item (ID: %v) deleted", id)
	return nil
}
