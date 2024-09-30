package db

import (
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
	ID   int
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
	ID                int
	ClientID          int
	Location          string
	TripStartLocation string
	TripEndLocation   string
	StartAtDateTime   time.Time
	EndAtDateTime     time.Time
}

func (s Session) String() string {
	var format string
	if s.TripStartLocation != "" {
		format += s.TripStartLocation + " > "
	}
	format += fmt.Sprintf("[%v]", s.Location)
	if s.TripEndLocation != "" {
		format += " > " + s.TripEndLocation
	}
	format += "\n[ "
	if s.StartAtDateTime.IsZero() {
		format += s.StartAtDateTime.Format(time.DateOnly)
	}
	format += " - "
	if s.EndAtDateTime.IsZero() {
		format += s.EndAtDateTime.Format(time.DateOnly)
	}
	format += " ]"

	return fmt.Sprint(format)
}

func (s Session) PreInsertValid() error {
    switch {
    case s.ClientID <= 0 || s.Location == "":
		return utils.LogError("client ID and location cannot be empty or negative")
    case !s.EndAtDateTime.IsZero() && s.StartAtDateTime.After(s.EndAtDateTime):
		return utils.LogError("start date must be before end date")
    case    len([]rune(s.Location)) > 100 ||
            len([]rune(s.TripStartLocation)) > 100 ||
            len([]rune(s.TripEndLocation)) > 100:
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
	if s.StartAtDateTime.IsZero() || s.EndAtDateTime.IsZero() {
		return utils.LogError("start date and end date cannot be empty")
	}
	return s.Valid()
}

// CarTrip
// Methods: String, PreInsertValid, Valid
type CarTrip struct {
	ID         int
	SessionID  int
	DistanceKM float64
	DateOnly   string
	// TODO in crud: UNIQUE validation
}

func (ct CarTrip) String() string {
	var format string
	if ct.SessionID != 0 {
		format += fmt.Sprintf("Session ID: %d - ", ct.SessionID)
	}
	format += fmt.Sprintf("%v km @ %v", ct.DistanceKM, ct.DateOnly)
	return format
}

func (ct CarTrip) PreInsertValid() error {
	dateTimeFormat, err := time.Parse(time.DateOnly, ct.DateOnly)
	if err != nil {
		return utils.LogError("invalid date format, expected yyyy-mm-dd: %v", err)
	}

	switch {
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
	ID   int
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
	ID             int
	SessionID      int
	TypeID         int
	Currency       string
	ReceiptRelPath string
	Notes          string
	DateTime       time.Time
}

func (e Expense) String() string {
	format := fmt.Sprintf(
		"Type: %v (%v) @ %v",
		e.TypeID, e.Currency, e.DateTime.Format(time.DateOnly),
	)
	if e.ReceiptRelPath != "" {
		format += fmt.Sprintf("\n%v", e.ReceiptRelPath)
	}
	if e.Notes != "" {
		format += fmt.Sprintf("\nNotes: %v", e.Notes)
	}
	return format
}

func (e Expense) PreInsertValid() error {
	switch {
	case e.TypeID == 0 || e.Currency == "" || e.DateTime.IsZero():
		return utils.LogError(
			"type id, currency and date time  must be non-zero",
		)
	case e.SessionID < 0:
		return utils.LogError("session ID can't be negative")
	case e.TypeID <= 0:
		return utils.LogError("type ID must be positive and non-zero.")
	case len([]rune(e.Currency)) > 10:
		return utils.LogError("currency can't exceeds 10 characters")
	case len([]rune(e.ReceiptRelPath)) > 50:
		return utils.LogError("receipt URL can't exceeds 50 characters")
	case len([]rune(e.Notes)) > 150:
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
	receiptFullPath := filepath.Join(config.ReceiptsDir, e.ReceiptRelPath)
	_, errOs := os.Stat(receiptFullPath)
	errIsImageFile := isImageFile(receiptFullPath)

	switch {
	case e.ReceiptRelPath == "":
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
	ID        int
	ExpenseID int
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
