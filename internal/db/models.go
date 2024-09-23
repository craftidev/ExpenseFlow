package db

import (
    "time"
)

type Client struct {
    ID      int
    Name    string
    Address string
    Phone   string
    Email   string
}

type Session struct {
    ID       int
    ClientID int
    StartAt  time.Time
    EndAt    time.Time
}

type Category struct {
    ID          int
    Name        string
    Description string
}

type Reason struct {
    ID          int
    CategoryID  int
    Name        string
    Description string
}

type Expense struct {
    ID          int
    SessionID   int
    ReasonID    int
    Description string
    Amount      float64
    DateTime    time.Time
}
