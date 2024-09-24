# ExpenseFlow
**Track and report professional expenses**
*ExpenseFlow is an app designed to help professionals track and report their expenses. It simplifies the process of managing receipts, categorizing expenses, and generating reports, all from a single platform.*

### Expected Features
- 🚧 **Session-based expense tracking:** Log expenses against a session (client or mission).
- 🚧 **Expense details:** Track the reason, value, date, and time of each expense.
- 🚧 **Receipt capture:** Upload photos of receipts.
- 🚧 **Reports:** Generate reports in a specific format.
- 🚧 **Cross-platform support:** Available on Android, iOS, Windows, Mac, Linux and Online.

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
### Install
🚧
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

## Problems and solutions
🚧
## Design choices
🚧
## Pending question
- Do I need a currency converter?
- Does the report show multiple totals depending on currency?
- What are the models/properties I'm missing?
- Is there a use to know outgoing or incoming travels?
- Use for town goal?
- Expense.Location just town names?
- Does Client need contact info?
- Does Expense need Description or OtherInfo?
