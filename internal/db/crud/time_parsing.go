package crud

import (
	"database/sql"
	"strings"
	"time"

	"github.com/craftidev/expenseflow/internal/db"
	"github.com/craftidev/expenseflow/internal/utils"
)


func ParsingStrToTime(timestamp string) (time.Time, error) {
	// The full timestamp string looks like this with time.Time
    // "2024-10-03 14:28:59.052881354 +0200 CEST m=+0.004487003"
	parts := strings.Split(timestamp, " m=")
	cleanedTimestamp := parts[0]

    layout := "2006-01-02 15:04:05.999999999 -0700 MST"
	parsedTime, err := time.Parse(layout, cleanedTimestamp)
	if err != nil {
		return time.Time{}, utils.LogError("Error parsing time: %v", err)
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
