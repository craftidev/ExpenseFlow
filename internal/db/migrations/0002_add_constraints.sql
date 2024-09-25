-- No zero or negative
ALTER TABLE expenses
ADD CONSTRAINT ck_positive_amount
CHECK (amount_value            >  0);

-- No empty string
ALTER TABLE clients
ADD CONSTRAINT ck_non_empty_fields
CHECK (LENGTH(name)            >  0);

ALTER TABLE sessions
ADD CONSTRAINT ck_non_empty_fields
CHECK (
    LENGTH(name)               >  0 AND
    LENGTH(address)            >  0 AND
    LENGTH(start_at)           >  0 AND
    LENGTH(end_at)             >  0
);

ALTER TABLE expense_types
ADD CONSTRAINT ck_non_empty_name
CHECK (LENGTH(name)            >  0);

ALTER TABLE expenses
ADD CONSTRAINT ck_non_empty_fields
CHECK (
    LENGTH(amount_currency)    >  0 AND
    LENGTH(location)           >  0 AND
    LENGTH(datetime)           >  0 AND
    LENGTH(receipt_url)        >  0
);

-- Size limitation
ALTER TABLE clients
ADD CONSTRAINT ck_normal_size_name_100
CHECK (LENGTH(name)            <= 100);

ALTER TABLE sessions
ADD CONSTRAINT ck_normal_size_name_address_100
CHECK (
    LENGTH(address) <          =  100 AND
    LENGTH(name)               <= 100
);

ALTER TABLE sessions
ADD CONSTRAINT ck_normal_size_start_at_end_at_19
CHECK (
    LENGTH(start_at)           == 19 AND
    LENGTH(end_at)             == 19
);

ALTER TABLE expense_types
ADD CONSTRAINT ck_normal_size_name_50
CHECK (LENGTH(name)            <= 50);

ALTER TABLE expenses
ADD CONSTRAINT ck_normal_size_currency_10
CHECK (LENGTH(amount_currency) <= 10);

ALTER TABLE expenses
ADD CONSTRAINT ck_normal_size_location_100
CHECK (LENGTH(location)        <= 100);

ALTER TABLE expenses
ADD CONSTRAINT ck_normal_size_datetime_19
CHECK (LENGTH(datetime)        == 19);

ALTER TABLE expenses
ADD CONSTRAINT ck_normal_size_receipt_url_150
CHECK (LENGTH(receipt_url)     <= 150);
