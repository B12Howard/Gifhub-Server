-- +goose Up
-- +goose StatementBegin
CREATE TABLE credit_cards (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name_on_card VARCHAR NOT NULL,
    card_number VARCHAR NOT NULL,
    expiration_date VARCHAR(5) NOT NULL,
    ccv VARCHAR(4) CHECK (LENGTH(ccv) IN (3, 4)) NOT NULL,
    card_postal_code VARCHAR,
    user_type VARCHAR NOT NULL CHECK (user_type IN ('member', 'kith'))
);

CREATE TABLE vehicle (
    id SERIAL PRIMARY KEY,
    make VARCHAR NOT NULL,
    model VARCHAR NOT NULL,
    year VARCHAR NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS vehicle;
DROP TABLE IF EXISTS credit_cards;
-- +goose StatementEnd
