package crud

import (
	"database/sql"
	"log"

	"github.com/craftidev/expenseflow/internal/db"
	"github.com/craftidev/expenseflow/internal/utils"
)

func CreateExpenseType(database *sql.DB, expenseType db.ExpenseType) (int64, error) {
	if err := expenseType.PreInsertValid(); err != nil {
		return 0, err
	}
	ok, err := expenseTypeNameIsUnique(database, expenseType)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, utils.LogError(
			"expense type name already exists: %s", expenseType.Name,
		)
	}

	sqlQuery := "INSERT INTO expense_types(name) VALUES (?)"
	stmt, err := database.Prepare(sqlQuery)
	if err != nil {
		return 0, utils.LogError(
			"rejected querry: %v, error: %v", sqlQuery, err,
		)
	}
	defer stmt.Close()

	res, err := stmt.Exec(expenseType.Name)
	if err != nil {
		return 0, utils.LogError(
			"unable to create expense type: %v, error: %v",
			expenseType, err,
		)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, utils.LogError(
			"new expense type created, but failed to get last inserted ID: %v, error: %v",
			expenseType, err,
		)
	}

	log.Printf("[info] new expense type (ID: %v) created", id)
	return id, nil
}

func GetExpenseTypeByID(database *sql.DB, id int64) (*db.ExpenseType, error) {
	sqlQuery := "SELECT id, name FROM expense_types WHERE id = ?"
	stmt, err := database.Prepare(sqlQuery)
	if err != nil {
		return nil, utils.LogError(
			"rejected querry: %v, error: %v", sqlQuery, err,
		)
	}
	defer stmt.Close()

	var expenseType db.ExpenseType
	err = stmt.QueryRow(id).Scan(&expenseType.ID, &expenseType.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.LogError("expense type not found (ID: %d)", id)
		}
		return nil, utils.LogError("failed to fetch expense type by ID: %v", err)
	}

	if err := expenseType.Valid(); err != nil {
		return nil, err // Integrity of data is breached
	}
	return &expenseType, nil
}

func UpdateExpenseType(database *sql.DB, expenseType db.ExpenseType) error {
	if err := expenseType.Valid(); err != nil {
		return err
	}
	ok, err := expenseTypeNameIsUnique(database, expenseType)
	if err != nil {
		return err
	}
	if !ok {
		return utils.LogError("expense type name already exists: %s",
        expenseType.Name,
    )
	}

	sqlQuery := "UPDATE expense_types SET name = ? WHERE id = ?"
	stmt, err := database.Prepare(sqlQuery)
	if err != nil {
		return utils.LogError("rejected querry: %v, error: %v", sqlQuery, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(expenseType.Name, expenseType.ID)
	if err != nil {
		return utils.LogError(
            "unable to update expense type: %v, error: %v", expenseType, err,
        )
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return utils.LogError("failed to check affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return utils.LogError(
            "no expense type found with ID: %d", expenseType.ID,
        )
	}

	log.Printf("[info] expense type (ID: %v) updated", expenseType.ID)
	return nil
}

func DeleteExpenseTypeByID(database *sql.DB, id int64) error {
	if id <= 0 {
		return utils.LogError("expense type ID must be positive and non-zero")
	}

	ok, err := expenseTypeIsNotRefAsAnFK(database, id)
	if err != nil {
		return err
	}
	if !ok {
		return utils.LogError(
			"expense type (ID: %v) is still referenced by expenses", id,
		)
	}

	sqlQuery := "DELETE FROM expense_types WHERE id = ?"
	stmt, err := database.Prepare(sqlQuery)
	if err != nil {
		return utils.LogError("rejected querry: %v, error: %v", sqlQuery, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return utils.LogError(
            "unable to delete expense type with ID: %v, error: %v", id, err,
        )
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return utils.LogError("failed to check affected rows: %v", err)
	}
	if rowsAffected == 0 {
		return utils.LogError("no expense type found with ID: %d", id)
	}

	log.Printf("[info] expense type (ID: %v) deleted", id)
	return nil
}

func expenseTypeNameIsUnique(database *sql.DB, expenseType db.ExpenseType) (bool, error) {
	sqlQuery := "SELECT COUNT(*) FROM expense_types WHERE name = ? AND id != ?"

	stmt, err := database.Prepare(sqlQuery)
	if err != nil {
		return false, utils.LogError(
            "rejected querry: %v, error: %v", sqlQuery, err,
        )
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(expenseType.Name, expenseType.ID).Scan(&count)
	if err != nil {
		return false, utils.LogError(
            "failed to count expense types with name: %v, error: %v",
            expenseType.Name, err,
        )
	}

	return count == 0, nil
}

func expenseTypeIsNotRefAsAnFK(database *sql.DB, id int64) (bool, error) {
	sqlQuery := "SELECT COUNT(*) FROM expenses WHERE type_id = ?"

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
			"failed to count expenses with expense type ID: %v, error: %v",
			id, err,
		)
	}

	return count == 0, nil
}
