package models_test

import (
	"testing"
	"time"

	"github.com/craftidev/expenseflow/internal/db"
	"github.com/craftidev/expenseflow/tests"
)


func getValidSession() db.Session {
    return db.Session{
        ID: 1,
        ClientID: 1,
        Location: "New York City (NY)",
        TripStartLocation: "Philadelphia (PA)",
        TripEndLocation: "Boston (MA)",
        StartAtDateTime: time.Now(),
    }
}

func TestSessionPreInsertValid(t *testing.T) {
    validSession := getValidSession()
    err := validSession.PreInsertValid()
    if err != nil {
        t.Errorf("expected valid session, got error: %v", err)
    }

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
    tests.ValidateInvalidEntities(t, invalidSessions, func(s db.Session) error {
        return s.PreInsertValid()
    })
}
