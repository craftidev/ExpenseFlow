package crud

import (
	"database/sql"
	"log"

	"github.com/craftidev/expenseflow/internal/db"
	"github.com/craftidev/expenseflow/internal/utils"
)


func CreateSession(database *sql.DB, session db.Session) (int64, error) {
    if err := session.PreInsertValid(); err != nil {
        return 0, err
    }

    sqlQuery :=
        `INSERT INTO sessions(
            client_id,
            location,
            trip_start_location,
            trip_end_location,
            start_at_date_time
            end_at_date_time
        ) VALUES (?, ?, ?, ?, ?, ?)`
    stmt, err := database.Prepare(sqlQuery)
    if err != nil {
        return 0, utils.LogError(
            "rejected query: %v, error: %v", sqlQuery, err,
        )
    }
    defer stmt.Close()

    res, err := stmt.Exec(
        session.ClientID,
        session.Location,
        session.TripStartLocation,
        session.TripEndLocation,
        session.StartAtDateTime,
        session.EndAtDateTime,
    )
    if err != nil {
        return 0, utils.LogError(
            "unable to create session: %v, error: %v",
            session, err,
        )
    }

    id, err := res.LastInsertId()
    if err != nil {
        log.Println("[error] new session created with unknown ID") // The error did not prevent the INSERT
        return 0, utils.LogError(
            "failed to get last inserted ID: %v, error: %v",
            session, err,
        )
    }
    log.Printf("[info] new session (ID: %v) created", id)

    return id, nil
}

func GetSessionByID(database *sql.DB, id int64) (*db.Session, error) {
    sqlQuery := `SELECT
                    id,
                    client_id,
                    location,
                    trip_start_location,
                    trip_end_location,
                    start_at_date_time,
                    end_at_date_time
                FROM sessions WHERE id = ?`
    stmt, err := database.Prepare(sqlQuery)
    if err != nil {
		return nil, utils.LogError(
            "rejected querry: %v, error: %v", sqlQuery, err,
        )
    }
    defer stmt.Close()

    var session db.Session
    err = stmt.QueryRow(id).Scan(
        &session.ID,
        &session.ClientID,
        &session.Location,
        &session.TripStartLocation,
        &session.TripEndLocation,
        &session.StartAtDateTime,
        &session.EndAtDateTime,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, utils.LogError("session not found (ID: %d)", id)
        }
        return nil, utils.LogError("failed to fetch session by ID: %v", err)
    }

    if err := session.Valid(); err != nil {
		return nil, err // Integrity of data is breached
	}
	return &session, nil
}

func UpdateSession(database *sql.DB, session db.Session) error {
	if err := session.Valid(); err != nil {
		return err
	}

	sqlQuery := `UPDATE sessions SET
                    client_id = ?,
                    location = ?,
                    trip_start_location = ?,
                    trip_end_location = ?,
                    start_at_date_time = ?,
                    end_at_date_time = ?
                WHERE id = ?`
	stmt, err := database.Prepare(sqlQuery)
	if err != nil {
		return utils.LogError("rejected querry: %v, error: %v", sqlQuery, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(
        session.ClientID,
        session.Location,
        session.TripStartLocation,
        session.TripEndLocation,
        session.StartAtDateTime,
        session.EndAtDateTime,
        session.ID,
    )
	if err != nil {
		return utils.LogError("unable to update session: %v, error: %v", session, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return utils.LogError("failed to check affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return utils.LogError("no session found with ID: %d", session.ID)
	}

	log.Printf("[info] session (ID: %v) updated", session.ID)
	return nil
}

// TODO restrict ON DELETE for CarTrip and Expense using SessionID as FK
func DeleteSessionByID(database *sql.DB, id int64) error {
	if id <= 0 {
		return utils.LogError("session ID must be positive and non-zero")
	}

	sqlQuery := "DELETE FROM sessions WHERE id = ?"
	stmt, err := database.Prepare(sqlQuery)
	if err != nil {
		return utils.LogError("rejected querry: %v, error: %v", sqlQuery, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return utils.LogError("unable to delete session with ID: %v, error: %v", id, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return utils.LogError("failed to check affected rows: %v", err)
	}
	if rowsAffected == 0 {
		return utils.LogError("no session found with ID: %d", id)
	}

	log.Printf("[info] session (ID: %v) deleted", id)
	return nil
}
