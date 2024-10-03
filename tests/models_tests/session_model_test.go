package models_tests

import (
	"database/sql"
	"testing"
	"time"

	"github.com/craftidev/expenseflow/internal/db"
	"github.com/craftidev/expenseflow/tests"
)

func GetValidSession() db.Session {
	return db.Session{
		ID:                1,
		ClientID:          1,
		Location:          "New York City (NY)",
		TripStartLocation: sql.NullString{
            String: "Philadelphia (PA)",
            Valid: true,
        },
		TripEndLocation: sql.NullString{
            String: "Boston (MA)",
            Valid: true,
        },
		StartAtDateTime:   db.NullableTime{
            Time: time.Now(),
            Valid: true,
        },
		EndAtDateTime:   db.NullableTime{
            Time: time.Now().Add(24 * time.Hour),
            Valid: true,
        },
    }
}

func TestSessionPreInsertValid(t *testing.T) {
	validSession := GetValidSession()

    validSessions := tests.InitializeSliceOfValidAny(6, validSession)
    validSessions[1].ID = 0
    validSessions[2].TripStartLocation = sql.NullString{Valid: false}
    validSessions[3].TripEndLocation = sql.NullString{Valid: false}
    validSessions[4].StartAtDateTime = db.NullableTime{Valid:false}
    validSessions[5].EndAtDateTime = db.NullableTime{Valid:false}
    tests.ValidateEntities(t, validSessions, false, func(s db.Session) error {
        return s.PreInsertValid()
    })

	invalidSessions := tests.InitializeSliceOfValidAny(11, validSession)
	invalidSessions[0].ClientID = 0
	invalidSessions[1].ClientID = -1
	invalidSessions[2].Location = ""
	invalidSessions[3].Location = "a" + string(make([]rune, 100))
	invalidSessions[4].TripStartLocation.String =  ""
	invalidSessions[5].TripStartLocation.String =  "a" + string(make([]rune, 100))
	invalidSessions[6].TripEndLocation.String =  ""
	invalidSessions[7].TripEndLocation.String =  "a" + string(make([]rune, 100))
	invalidSessions[8].StartAtDateTime.Time =  time.Time{}
	invalidSessions[9].EndAtDateTime.Time =  time.Time{}
	// Test with start time after end time
	invalidSessions[10].StartAtDateTime.Time = time.Now().Add(24 * time.Hour)
	invalidSessions[10].EndAtDateTime.Time = time.Now()
	tests.ValidateEntities(t, invalidSessions, true, func(s db.Session) error {
		return s.PreInsertValid()
	})
}

func TestSessionValid(t *testing.T) {
	validSession := GetValidSession()

    validSessions := tests.InitializeSliceOfValidAny(5, validSession)
    validSessions[1].TripStartLocation.Valid = false
    validSessions[2].TripEndLocation.Valid = false
    validSessions[3].StartAtDateTime.Valid = false
    validSessions[4].EndAtDateTime.Valid = false
    tests.ValidateEntities(t, validSessions, false, func(s db.Session) error {
        return s.Valid()
    })

	invalidSessions := tests.InitializeSliceOfValidAny(2, validSession)
	invalidSessions[0].ID = 0
	invalidSessions[1].ID = -1
    tests.ValidateEntities(t, invalidSessions, true, func(s db.Session) error {
        return s.Valid()
    })
}

func TestSesionPreReportValid(t *testing.T) {
    validSession := GetValidSession()

    validSessions := tests.InitializeSliceOfValidAny(3, validSession)
    validSessions[1].TripStartLocation.Valid = false
    validSessions[2].TripEndLocation.Valid = false
    tests.ValidateEntities(t, validSessions, false, func(s db.Session) error {
        return s.PreReportValid()
    })

    invalidSessions := tests.InitializeSliceOfValidAny(2, validSession)
    invalidSessions[0].StartAtDateTime.Valid = false
    invalidSessions[1].EndAtDateTime.Valid = false
    tests.ValidateEntities(t, invalidSessions, true, func(s db.Session) error {
        return s.PreReportValid()
    })
}
