# ExpenseFlow
**Track and report professional expenses**
*ExpenseFlow is an app designed to help professionals track and report their expenses. It simplifies the process of managing receipts, categorizing expenses, and generating reports, all from a single platform.*

### Features
- **Session-based expense tracking:** Log expenses against a session (client or mission).
- **Expense details:** Track the reason, value, date, and time of each expense.
- **Receipt capture:** Upload photos of receipts.
- **Custom reports:** Generate reports in a specific format.
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
   - [ ] Integrate a database (SQLite)
   - [ ] Handle authentication (OAuth2 or JWT)
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
## Problems and solutions
🚧
## Design choices
🚧
