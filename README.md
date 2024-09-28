# ExpenseFlow
**Track and report professional expenses**
*ExpenseFlow is an app designed to help professionals track and report their expenses. It simplifies the process of managing receipts, categorizing expenses, and generating reports, all from a single platform.*

### 🚧 Features 🚧
- **Session-based expense tracking:** Log expenses against a session (client or mission).
- **Expense details:** Track the reason, value, date, and time of each expense.
- **Receipt capture:** Upload photos of receipts.
- **Reports:** Generate reports in a specific format.
- **Cross-platform support:** Available on Android, iOS, Windows, Mac, Linux and Online.

### Tech Stack
- **Backend:** Go (Golang)
- **Frontend:** Flutter
- **Database:** SQLite
- **License:** Apache 2.0

### Roadmap
1. **Backend Development (Go)**
   - [ ] Define data models (expenses, sessions, reasons, etc.)
   - [ ] Set up API routes for expense tracking:
     - [ ] POST `/expenses`: Add a new expense
     - [ ] GET `/expenses`: Get a list of all expenses
     - [ ] POST `/receipts`: Upload a receipt
     - [ ] GET `/reports`: Generate a report
   - [ ] Set up error handling, logging, and unit testing
   - [ ] Write documentation for the backend API (OpenAPI/Swagger)

2. **Frontend Development (Flutter)**
   - [ ] Create basic UI for inputting expenses and sessions
   - [ ] Add functionality to upload receipt photos
   - [ ] Implement a report viewer
   - [ ] Handle authentication (connect to backend)
   - [ ] Perform user testing

3. **Deployment**
   - [ ] Set up web hosting (use platforms like Firebase, DigitalOcean, etc.)
   - [ ] Prepare Android/iOS builds for Google Play and the Apple App Store
   - [ ] Ensure cross-platform compatibility for desktop versions

4. **Future Features**
   - [ ] Push notifications for expense reminders
   - [ ] Multi-language support
   - [ ] Integration with other tools (e.g., Google Drive for backup)
   - [ ] Customizable filters/presentation for the report
   - [ ] Export in CSV/PDF/...

---

# Development Journal
## First impressions
### 1. **Roadmap Details**

- **Data Models (Go)**:
  ```go
  type Expense struct {
      ID         int
      SessionID  int
      Reason     string
      Value      float64
      DateTime   time.Time
      ReceiptURL string
  }

  type Session struct {
      ID      int
      Client  string
      Mission string
  }
  ```

- **API Design**: Plan API endpoints. Need routes for adding expenses, uploading receipts, and generating reports.

- **Authentication**: Can skip authentication early on, but eventually, want to handle user accounts. JWT (JSON Web Tokens) could be a good fit.

- **Testing**: Write tests for the API routes as I go. Go has built-in testing functionality (`go test`).

### 2. **Frontend Focus (Flutter)**
Once the backend is functional, shift focus to Flutter. Keep it simple at first:
- A form to add expenses
- A file picker for receipt images
- A report generator that formats data nicely

## Reminders and hints
### Testing:
```bash
go test ./tests
```

Change go test timeout in VS Code:
```json
{
    "go.testTimeout": "5s"
}

```

### Running:
```bash
go run ./cmd/expenseflow
```

### Install
Modify path in `/config/config.go`

## Important decisions
### Choice of stack: Go / Flutter / SQLite:
**Go** is fast, simple and clear. After some courses on [boot.dev](https://boot.dev) I got attracted to it and wanted to make a full project to really compare to my other experiences in Python, my language of choice for many years, but also PHP, Java, TypeScript...

**Flutter** is my default choice in frontend. It's cross-platform most of all. I need more experience in it.

**SQLite** is also my default choice for database. Lightweight, it integrates well in mobile. ExpenseFlow doesn't have a huge amount of complex data to handle.

### Codebase architecture
```tree
ExpenseFlow/
│
├── cmd/
│   └── expenseflow/
│       └── main.go         # Entry point of the application
│
├── internal/
│   ├── db/
│   │   ├── models.go       # Database models (e.g., Expense, Session)
│   │   ├── queries.go      # Database query functions
│   │   └── init.go         # Database connection setup
│   ├── handlers/
│   │   ├── expenses.go     # API route handlers (e.g., for adding expenses)
│   │   └── sessions.go     # API route handlers (e.g., for sessions)
│   └── services/
│       └── report.go       # Business logic (e.g., generating reports)
│
├── pkg/
│   └── middleware/         # Any reusable middleware (e.g., authentication, logging)
│
├── api/
│   └── routes.go           # Routes definition (e.g., registering API endpoints)
│
├── config/
│   └── config.go           # Configuration (e.g., environment variables, app settings)
│
├── .gitignore
├── go.mod
└── README.md
```

### Date/Time
Everything will use the standard library "time" and will record UTC timestamps. Making the user able travel in different timezone and not confuse any logic. Flutter will be the one taking care of the preferred display for the user.

### Receipt handling
I chose to have non-nullable in the DB. A placeholder IMG. So when I test Expense for validation it's not the same as when I `CheckReceipt()`.
The latter will make sure the file exist and is not the placeholder.
I think the user could add the picture later in their workflow when adding expense.

### Hard limiting Float (for `Amount.Value` operations)
After creating some test to identify when Add or Sum were creating a float `> math.MaxFloat64`. I realized they were weird behaviors. You can't subtract a small float from a giant one, the result is unchanged. So I decided to hard code an unrealistic max in `/config/config.go` at `1_000_000_000.0`

### Datatype for IDs
sqlite3 drivers are returning `int64` for ID columns. I decided to stick with `int` datatype in go (32 or 64 depending on the machine running.) It's very unlikely that I'll ever need `int64`, but I added a validation with `Fatal` if it ever occurs.

### Logging
I will use the tag `[info]` for what's not fatal/error. I'll try to not crowd the logs with useless logic event. But for now any change to the DB is logged. And to avoid any security/privacy breach I'll only log IDs for information.

## Problems and solutions
🚧
## Design choices
🚧
## Questions to users
- Do I need a currency converter? Does the report show multiple totals depending on currency?
-> `Multiple reports in case of multiple currencies.`
- What are the models/properties I'm missing?
-> `Standard categories`, `Additional expense comment: Observation`, `Taxes by categories`, `Expenses can have Session = NULL`, `KM by session`, `opt VILLE > MISSION > opt VILLE`, `Session.Town, Session.ZipCode`
