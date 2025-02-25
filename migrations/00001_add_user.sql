-- +goose Up
-- +goose StatementBegin
CREATE TABLE times_available (
    id SERIAL PRIMARY KEY,
    m BOOLEAN DEFAULT FALSE,
    t BOOLEAN DEFAULT FALSE,
    w BOOLEAN DEFAULT FALSE,
    th BOOLEAN DEFAULT FALSE,
    f BOOLEAN DEFAULT FALSE,
    sa BOOLEAN DEFAULT FALSE,
    su BOOLEAN DEFAULT FALSE,
    time1 BOOLEAN DEFAULT FALSE,
    time2 BOOLEAN DEFAULT FALSE,
    time3 BOOLEAN DEFAULT FALSE,
    time4 BOOLEAN DEFAULT FALSE
);

CREATE TABLE kith_available (
    id SERIAL PRIMARY KEY,
    m BOOLEAN DEFAULT FALSE,
    t BOOLEAN DEFAULT FALSE,
    w BOOLEAN DEFAULT FALSE,
    th BOOLEAN DEFAULT FALSE,
    f BOOLEAN DEFAULT FALSE,
    sa BOOLEAN DEFAULT FALSE,
    su BOOLEAN DEFAULT FALSE
);

CREATE TABLE members (
    id SERIAL PRIMARY KEY,
    my_headline VARCHAR NOT NULL,
    about_me VARCHAR,
    postal_code VARCHAR,
    street_address VARCHAR,
    apt_number VARCHAR,
    city VARCHAR,
    state VARCHAR(15),
    additional_information VARCHAR
);

CREATE TABLE kiths (
    id SERIAL PRIMARY KEY,
    public_profile VARCHAR NOT NULL,
    travel_distance INT NOT NULL,
    start_work_date DATE NOT NULL,
    additional_information VARCHAR NULL,
    my_headline VARCHAR NOT NULL,
    about_me VARCHAR NULL,
    postal_code VARCHAR(10) NULL, -- Assuming US ZIP codes
    street_address VARCHAR NULL,
    apt_number VARCHAR NULL,
    city VARCHAR NULL,
    state VARCHAR(15) NULL,
    days_available INT REFERENCES kith_available(id),
    rate NUMERIC(9,2),
    times_available INT REFERENCES times_available(id)
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(15) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    kith INT UNIQUE REFERENCES kiths(id) ON DELETE SET NULL,
    member INT UNIQUE REFERENCES members(id) ON DELETE SET NULL,
    first_name VARCHAR(80),
    last_name VARCHAR(80),
    email_confirmed BOOLEAN,
    active_subscription BOOLEAN,
    external_id VARCHAR(128)
);

CREATE TABLE unavailable_days (
    id SERIAL PRIMARY KEY,
    kith INT REFERENCES kiths(id) ON DELETE CASCADE,
    unavailable_date DATE NOT NULL
);

CREATE TABLE experiences (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE
);

CREATE TABLE certifications (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE
);

CREATE TABLE kith_experiences (
    id SERIAL PRIMARY KEY,
    kith INT REFERENCES kiths(id) ON DELETE CASCADE,
    experience INT REFERENCES experiences(id) ON DELETE CASCADE
);

CREATE TABLE kith_certifications (
    id SERIAL PRIMARY KEY,
    kith INT REFERENCES kiths(id) ON DELETE CASCADE,
    certification INT REFERENCES certifications(id) ON DELETE CASCADE
);

-- DTO and demo storage only, not stored in real DB because of compliance concerns
CREATE TABLE bank_accounts (
    id SERIAL PRIMARY KEY,
    kith INT UNIQUE REFERENCES kiths(id) ON DELETE CASCADE,
    name VARCHAR NOT NULL,
    routing_number VARCHAR(9) NOT NULL CHECK (routing_number ~ '^[0-9]{9}$'), -- US standard
    account_number VARCHAR NOT NULL -- Length varies by bank
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bank_accounts;
DROP TABLE kith_certifications;
DROP TABLE kith_experiences;
DROP TABLE certifications;
DROP TABLE experiences;
DROP TABLE kith_available;
DROP TABLE unavailable_days;
DROP TABLE times_available;
DROP TABLE users;
DROP TABLE kiths;
DROP TABLE members;
-- +goose StatementEnd
