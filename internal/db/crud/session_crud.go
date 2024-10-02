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
        log.Println("[error] New session created with unknown ID") // The error did not prevent the INSERT
        return 0, utils.LogError(
            "failed to get last inserted ID: %v, error: %v",
            session, err,
        )
    }
    log.Printf("[info] New session (ID: %v) created", id)

    return id, nil
}
