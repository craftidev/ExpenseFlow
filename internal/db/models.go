package db

import (
	"fmt"
	"time"
)


// List of models: Client, Session, Amount, Expense
// (Not its own DB table: Amount)

// Client
type Client struct {
    ID   int
    Name string
}

func (c Client) String() string {
    return fmt.Sprintf(c.Name)
}


// Session
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


// Amount
type Amount struct {
    Value    float64
    Currency string
}

func (a Amount) String() string {
    return fmt.Sprintf("%.2f %s", a.Value, a.Currency)
}


// Expense
type Expense struct {
    ID         int
    SessionID  int
    Type       string
    Amount     Amount
    Location   string
    DateTime   time.Time
    ReceiptURL string
}

func (e Expense) String() string {
    return fmt.Sprintf("%s (%v)", e.Amount, e.Type)
}
