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
    ID      int
    Client  Client
    StartAt time.Time
    EndAt   time.Time
}

type Category struct {
    ID          int
    Name        string
    Description string
}

type Reason struct {
    ID          int
    Name        string
    Description string
    Category    Category
}

type Expense struct {
    ID          int
    Session     Session
    Description string
    Amount      float64
    Reason      Reason
    DateTime    time.Time
}
