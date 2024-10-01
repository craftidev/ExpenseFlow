package models_test

import (
	"testing"
	"time"

	"github.com/craftidev/expenseflow/internal/db"
	"github.com/craftidev/expenseflow/tests"
)

func getValidSession() db.Session {
	return db.Session{
		ID:                1,
		ClientID:          1,
		Location:          "New York City (NY)",
		TripStartLocation: "Philadelphia (PA)",
		TripEndLocation:   "Boston (MA)",
		StartAtDateTime:   time.Now(),
		EndAtDateTime:     time.Now().Add(24 * time.Hour),
	}
}

func TestSessionPreInsertValid(t *testing.T) {
	validSession := getValidSession()

    validSessions := tests.InitializeSliceOfValidAny(6, validSession)
    validSessions[1].ID = 0
    validSessions[2].TripStartLocation = ""
    validSessions[3].TripEndLocation = ""
    validSessions[4].StartAtDateTime = time.Time{}
    validSessions[5].EndAtDateTime = time.Time{}
    tests.ValidateEntities(t, validSessions, false, func(s db.Session) error {
        return s.PreInsertValid()
    })

	invalidSessions := tests.InitializeSliceOfValidAny(7, validSession)
	invalidSessions[0].ClientID = 0
	invalidSessions[1].ClientID = -1
	invalidSessions[2].Location = ""
	invalidSessions[3].Location = "a" + string(make([]rune, 100))
	invalidSessions[4].TripStartLocation = "a" + string(make([]rune, 100))
	invalidSessions[5].TripEndLocation = "a" + string(make([]rune, 100))
	// Test with start time after end time
	invalidSessions[6].StartAtDateTime = time.Now()
	invalidSessions[6].EndAtDateTime = time.Now().Add(-24 * time.Hour)
	tests.ValidateEntities(t, invalidSessions, true, func(s db.Session) error {
		return s.PreInsertValid()
	})
}

func TestSessionValid(t *testing.T) {
	validSession := getValidSession()

    validSessions := tests.InitializeSliceOfValidAny(5, validSession)
    validSessions[1].TripStartLocation = ""
    validSessions[2].TripEndLocation = ""
    validSessions[3].StartAtDateTime = time.Time{}
    validSessions[4].EndAtDateTime = time.Time{}
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
    validSession := getValidSession()

    validSessions := tests.InitializeSliceOfValidAny(3, validSession)
    validSessions[1].TripStartLocation = ""
    validSessions[2].TripEndLocation = ""
    tests.ValidateEntities(t, validSessions, false, func(s db.Session) error {
        return s.PreReportValid()
    })

    invalidSessions := tests.InitializeSliceOfValidAny(2, validSession)
    invalidSessions[0].StartAtDateTime = time.Time{}
    invalidSessions[1].EndAtDateTime = time.Time{}
    tests.ValidateEntities(t, invalidSessions, true, func(s db.Session) error {
        return s.PreReportValid()
    })
}
