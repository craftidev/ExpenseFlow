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

// List of models: Client, Session, Amount, ExpenseType, Expense
// (Not its own DB table: Amount)

// Client
// Methods: String, Valid
type Client struct {
    ID   int
    Name string
}

func (c Client) String() string {
    return fmt.Sprintf(c.Name)
}

func(c Client) Valid() error {
    if c.ID == 0 || c.Name == "" {
        return utils.LogError("client ID and name must be non-zero and non-empty")
    }
    return nil
}

// Session
// Methods: String, Valid
type Session struct {
    ID       int
    ClientID int
    Name     string
    Address  string
    StartAt  time.Time
    EndAt    time.Time
}

func (s Session) String() string {
    return fmt.Sprintf(s.Name)
}

func (s Session) Valid() error {
    if s.ID == 0 || s.ClientID == 0 || s.Name == "" || s.Address == "" || s.StartAt.IsZero() || s.EndAt.IsZero() {
        return utils.LogError("session ID, client ID, name, address, start and end times must be non-zero and non-empty")
    }
    if s.StartAt.After(s.EndAt) {
        return utils.LogError("start time must be before end time (session ID: %d)", s.ID)
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
    return fmt.Sprintf("%.2f %s", a.Value, a.Currency)
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
    }

    a.Value += other.Value
    return nil
}

func (Amount) Sum(amounts []Amount) ([]Amount, error) {
    sumsByCurrency := make(map[string]float64)
    for _, amount := range amounts {
        if err := amount.Valid(); err != nil {
            return nil, err
        }
        sumsByCurrency[amount.Currency] += amount.Value
    }

    var resultFormat []Amount
    for currency, sum := range sumsByCurrency {
        resultFormat = append(resultFormat, Amount{sum, currency})
    }
    return resultFormat, nil
}


// ExpenseType
// Methods: String, Valid
type ExpenseType struct {
    ID   int
    Name string
}

func (et ExpenseType) String() string {
    return fmt.Sprintf(et.Name)
}

func (et ExpenseType) Valid() error {
    if et.ID == 0 || et.Name == "" {
        return utils.LogError("expense type ID and name must be non-zero and non-empty")
    }
    return nil
}


// Expense
// Methods: String, Valid, CheckReceipt
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
    return fmt.Sprintf("%v (%d)", e.Amount, e.TypeID)
}

func (e Expense) Valid() error {
    err := e.Amount.Valid()
    switch {
    case e.ID == 0 || e.SessionID == 0 || e.DateTime.IsZero() || e.ReceiptURL == "":
        return utils.LogError("expense ID, session ID, date and time, and receipt URL must be non-zero and non-empty")
    case err != nil:
        return  err
    case e.ID < 0 || e.SessionID < 0:
        return utils.LogError("expense ID and session ID can't be negative")
    default:
        return nil
    }
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
        return errOs
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
        return utils.LogError("error reading headers of receipt image: %v", err)
    }

    buffer = buffer[:n] // Adjust buffer size to the actual number of bytes read
    contentType := http.DetectContentType(buffer)
    switch contentType {
    case "image/jpeg", "image/png", "image/gif", "image/bmp", "image/webp":
        return nil
    default:
        return utils.LogError("invalid receipt image type %s: %s", filePath, contentType)
    }
}
