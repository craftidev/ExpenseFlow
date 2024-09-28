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


// List of models: Client, Session, ExpenseType, Expense
// Plus (not its own DB table): Amount, AmountList

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
// Methods: String, PreInsertValid, Valid
type Session struct {
    ID       int
    ClientID int // TODO delete
    Name     string
    Address  string
    StartAt  time.Time
    EndAt    time.Time
}

func (s Session) String() string {
    return fmt.Sprintf(s.Name)
}

func (s Session) PreInsertValid() error {
    if s.ClientID == 0 || s.Name == "" || s.Address == "" || s.StartAt.IsZero() || s.EndAt.IsZero() {
        return utils.LogError("client ID, name, address, start and end times must be positive and non-zero")
    }
    if s.StartAt.After(s.EndAt) {
        return utils.LogError("start time must be before end time (session ID: %d)", s.ID)
    }
    return nil
}

func (s Session) Valid() error {
    if s.ID <= 0 {
        return utils.LogError("session ID must be positive and non-zero")
    }
    return s.PreInsertValid()
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
}

func (et ExpenseType) String() string {
    return fmt.Sprintf(et.Name)
}

func (et ExpenseType) PreInsertValid() error {
    if et.Name == "" {
        return utils.LogError("ame must be non-zero")
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
}

func (e Expense) String() string {
    return fmt.Sprintf(
        "%v (%d)\nat %v (%v)\nreceipt: %v",
        e.Amount, e.TypeID, e.Location, e.DateTime, e.ReceiptURL,
    )
}

func (e Expense) PreInsertValid() error {
    err := e.Amount.Valid()
    switch {
    case e.SessionID <= 0 || e.DateTime.IsZero() || e.ReceiptURL == "":
        return utils.LogError(
            "session ID, date and time, and receipt URL must be positive and non-zero",
        )
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
    receiptPath := filepath.Join(config.Path, e.ReceiptURL)
    _, errOs := os.Stat(receiptPath)
    errIsImageFile := isImageFile(config.Path + e.ReceiptURL)

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
