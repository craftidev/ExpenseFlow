package models_test

import (
	"testing"
	"time"

	"github.com/craftidev/expenseflow/internal/db"
	"github.com/craftidev/expenseflow/tests"
)

func getValidCarTrip() db.CarTrip {
    return db.CarTrip{
        ID:         1,
        SessionID:  1,
        DistanceKM: 50.5,
        DateOnly:  "2022-01-01",
    }
}

func TestCarTripPreInsertValid(t *testing.T) {
    validCarTrip := getValidCarTrip()

    validCarTrips := tests.InitializeSliceOfValidAny(3, validCarTrip)
    validCarTrips[1].ID = 0
    validCarTrips[2].SessionID = 0
    tests.ValidateEntities(t, validCarTrips, false, func(ct db.CarTrip) error {
        return ct.PreInsertValid()
    })

    invalidCarTrips := tests.InitializeSliceOfValidAny(6, validCarTrip)
    invalidCarTrips[0].SessionID = -1
    invalidCarTrips[1].DistanceKM = -1
    invalidCarTrips[2].DistanceKM = 0
    invalidCarTrips[3].DateOnly = time.Time{}.Format(time.DateOnly)
    invalidCarTrips[4].DateOnly = "non sens"
    invalidCarTrips[5].DateOnly = "nons-en-s2"
    tests.ValidateEntities(t, invalidCarTrips, true, func(ct db.CarTrip) error {
        return ct.PreInsertValid()
    })
}

func TestCarTripValid(t *testing.T) {
    validCarTrip := getValidCarTrip()

    validCarTrips := tests.InitializeSliceOfValidAny(2, validCarTrip)
    validCarTrips[1].SessionID = 0
    tests.ValidateEntities(t, validCarTrips, false, func(ct db.CarTrip) error {
        return ct.Valid()
    })

    invalidCarTrips := tests.InitializeSliceOfValidAny(2, validCarTrip)
    invalidCarTrips[0].ID = -1
    invalidCarTrips[1].ID = 0
    tests.ValidateEntities(t, invalidCarTrips, true, func(ct db.CarTrip) error {
        return ct.Valid()
    })
}
