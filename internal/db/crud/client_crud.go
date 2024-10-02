package crud

import (
	"database/sql"
	"log"

	"github.com/craftidev/expenseflow/internal/db"
	"github.com/craftidev/expenseflow/internal/utils"
)

func CreateClient(database *sql.DB, client db.Client) (int64, error) {
	if err := client.PreInsertValid(); err != nil {
		return 0, err
	}
	ok, err := checkClientNameIsUnique(database, client)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, utils.LogError(
            "client name already exists: %s", client.Name,
        )
	}

	sqlQuery := "INSERT INTO clients(name) VALUES (?)"
	stmt, err := database.Prepare(sqlQuery)
	if err != nil {
		return 0, utils.LogError(
            "rejected querry: %v, error: %v", sqlQuery, err,
        )
	}
	defer stmt.Close()

	res, err := stmt.Exec(client.Name)
	if err != nil {
		return 0, utils.LogError(
            "unable to create client: %v, error: %v",
            client, err,
        )
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println("[error] New client created with unknown ID") // The error did not prevent the INSERT
		return 0, utils.LogError(
            "failed to get last inserted ID: %v, error: %v",
            client, err,
        )
	}
	log.Printf("[info] New client (ID: %v) created", id)

	return id, nil
}

func GetClientByID(database *sql.DB, id int64) (*db.Client, error) {
	sqlQuery := "SELECT id, name FROM clients WHERE id = ?"
	stmt, err := database.Prepare(sqlQuery)
	if err != nil {
		return nil, utils.LogError(
            "rejected querry: %v, error: %v", sqlQuery, err,
        )
	}
	defer stmt.Close()

	var client db.Client
	err = stmt.QueryRow(id).Scan(&client.ID, &client.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.LogError("client not found (ID: %d)", id)
		}
		return nil, utils.LogError("failed to fetch client by ID: %v", err)
	}

	if err := client.Valid(); err != nil {
		return nil, err // Integrity of data is breached
	}
	return &client, nil
}

func UpdateClient(database *sql.DB, client db.Client) error {
	if err := client.Valid(); err != nil {
		return err
	}
	ok, err := checkClientNameIsUnique(database, client)
	if err != nil {
		return err
	}
	if !ok {
		return utils.LogError("client name already exists: %s", client.Name)
	}

	sqlQuery := "UPDATE clients SET name = ? WHERE id = ?"
	stmt, err := database.Prepare(sqlQuery)
	if err != nil {
		return utils.LogError("rejected querry: %v, error: %v", sqlQuery, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(client.Name, client.ID)
	if err != nil {
		return utils.LogError("unable to update client: %v, error: %v", client, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return utils.LogError("failed to check affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return utils.LogError("no client found with ID: %d", client.ID)
	}

	log.Printf("[info] Client (ID: %v) updated", client.ID)

	return nil
}

func DeleteClientByID(database *sql.DB, id int64) error {
	if id <= 0 {
		return utils.LogError("client ID must be positive and non-zero")
	}

	sqlQuery := "DELETE FROM clients WHERE id = ?"
	stmt, err := database.Prepare(sqlQuery)
	if err != nil {
		return utils.LogError("rejected querry: %v, error: %v", sqlQuery, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return utils.LogError("unable to delete client with ID: %v, error: %v", id, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return utils.LogError("failed to check affected rows: %v", err)
	}
	if rowsAffected == 0 {
		return utils.LogError("no client found with ID: %d", id)
	}

	log.Printf("[info] Client (ID: %v) deleted", id)

	return nil
}

func checkClientNameIsUnique(database *sql.DB, client db.Client) (bool, error) {
	sqlQuery := "SELECT COUNT(*) FROM clients WHERE name = ? AND id != ?"

	stmt, err := database.Prepare(sqlQuery)
	if err != nil {
		return false, utils.LogError("rejected querry: %v, error: %v", sqlQuery, err)
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(client.Name, client.ID).Scan(&count) // AutoIncrementation in sqlite make sure 0 never exist
	if err != nil {
		return false, utils.LogError("failed to count clients with name: %v, error: %v", client.Name, err)
	}

	return count == 0, nil
}
