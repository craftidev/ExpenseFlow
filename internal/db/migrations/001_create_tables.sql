CREATE TABLE IF NOT EXISTS clients (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    name            TEXT     NOT NULL UNIQUE,

    CONSTRAINT ck_non_empty_fields               CHECK (LENGTH(name) > 0),
    CONSTRAINT ck_normal_size_name_100           CHECK (LENGTH(name) <= 100)
);

CREATE TABLE IF NOT EXISTS sessions (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    client_id       INTEGER  NOT NULL,
    name            TEXT     NOT NULL UNIQUE,
    address         TEXT     NOT NULL,
    start_at        TEXT     NOT NULL,
    end_at          TEXT     NOT NULL,

    FOREIGN KEY (client_id)  REFERENCES clients(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT,

    CONSTRAINT ck_non_empty_fields               CHECK (
        LENGTH(name) > 0 AND
        LENGTH(address) > 0 AND
        LENGTH(start_at) > 0 AND
        LENGTH(end_at) > 0
    ),
    CONSTRAINT ck_normal_size_name_address_100   CHECK (
        LENGTH(address) <= 100 AND
        LENGTH(name) <= 100
    ),
    CONSTRAINT ck_normal_size_start_at_end_at_19 CHECK (
        LENGTH(start_at) == 19 AND
        LENGTH(end_at) == 19
    )
);

CREATE TABLE IF NOT EXISTS expense_types (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    name            TEXT     NOT NULL UNIQUE,
    taxe_rate       REAL     NOT NULL,

    CONSTRAINT ck_non_empty_name                 CHECK (LENGTH(name) > 0),
    CONSTRAINT ck_normal_size_name_50            CHECK (LENGTH(name) <= 50),
    CONSTRAINT ck_positive_taxe_rate             CHECK (taxe_rate >= 0),
    CONSTRAINT ck_normal_size_taxe_rate_60       CHECK (taxe_rate <= 60)
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
    notes           TEXT     NOT NULL,

    FOREIGN KEY (session_id) REFERENCES sessions(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT,

    FOREIGN KEY (type_id)    REFERENCES expense_types(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT,

    CONSTRAINT ck_positive_amount                CHECK (amount_value > 0),
    CONSTRAINT ck_non_empty_fields               CHECK (
        LENGTH(amount_currency)    >  0 AND
        LENGTH(location)           >  0 AND
        LENGTH(datetime)           >  0 AND
        LENGTH(receipt_url)        >  0
    ),
    CONSTRAINT ck_normal_size_currency_10        CHECK (LENGTH(amount_currency) <= 10),
    CONSTRAINT ck_normal_size_location_100       CHECK (LENGTH(location) <= 100),
    CONSTRAINT ck_normal_size_datetime_19        CHECK (LENGTH(datetime) == 19),
    CONSTRAINT ck_normal_size_receipt_url_150    CHECK (LENGTH(receipt_url) <= 150), -- TODO reduce
    CONSTRAINT ck_normal_notes_150               CHECK (LENGTH(notes) <= 150)
);
