package crud

import (
	"database/sql"
	"time"

	"github.com/craftidev/expenseflow/internal/db"
	"github.com/craftidev/expenseflow/internal/utils"
)


func ParsingStrToTime(timestamp string) (time.Time, error) {
	// The full timestamp string looks like this with time.Time
    // "2024-10-03 16:57:00.013887899+02:00"
    layout := "2006-01-02 15:04:05.999999999-07:00"
	parsedTime, err := time.Parse(layout, timestamp)
	if err != nil {
		return time.Time{}, utils.LogError("error parsing time: %v", err)
	}

	return parsedTime, nil
}

func ParsingNullableStrToTime(timestamp sql.NullString) (
    db.NullableTime, error,
) {
    result := db.NullableTime{Valid: false}
    if timestamp.Valid {
        parsedTime, err := ParsingStrToTime(timestamp.String)
        if err != nil {
            return result, err
        }
        result.Time = parsedTime
        result.Valid = true
    }
    return result, nil
}
