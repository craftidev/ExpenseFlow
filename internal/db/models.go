package db

import (
	"fmt"
	"time"
    "os"
    "errors"
	"github.com/craftidev/expenseflow/config"
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
        return fmt.Errorf("client ID and name must be non-zero and non-empty")
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
        return fmt.Errorf("session ID, client ID, name, address, start and end times must be non-zero and non-empty")
    }
    if s.StartAt.After(s.EndAt) {
        return fmt.Errorf("start time must be before end time (session ID: %d)", s.ID)
    }
    return nil
}


// Amount
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
        return fmt.Errorf("amount value and currency must be non-empty and non-zero")
    }
    if a.Value < 0 {
        return fmt.Errorf("amount value must be positive")
    }
    return nil
}

func (a *Amount) Add(other Amount) error {
    switch {
    case a.Valid() != nil:
        return fmt.Errorf("invalid first amount: %v", a.Valid())
    case other.Valid() != nil:
        return fmt.Errorf("invalid second amount: %v", other.Valid())
    case a.Currency != other.Currency:
        return fmt.Errorf("currencies don't match: %v and %v", a.Currency, other.Currency)
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
        return fmt.Errorf("expense type ID and name must be non-zero and non-empty")
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
    switch {
    case e.ID == 0 || e.SessionID == 0 || e.DateTime.IsZero() || e.ReceiptURL == "":
        return fmt.Errorf("expense ID, session ID, date and time, and receipt URL must be non-zero and non-empty")
    case e.Amount.Valid() != nil:
        return fmt.Errorf("invalid amount: %v", e.Amount.Valid())
    case e.ID < 0 || e.SessionID < 0:
        return fmt.Errorf("expense ID and session ID can't be negative")
    default:
        return nil
    }
}

// TODO probably will have to test with Flutter if the img is corrupted and can't show
func (e Expense) CheckReceipt() error {
    _, err := os.Stat(config.Path + e.ReceiptURL)
    switch {
    case e.ReceiptURL == config.DefaultReceiptURL:
        return fmt.Errorf("receipt is the default placeholder image")
    case e.ReceiptURL == "" :
        return fmt.Errorf("receipt URL is empty")
    case errors.Is(err, os.ErrNotExist):
        return fmt.Errorf("invalid receipt URL (with path config): %s%s", config.Path, e.ReceiptURL)
    case err != nil:
        return err
    default:
        return nil
    }
}
