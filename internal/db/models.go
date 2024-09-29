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


// List of models: Client, Session, CarTrip, ExpenseType, Expense
// Plus (not its own DB table): Amount, AmountList
// By order of less strict to more strict for validation:
// PreinsertValid < Valid < PreReportValid

// Client
// Methods: String, PreInsertValid, Valid
type Client struct {
    ID   int
    Name string
}

func (c Client) String() string {
    return fmt.Sprintf(c.Name)
}

func(c Client) PreInsertValid() error {
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
    ID       int
    ClientID int
    Location string
    TripStartLocation string
    TripEndLocation string
    StartAtDateTime  time.Time
    EndAtDateTime    time.Time
}

func (s Session) String() string {
    var format string
    if s.TripStartLocation != "" {
        format += s.TripStartLocation + " > "
    }
    format += fmt.Sprintf("[%v]", s.Location)
    if s.TripEndLocation!= "" {
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
    if s.ClientID == 0 || s.Location == "" {
        return utils.LogError("client ID and location cannot be empty")
    }
    if s.StartAtDateTime.After(s.EndAtDateTime) {
        return utils.LogError("start date must be before end date")
    }
    if (len([]rune(s.Location)) > 100 ||
        len([]rune(s.TripStartLocation)) > 100 ||
        len([]rune(s.TripEndLocation)) > 100) {
        return utils.LogError(
            "location, trip start location, and trip end location " +
            "cannot exceed maximum length of 100 characters")
    }
    return nil
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
// Methods:
type CarTrip struct {
    ID int
    SessionID int
    DistanceKM float64
    DateTime time.Time // TODO in crud: UNIQUE validation
}

func (ct CarTrip) String() string {
    var format string
    if ct.SessionID != 0 {
        format += fmt.Sprintf("Session ID: %d - ", ct.SessionID)
    }
    format += fmt.Sprintf("%v km @ %v", ct.DistanceKM, ct.DateTime.Format(time.DateOnly))
    return format
}

func (ct CarTrip) PreInsertValid() error {
    if ct.DistanceKM == 0 || ct.DateTime.IsZero() {
        return utils.LogError("distance km and datetime must be non zero")
    }
    if ct.DistanceKM < 0 {
        return utils.LogError("distance km must be positive")
    }
    return nil
}


// Amount (Not in DB)
// Methods: String, Valid, Add
type Amount struct {
    Value    float64
    Currency string
}

func (a Amount) String() string {
    return fmt.Sprintf("{%.2f %s}", a.Value, a.Currency)
}

func (a Amount) Valid() error {
    if a.Currency == "" || a.Value == 0 {
        return utils.LogError("amount value and currency must be non-empty and non-zero")
    }
    if a.Value < 0 {
        return utils.LogError("amount value must be positive")
    }
    return nil
}

func (a *Amount) Add(other Amount) error {
    errA := a.Valid()
    errOther := other.Valid()

    switch {
    case errA != nil:
        return errA
    case errOther != nil:
        return errOther
    case a.Currency != other.Currency:
        return utils.LogError("currencies don't match: %v and %v", a.Currency, other.Currency)
    case a.Value + other.Value > config.MaxFloat:
        return utils.LogError("sum exceeds maximum float64 value")
    default:
        a.Value += other.Value
        return nil
    }
}


// AmountList (Not in DB)
// Methods: String, Valid, Equal, Sum
type AmountList []Amount

func (a AmountList) String() string {
    var result string
    for _, amount := range a {
        result += fmt.Sprintf("%s, ", amount.String())
    }
    return result[:len(result) - 2]
}

func (a AmountList) Valid() error {
    if len(a) == 0 {
        return utils.LogError("list of amounts is empty")
    }

    for _, amount := range a {
        if err := amount.Valid(); err != nil {
            return err
        }
    }
    return nil
}

func (a AmountList) Equal(other AmountList) (bool, error) {
    if err := a.Valid(); err != nil {
        return false, err
    }
    if err := other.Valid(); err != nil {
        return false, err
    }
    if len(a) != len(other) {
        return false, nil
    }

    temp := make(AmountList, len(other))
    copy(temp, other)

    for _, amountA := range a {
        found := false
        for i, v := range temp {
            if v == amountA {
                temp = append(temp[:i], temp[i + 1:]...)
                found = true
                break
            }
        }

        if !found {
            return false, nil
        }
    }
    return true, nil
}

func (a AmountList) Sum() (AmountList, error) {
    if err := a.Valid(); err != nil {
        return nil, err
    }

    sumsByCurrency := make(map[string]float64)
    for _, amount := range a {
        if  sum := sumsByCurrency[amount.Currency] + amount.Value;
            sum > config.MaxFloat {
            return nil, utils.LogError(
                "sum exceeds maximum float64 value: %f + %f",
                sumsByCurrency[amount.Currency], amount.Value,
            )
        }
        sumsByCurrency[amount.Currency] += amount.Value
    }

    result := make([]Amount, 0, len(sumsByCurrency))
    for currency, sum := range sumsByCurrency {
        result = append(result, Amount{Value: sum, Currency: currency})
    }
    return result, nil
}


// ExpenseType
// Methods: String, PreInsertValid, Valid
type ExpenseType struct {
    ID   int
    Name string
    TaxeRate float64
}

func (et ExpenseType) String() string {
    return fmt.Sprintf("%v taxed at %v%%", et.Name, et.TaxeRate)
}

func (et ExpenseType) PreInsertValid() error {
    if et.Name == "" {
        return utils.LogError("name must be non-zero")
    }
    if et.TaxeRate < 0 || et.TaxeRate > 60 {
        return utils.LogError(
            "taxe rate must be positive and less than or equal to 60",
        )
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
// Methods: String, PreInsertValid, CheckReceipt
type Expense struct {
    ID         int
    SessionID  int
    TypeID     int
    Amount     Amount       // Mapped data: will be 2 column in DB
    Location   string
    DateTime   time.Time
    ReceiptURL string
    Notes      string
}

func (e Expense) String() string {
    format := fmt.Sprintf(
        "%v (%d)\nat %v (%v)\nreceipt: %v",
        e.Amount, e.TypeID, e.Location, e.DateTime, e.ReceiptURL,
    )
    if e.Notes != "" {
        format += fmt.Sprintf("\nNotes: %v", e.Notes)
    }
    return format
}

func (e Expense) PreInsertValid() error {
    err := e.Amount.Valid()
    switch {
    case e.DateTime.IsZero() || e.ReceiptURL == "":
        return utils.LogError(
            "date and time, and receipt URL must be positive and non-zero",
        )
    case e.SessionID < 0:
        return utils.LogError("session ID must be 0 or positive.")
    case len([]rune(e.Notes)) >= 150:
        return utils.LogError("notes must be under 150 characters")
    case err != nil:
        return  err
    case e.ID < 0 || e.SessionID < 0:
        return utils.LogError("expense ID and session ID can't be negative")
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

func (e Expense) CheckReceipt() error {
    receiptPath := filepath.Join(config.GetAppPath(), e.ReceiptURL)
    _, errOs := os.Stat(receiptPath)
    errIsImageFile := isImageFile(receiptPath)

    switch {
    case e.ReceiptURL == config.DefaultReceiptURL:
        return utils.LogError("receipt is the default placeholder image")
    case e.ReceiptURL == "" :
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
            "invalid receipt image type %s: %s",filePath, contentType,
        )
    }
}
