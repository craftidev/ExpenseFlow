package crud

import (
	"database/sql"
	"log"

	"github.com/craftidev/expenseflow/internal/db"
	"github.com/craftidev/expenseflow/internal/utils"
)

func CreateExpense(database *sql.DB, expense db.Expense) (int64, error) {
	if err := expense.PreInsertValid(); err != nil {
		return 0, err
	}

	sqlQuery := `INSERT INTO expenses(
                    session_id,
                    type_id,
                    currency,
                    receipt_rel_path,
                    notes,
                    date_time
                ) VALUES (?, ?, ?, ?, ?, ?)`
	stmt, err := database.Prepare(sqlQuery)
	if err != nil {
		return 0, utils.LogError(
			"rejected querry: %v, error: %v", sqlQuery, err,
		)
	}
	defer stmt.Close()

	res, err := stmt.Exec(
        expense.SessionID,
        expense.TypeID,
        expense.Currency,
        expense.ReceiptRelPath,
        expense.Notes,
        expense.DateTime,
    )
	if err != nil {
		return 0, utils.LogError(
			"unable to create expense: %v, error: %v",
			expense, err,
		)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, utils.LogError(
			"new expense created, but failed to get last inserted ID: %v, error: %v",
			expense, err,
		)
	}

	log.Printf("[info] new expense (ID: %v) created", id)
	return id, nil
}

func GetExpenseByID(database *sql.DB, id int64) (*db.Expense, error) {
	sqlQuery := `SELECT
                    id,
                    session_id,
                    type_id,
                    currency,
                    receipt_rel_path,
                    notes,
                    date_time
                FROM expenses WHERE id = ?`
	stmt, err := database.Prepare(sqlQuery)
	if err != nil {
		return nil, utils.LogError(
			"rejected querry: %v, error: %v", sqlQuery, err,
		)
	}
	defer stmt.Close()

	var expense db.Expense
    var dateTime string
	err = stmt.QueryRow(id).Scan(
        &expense.ID,
        &expense.SessionID,
        &expense.TypeID,
        &expense.Currency,
        &expense.ReceiptRelPath,
        &expense.Notes,
        &dateTime,
    )
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.LogError("expense not found (ID: %d)", id)
		}
		return nil, utils.LogError("failed to fetch expense by ID: %v", err)
	}

    expense.DateTime, err = ParsingStrToTime(dateTime)
    if err!= nil {
        return nil, err
    }

	if err := expense.Valid(); err != nil {
		return nil, err // Integrity of data is breached
	}
	return &expense, nil
}

func UpdateExpense(database *sql.DB, expense db.Expense) error {
	if err := expense.Valid(); err != nil {
		return err
	}

	sqlQuery := `UPDATE expenses SET
                    session_id,
                    type_id,
                    currency,
                    receipt_rel_path,
                    notes,
                    date_time
                WHERE id = ?`
	stmt, err := database.Prepare(sqlQuery)
	if err != nil {
		return utils.LogError("rejected querry: %v, error: %v", sqlQuery, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(
        expense.SessionID,
        expense.TypeID,
        expense.Currency,
        expense.ReceiptRelPath,
        expense.Notes,
        expense.DateTime,
        expense.ID,
    )
	if err != nil {
		return utils.LogError(
            "unable to update expense: %v, error: %v", expense, err,
        )
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return utils.LogError("failed to check affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return utils.LogError("no expense found with ID: %d", expense.ID)
	}

	log.Printf("[info] expense (ID: %v) updated", expense.ID)
	return nil
}

func DeleteExpenseByID(database *sql.DB, id int64) error {
	if id <= 0 {
		return utils.LogError("expense ID must be positive and non-zero")
	}

	ok, err := expenseIsNotRefAsAnFK(database, id)
	if err != nil {
		return err
	}
	if !ok {
		return utils.LogError(
			"expense (ID: %v) is still referenced by line items", id,
		)
	}

	sqlQuery := "DELETE FROM expenses WHERE id = ?"
	stmt, err := database.Prepare(sqlQuery)
	if err != nil {
		return utils.LogError("rejected querry: %v, error: %v", sqlQuery, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return utils.LogError(
            "unable to delete expense with ID: %v, error: %v", id, err,
        )
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return utils.LogError("failed to check affected rows: %v", err)
	}
	if rowsAffected == 0 {
		return utils.LogError("no expense found with ID: %d", id)
	}

	log.Printf("[info] expense (ID: %v) deleted", id)
	return nil
}

func expenseIsNotRefAsAnFK(database *sql.DB, id int64) (bool, error) {
	sqlQuery := "SELECT COUNT(*) FROM line_items WHERE expense_id = ?"

	stmt, err := database.Prepare(sqlQuery)
	if err != nil {
		return false, utils.LogError(
			"rejected querry: %v, error: %v", sqlQuery, err,
		)
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(id).Scan(&count)
	if err != nil {
		return false, utils.LogError(
			"failed to count line items with expense ID: %v, error: %v",
			id, err,
		)
	}

	return count == 0, nil
}
