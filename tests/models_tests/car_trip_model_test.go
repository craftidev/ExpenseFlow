package models_tests

import (
	"testing"
	"time"

	"github.com/craftidev/expenseflow/internal/db"
	"github.com/craftidev/expenseflow/tests"
)


func TestCarTripPreInsertValid(t *testing.T) {
    validCarTrip := tests.GetValidCarTrip()

    validCarTrips := tests.InitializeSliceOfValidAny(3, validCarTrip)
    validCarTrips[1].ID = 0
    validCarTrips[2].SessionID.Valid = false
    tests.ValidateEntities(t, validCarTrips, false, func(ct db.CarTrip) error {
        return ct.PreInsertValid()
    })

    invalidCarTrips := tests.InitializeSliceOfValidAny(7, validCarTrip)
    invalidCarTrips[0].SessionID.Int64 = 0
    invalidCarTrips[1].SessionID.Int64 = -1
    invalidCarTrips[2].DistanceKM = -1
    invalidCarTrips[3].DistanceKM = 0
    invalidCarTrips[4].DateOnly = time.Time{}.Format(time.DateOnly)
    invalidCarTrips[5].DateOnly = "non sens"
    invalidCarTrips[6].DateOnly = "nons-en-s2"
    tests.ValidateEntities(t, invalidCarTrips, true, func(ct db.CarTrip) error {
        return ct.PreInsertValid()
    })
}

func TestCarTripValid(t *testing.T) {
    validCarTrip := tests.GetValidCarTrip()

    validCarTrips := tests.InitializeSliceOfValidAny(2, validCarTrip)
    validCarTrips[1].SessionID.Valid = false
    tests.ValidateEntities(t, validCarTrips, false, func(ct db.CarTrip) error {
        return ct.Valid()
    })

    invalidCarTrips := tests.InitializeSliceOfValidAny(3, validCarTrip)
    invalidCarTrips[0].SessionID.Int64 = 0
    invalidCarTrips[1].ID = -1
    invalidCarTrips[2].ID = 0
    tests.ValidateEntities(t, invalidCarTrips, true, func(ct db.CarTrip) error {
        return ct.Valid()
    })
}
