package db

import (
	"fmt"
	"time"
	"github.com/craftidev/expenseflow/config"
)

// List of models: Client, Session, Amount, ExpenseType, Expense
// (Not its own DB table: Amount)
//
// TODO Validation/Error handling thoroughly

// Client
// Methods: String, Valid
type Client struct {
    ID   int
    Name string
}

func (c Client) String() string {
    return fmt.Sprintf(c.Name)
}

func(c Client) Valid() (Client, error) {
    if c.ID == 0 || c.Name == "" {
        return Client{}, fmt.Errorf("client ID and name must be non-zero and non-empty")
    }
    return c, nil
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

func (s Session) Valid() (Session, error) {
    if s.ID == 0 || s.ClientID == 0 || s.Name == "" || s.Address == "" || s.StartAt.IsZero() || s.EndAt.IsZero() {
        return Session{}, fmt.Errorf("session ID, client ID, name, address, start and end times must be non-zero and non-empty")
    }
    if s.StartAt.After(s.EndAt) {
        return Session{}, fmt.Errorf("start time must be before end time (session ID: %d)", s.ID)
    }
    return s, nil
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

func (a Amount) Valid() (Amount, error) {
    if a.Currency == "" || a.Value == 0 {
        return Amount{}, fmt.Errorf("amount value and currency must be non-empty and non-zero")
    }
    if a.Value < 0 {
        return Amount{}, fmt.Errorf("amount value must be positive")
    }
    return a, nil
}

func (a *Amount) Add(other Amount) error {
    if a.Currency != other.Currency {
        return fmt.Errorf("currencies don't match: %v and %v", a, other)
    }
    a.Value += other.Value
    return nil
}

func (Amount) Sum(amounts []Amount) ([]Amount, error) {
    sumsByCurrency := make(map[string]float64)
    for _, amount := range amounts {
        if _, err := amount.Valid(); err != nil {
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

func (et ExpenseType) Valid() (ExpenseType, error) {
    if et.ID == 0 || et.Name == "" {
        return ExpenseType{}, fmt.Errorf("expense type ID and name must be non-zero and non-empty")
    }
    return et, nil
}


// Expense
// Methods: String, Valid, HasReceipt
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

// TODO in services, even if we have a default URL (pointing at a default img representing empty receipt),
// we need to check before creating a report that ALL Receipts are real ones not default
func (e Expense) Valid() (Expense, error) {
    var amountZeroValue Amount
    if e.ID == 0 || e.SessionID == 0 || e.DateTime.IsZero() || e.Amount == amountZeroValue || e.ReceiptURL == "" {
        return Expense{}, fmt.Errorf("expense ID, session ID, amount, date and time, and receipt URL must be non-zero and non-empty")
    }
    if _, err := e.Amount.Valid(); err != nil {
        return Expense{}, fmt.Errorf("invalid amount: %v", err)
    }
    return e, nil
}

func (e Expense) HasReceipt() bool {
    return e.ReceiptURL!= config.DefaultReceiptURL
}
