package db

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/craftidev/expenseflow/config"
	"github.com/craftidev/expenseflow/internal/utils"
)


// List of models: Client, Session, CarTrip, ExpenseType, Expense, LineItem
// Iterables: ExpenseList, LineItemList

// By order of less strict to more strict for validation:
// - PreInsertValid (no ID is ok for insert) <
// - Valid (need ID, zero value for NULLable column is ok) <
// - PreReportValid (zero value for certain NULLable is not ok)


// Client
// Methods: String, PreInsertValid, Valid
type Client struct {
	ID   int64
	Name string
}

func (c Client) String() string {
	return fmt.Sprintf(c.Name)
}

func (c Client) PreInsertValid() error {
	if c.Name == "" {
		return utils.LogError("name must be non-zero")
	}
	if len([]rune(c.Name)) > 100 {
		return utils.LogError("client name exceeds maximum length of 100 characters")
	}
	return nil
}

func (c Client) Valid() error {
	if c.ID <= 0 {
		return utils.LogError("client ID must be positive and non-zero")
	}
	return c.PreInsertValid()
}

// Session
// Methods: String, PreInsertValid, Valid, PreReportValid
type Session struct {
	ID                int64
	ClientID          int64
	Location          string
	TripStartLocation sql.NullString
	TripEndLocation   sql.NullString
	StartAtDateTime   NullableTime
	EndAtDateTime     NullableTime
}

func (s Session) String() string {
	var format string
	if s.TripStartLocation.Valid {
		format += s.TripStartLocation.String + " > "
	}
	format += fmt.Sprintf("[%v]", s.Location)
	if s.TripEndLocation.Valid {
		format += " > " + s.TripEndLocation.String
	}
	format += "\n[ "
	if s.StartAtDateTime.Valid {
		format += s.StartAtDateTime.Time.Format(time.DateOnly)
	}
	format += " - "
	if s.EndAtDateTime.Valid {
		format += s.EndAtDateTime.Time.Format(time.DateOnly)
	}
	format += " ]"

	return fmt.Sprint(format)
}

func (s Session) PreInsertValid() error {
    switch {
    case    s.ClientID <= 0 ||
            s.Location == "" ||
            s.TripStartLocation.Valid && s.TripStartLocation.String == "" ||
            s.TripEndLocation.Valid && s.TripEndLocation.String == "" ||
            s.StartAtDateTime.Valid && s.StartAtDateTime.Time.IsZero() ||
            s.EndAtDateTime.Valid && s.EndAtDateTime.Time.IsZero():
		return utils.LogError(`
            client ID, location, trip start/end location and
            start/end at date time, cannot be empty or negative
        `)
    case    s.EndAtDateTime.Valid &&
            s.StartAtDateTime.Valid &&
            s.StartAtDateTime.Time.After(s.EndAtDateTime.Time):
		return utils.LogError("start date must be before end date")
    case    len([]rune(s.Location)) > 100 ||
            (s.TripStartLocation.Valid && len([]rune(s.TripStartLocation.String)) > 100) ||
            (s.TripEndLocation.Valid &&len([]rune(s.TripEndLocation.String)) > 100):
		return utils.LogError(
			"location, trip start location, and trip end location " +
			"cannot exceed maximum length of 100 characters",
        )
    default:
        return nil
    }
}

func (s Session) Valid() error {
	if s.ID <= 0 {
		return utils.LogError("session ID must be positive and non-zero")
	}
	return s.PreInsertValid()
}

func (s Session) PreReportValid() error {
	if !s.StartAtDateTime.Valid || !s.EndAtDateTime.Valid {
		return utils.LogError("start date and end date cannot be empty")
	}
	return s.Valid()
}

// CarTrip
// Methods: String, PreInsertValid, Valid
type CarTrip struct {
	ID         int64
	SessionID  sql.NullInt64
	DistanceKM float64
	DateOnly   string
	// TODO in crud: UNIQUE validation
}

func (ct CarTrip) String() string {
	var format string
	if ct.SessionID.Valid {
		format += fmt.Sprintf("Session ID: %d - ", ct.SessionID.Int64)
	}
	format += fmt.Sprintf("%v km @ %v", ct.DistanceKM, ct.DateOnly)
	return format
}

func (ct CarTrip) PreInsertValid() error {
	dateTimeFormat, err := time.Parse(time.DateOnly, ct.DateOnly)
	if err != nil || len([]rune(ct.DateOnly)) != 10 {
		return utils.LogError(
            "invalid date format, expected yyyy-mm-dd, got: %v. Error: %v",
            err, ct.DateOnly,
        )
	}

	switch {
    case ct.SessionID.Valid && ct.SessionID.Int64 <= 0:
        return utils.LogError("session ID must be positive and non-zero")
	case ct.DistanceKM == 0 || dateTimeFormat.IsZero():
		return utils.LogError("distance km and datetime must be non zero")
	case ct.DistanceKM < 0:
		return utils.LogError("distance km must be positive")
	case ct.DistanceKM > config.MaxFloat:
		return utils.LogError(
			"distance km must be less than custom float limit: %v",
			config.MaxFloat,
		)
	default:
		return nil
	}
}

func (ct CarTrip) Valid() error {
	if ct.ID <= 0 {
		return utils.LogError("car trip ID must be positive and non-zero")
	}
	return ct.PreInsertValid()
}

// ExpenseType
// Methods: String, PreInsertValid, Valid
type ExpenseType struct {
	ID   int64
	Name string // TODO: check UNIQUE in crud
}

func (et ExpenseType) String() string {
	return fmt.Sprint(et.Name)
}

func (et ExpenseType) PreInsertValid() error {
	if et.Name == "" {
		return utils.LogError("name must be non-zero")
	}
	if len([]rune(et.Name)) > 50 {
		return utils.LogError("name cannot exceed maximum length of 50 characters")
	}
	return nil
}

func (et ExpenseType) Valid() error {
	if et.ID <= 0 {
		return utils.LogError("expense type ID must be positive and non-zero")
	}
	return et.PreInsertValid()
}

// Expense
// Methods: String, PreInsertValid, Valid, PreReportValid
type Expense struct {
	ID             int64
	SessionID      sql.NullInt64
	TypeID         int64
	Currency       string
	ReceiptRelPath sql.NullString
	Notes          sql.NullString
	DateTime       time.Time
}

func (e Expense) String() string {
	format := fmt.Sprintf(
		"Type: %v (%v) @ %v",
		e.TypeID, e.Currency, e.DateTime.Format(time.DateOnly),
	)
	if e.ReceiptRelPath.Valid {
		format += fmt.Sprintf("\n%v", e.ReceiptRelPath)
	}
	if e.Notes.Valid {
		format += fmt.Sprintf("\nNotes: %v", e.Notes)
	}
	return format
}

func (e Expense) PreInsertValid() error {
	switch {
	case    e.TypeID == 0 ||
            e.Currency == "" ||
            e.DateTime.IsZero() ||
            (e.ReceiptRelPath.Valid && e.ReceiptRelPath.String == "") ||
            (e.Notes.Valid && e.Notes.String == ""):
		return utils.LogError(
			"type id, currency, receipt url, notes and date time must be non-zero",
		)
	case e.SessionID.Valid && e.SessionID.Int64 <= 0:
		return utils.LogError("session ID must be positive and non-zero")
	case e.TypeID <= 0:
		return utils.LogError("type ID must be positive and non-zero.")
	case len([]rune(e.Currency)) > 10:
		return utils.LogError("currency can't exceeds 10 characters")
	case e.ReceiptRelPath.Valid && len([]rune(e.ReceiptRelPath.String)) > 50:
		return utils.LogError("receipt URL can't exceeds 50 characters")
	case e.Notes.Valid && len([]rune(e.Notes.String)) > 150:
		return utils.LogError("notes can't exceeds 150 characters")
	default:
		return nil
	}
}

func (e Expense) Valid() error {
	if e.ID <= 0 {
		return utils.LogError("expense ID must be positive and non-zero")
	}
	return e.PreInsertValid()
}

func (e Expense) PreReportValid() error {
	if err := e.checkReceipt(); err != nil {
		return err
	}
	return e.Valid()
}

func (e Expense) checkReceipt() error {
    if !e.ReceiptRelPath.Valid {
		return utils.LogError("receipt URL is empty")
    }
	receiptFullPath := filepath.Join(config.ReceiptsDir, e.ReceiptRelPath.String)
	_, errOs := os.Stat(receiptFullPath)
	errIsImageFile := isImageFile(receiptFullPath)

	switch {
    case e.ReceiptRelPath.String == "":
		return utils.LogError("receipt URL is empty")
	case errors.Is(errOs, os.ErrNotExist):
		return utils.LogError("invalid receipt URL: %v", errOs)
	case errOs != nil:
		return utils.LogError("undefined file error: %v", errOs)
	case errIsImageFile != nil:
		return errIsImageFile
	default:
		return nil
	}
}

func isImageFile(filePath string) error {
	receiptImage, err := os.Open(filePath)
	if err != nil {
		return utils.LogError("error opening receipt image: %v", err)
	}
	defer receiptImage.Close()

	// Read file header to determine content type
	buffer := make([]byte, 512)
	n, err := receiptImage.Read(buffer)
	if err != nil && err != io.EOF {
		return utils.LogError(
			"error reading headers of receipt image: %v", err,
		)
	}

	buffer = buffer[:n] // Adjust buffer size to the actual number of bytes read
	contentType := http.DetectContentType(buffer)
	switch contentType {
	case "image/jpeg", "image/png", "image/gif", "image/bmp", "image/webp":
		return nil
	default:
		return utils.LogError(
			"invalid receipt image type %s: %s", filePath, contentType,
		)
	}
}

// LineItem
// Method: String, PreInsertValid, Valid
type LineItem struct {
	ID        int64
	ExpenseID int64
	TaxeRate  float64
	Total     float64
}

func (li LineItem) String() string {
	return fmt.Sprintf(
		"Expense ID: %d - %.2f (taxe rate: %.2f%%)",
		li.ExpenseID, li.Total, li.TaxeRate*100,
	)
}

func (li LineItem) PreInsertValid() error {
	switch {
	case li.ExpenseID <= 0 || li.Total <= 0:
		return utils.LogError(
			"expense ID and total must be non-zero and  positive",
		)
	case li.TaxeRate < 0 || li.TaxeRate > 60:
		return utils.LogError(
			"taxe rate must be positive and not exceed 60",
		)
	case li.Total > config.MaxFloat:
		return utils.LogError(
			"total must not exceed maximum float64 value: %f",
			config.MaxFloat,
		)
	default:
		return nil
	}
}

func (li LineItem) Valid() error {
	if li.ID <= 0 {
		return utils.LogError("line item ID must be positive and non-zero")
	}
	return li.PreInsertValid()
}


// Iterables

// Method: MapExpensesByCurrency
type ExpenseList []Expense

// Method: SumByTaxeRates
type LineItemList []LineItem

func (eList ExpenseList) MapExpensesByCurrency() (map[string]ExpenseList, error) {
	result := make(map[string]ExpenseList)
	for _, expense := range eList {
		if err := expense.Valid(); err != nil {
			return nil, err
		}
		if _, ok := result[expense.Currency]; !ok {
			result[expense.Currency] = make(ExpenseList, 0)
		}
		result[expense.Currency] = append(result[expense.Currency], expense)
	}
	return result, nil
}

// This method consider Expense.Currency equality handled
func (liList LineItemList) SumByTaxeRates() (map[float64]float64, error) {
	result := make(map[float64]float64)
	for _, lineItem := range liList {
		if err := lineItem.Valid(); err != nil {
			return nil, err
		}
		result[lineItem.TaxeRate] += lineItem.Total
	}
	return result, nil
}

// Custom Nullable time.Time as the library sql doesn't have one
type NullableTime struct {
    Time time.Time
    Valid bool
}

func (nt *NullableTime) Scan(value interface{}) error {
    if value == nil {
        nt.Time, nt.Valid = time.Time{}, false
        return nil
    }
    nt.Valid = true
    switch v := value.(type) {
    case time.Time:
        nt.Time = v
        return nil
    case []byte:
        var err error
        nt.Time, err = time.Parse(time.RFC3339, string(v))
        return utils.LogError("Invalid time, error: %v", err)
    default:
        return utils.LogError("unable to scan NullableTime")
    }
}

func (nt NullableTime) Value() (driver.Value, error) {
    if !nt.Valid {
        return nil, nil
    }
    return nt.Time, nil
}

// TODO Equal func for my maps
// func (a AmountList) Equal(other AmountList) (bool, error) {
//     if err := a.Valid(); err != nil {
//         return false, err
//     }
//     if err := other.Valid(); err != nil {
//         return false, err
//     }
//     if len(a) != len(other) {
//         return false, nil
//     }

//     temp := make(AmountList, len(other))
//     copy(temp, other)

//     for _, amountA := range a {
//         found := false
//         for i, v := range temp {
//             if v == amountA {
//                 temp = append(temp[:i], temp[i + 1:]...)
//                 found = true
//                 break
//             }
//         }

//         if !found {
//             return false, nil
//         }
//     }
//     return true, nil
// }
