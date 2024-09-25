-- Limit TEXT size on some column

CREATE TABLE IF NOT EXISTS clients (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    name            TEXT     NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS sessions (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    client_id       INTEGER  NOT NULL,
    name            TEXT     NOT NULL UNIQUE,
    address         TEXT     NOT NULL,
    start_at        TEXT     NOT NULL,
    end_at          TEXT     NOT NULL,

    FOREIGN KEY (client_id) REFERENCES clients(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS expense_types (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    name            TEXT     NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS expenses (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id      INTEGER  NOT NULL,
    type_id         INTEGER  NOT NULL,
    amount_value    REAL     NOT NULL, -- mapped in Amount.value
    amount_currency TEXT     NOT NULL, -- mapped in Amount.currency
    location        TEXT     NOT NULL,
    datetime        TEXT     NOT NULL,
    receipt_url     TEXT     NOT NULL,

    FOREIGN KEY (session_id) REFERENCES sessions(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT

    FOREIGN KEY (type_id) REFERENCES expense_types(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
);
