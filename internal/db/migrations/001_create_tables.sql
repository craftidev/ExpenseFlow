CREATE TABLE IF NOT EXISTS clients (
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT    NOT NULL UNIQUE,

    CONSTRAINT ck_non_empty_fields     CHECK (LENGTH(name) > 0),
    CONSTRAINT ck_normal_size_name_100 CHECK (LENGTH(name) <= 100)
);

CREATE TABLE IF NOT EXISTS sessions (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    client_id           INTEGER  NOT NULL,
    location            TEXT     NOT NULL,
    trip_start_location TEXT         NULL,
    trip_end_location   TEXT         NULL,
    start_at_date       TEXT         NULL,
    end_at_date         TEXT         NULL,

    FOREIGN KEY (client_id)  REFERENCES clients(id),

    CONSTRAINT ck_non_empty_fields                 CHECK (
        LENGTH(location)    > 0 AND
        trip_start_location == NULL OR LENGTH(trip_start_location) > 0 AND
        trip_end_location   == NULL OR LENGTH(trip_end_location)   > 0
    ),
    CONSTRAINT ck_normal_size_name_locations_100   CHECK (
        LENGTH(location)            <= 100 AND
        LENGTH(trip_start_location) <= 100 AND
        LENGTH(trip_end_location)   <= 100
    ),
    CONSTRAINT ck_normal_size_start_end_at_date_19 CHECK (
        start_at_date == NULL OR LENGTH(start_at_date) == 19 AND
        end_at_date   == NULL OR LENGTH(end_at_date)   == 19
    )
);

CREATE TABLE IF NOT EXISTS car_trips(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id  INTEGER     NULL,
    distance_km REAL    NOT NULL,
    date_time   TEXT    NOT NULL UNIQUE,

    FOREIGN KEY (session_id) REFERENCES sessions(id),

    CONSTRAINT ck_normal_size_date_19  CHECK (LENGTH(date_time) == 19),
    CONSTRAINT ck_positive_distance_km CHECK (distance_km       > 0)
);

CREATE TABLE IF NOT EXISTS expense_types (
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT    NOT NULL UNIQUE,

    CONSTRAINT ck_non_empty_name      CHECK (LENGTH(name) > 0),
    CONSTRAINT ck_normal_size_name_50 CHECK (LENGTH(name) <= 50)
);

CREATE TABLE IF NOT EXISTS expenses (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id  INTEGER     NULL,
    type_id     INTEGER NOT NULL,
    currency    TEXT    NOT NULL,
    receipt_url TEXT        NULL,
    notes       TEXT        NULL,
    date_time   TEXT    NOT NULL,

    FOREIGN KEY (session_id) REFERENCES sessions(id),
    FOREIGN KEY (type_id)    REFERENCES expense_types(id),

    CONSTRAINT ck_non_empty_currency         CHECK (LENGTH(currency)    >  0),
    CONSTRAINT ck_normal_size_currency_10    CHECK (LENGTH(currency)    <= 10),
    CONSTRAINT ck_normal_size_receipt_url_50 CHECK (LENGTH(receipt_url) <= 50),
    CONSTRAINT ck_normal_notes_150           CHECK (LENGTH(notes)       <= 150),
    CONSTRAINT ck_normal_date_time           CHECK (LENGTH(date_time)   == 19)
);

CREATE TABLE IF NOT EXISTS line_items (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    expense_id INTEGER NOT NULL,
    taxe_rate  REAL    NOT NULL,
    total      REAL    NOT NULL,

    FOREIGN KEY (expense_id) REFERENCES expenses(id),

    CONSTRAINT ck_positive_total           CHECK (taxe_rate > 0),
    CONSTRAINT ck_positive_taxe_rate       CHECK (taxe_rate >= 0),
    CONSTRAINT ck_normal_size_taxe_rate_60 CHECK (taxe_rate <= 60)
);
