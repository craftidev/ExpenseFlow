package crud

import (
	"database/sql"
	"log"

	"github.com/craftidev/expenseflow/internal/db"
	"github.com/craftidev/expenseflow/internal/utils"
)

func CreateCarTrip(database *sql.DB, carTrip db.CarTrip) (int64, error) {
	if err := carTrip.PreInsertValid(); err != nil {
		return 0, err
	}
	ok, err := carTripDateOnlyIsUnique(database, carTrip)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, utils.LogError(
			"car trip at this date already exists: %s", carTrip.DateOnly,
		)
	}

	sqlQuery := `INSERT INTO car_trips(
                    session_id,
                    distance_km,
                    date_only
                ) VALUES (?, ?, ?)`
	stmt, err := database.Prepare(sqlQuery)
	if err != nil {
		return 0, utils.LogError(
			"rejected querry: %v, error: %v", sqlQuery, err,
		)
	}
	defer stmt.Close()

	res, err := stmt.Exec(
        carTrip.SessionID,
        carTrip.DistanceKM,
        carTrip.DateOnly,
    )
	if err != nil {
		return 0, utils.LogError(
			"unable to create car trip: %v, error: %v",
			carTrip, err,
		)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, utils.LogError(
			"new car trip created, but failed to get last inserted ID: %v, error: %v",
			carTrip, err,
		)
	}

	log.Printf("[info] new car trip (ID: %v) created", id)
	return id, nil
}

func GetCarTripByID(database *sql.DB, id int64) (*db.CarTrip, error) {
	sqlQuery := "SELECT id, name FROM car_trips WHERE id = ?"
	stmt, err := database.Prepare(sqlQuery)
	if err != nil {
		return nil, utils.LogError(
			"rejected querry: %v, error: %v", sqlQuery, err,
		)
	}
	defer stmt.Close()

	var carTrip db.CarTrip
	err = stmt.QueryRow(id).Scan(
        &carTrip.ID,
        &carTrip.SessionID,
        &carTrip.DistanceKM,
        &carTrip.DateOnly,
    )
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.LogError("car trip not found (ID: %d)", id)
		}
		return nil, utils.LogError("failed to fetch car trip by ID: %v", err)
	}

	if err := carTrip.Valid(); err != nil {
		return nil, err // Integrity of data is breached
	}
	return &carTrip, nil
}

func UpdateCarTrip(database *sql.DB, carTrip db.CarTrip) error {
	if err := carTrip.Valid(); err != nil {
		return err
	}
	ok, err := carTripDateOnlyIsUnique(database, carTrip)
	if err != nil {
		return err
	}
	if !ok {
		return utils.LogError(
            "car trip at this date already exists: %s", carTrip.DateOnly,
        )
	}

	sqlQuery := "UPDATE car_trips SET " + `
                    session_id = ?,
                    distance_km = ?,
                    date_only = ?` +
                " WHERE id = ?"
	stmt, err := database.Prepare(sqlQuery)
	if err != nil {
		return utils.LogError("rejected querry: %v, error: %v", sqlQuery, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(
        carTrip.SessionID, carTrip.DistanceKM, carTrip.DateOnly, carTrip.ID,
    )
	if err != nil {
		return utils.LogError("unable to update car trip: %v, error: %v", carTrip, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return utils.LogError("failed to check affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return utils.LogError("no car trip found with ID: %d", carTrip.ID)
	}

	log.Printf("[info] car trip (ID: %v) updated", carTrip.ID)
	return nil
}

func DeleteCarTripByID(database *sql.DB, id int64) error {
	if id <= 0 {
		return utils.LogError("car trip ID must be positive and non-zero")
	}

	sqlQuery := "DELETE FROM car_trips WHERE id = ?"
	stmt, err := database.Prepare(sqlQuery)
	if err != nil {
		return utils.LogError("rejected querry: %v, error: %v", sqlQuery, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return utils.LogError(
            "unable to delete car trip with ID: %v, error: %v", id, err,
        )
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return utils.LogError("failed to check affected rows: %v", err)
	}
	if rowsAffected == 0 {
		return utils.LogError("no car trip found with ID: %d", id)
	}

	log.Printf("[info] car trip (ID: %v) deleted", id)
	return nil
}

func carTripDateOnlyIsUnique(database *sql.DB, carTrip db.CarTrip) (bool, error) {
	sqlQuery := "SELECT COUNT(*) FROM car_trips WHERE date_only = ? AND id != ?"

	stmt, err := database.Prepare(sqlQuery)
	if err != nil {
		return false, utils.LogError("rejected querry: %v, error: %v", sqlQuery, err)
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(carTrip.DateOnly, carTrip.ID).Scan(&count)
	if err != nil {
		return false, utils.LogError(
            "failed to count car trips with date: %v, error: %v",
            carTrip.DateOnly, err,
        )
	}

	return count == 0, nil
}
