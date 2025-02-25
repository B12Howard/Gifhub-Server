-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN postal_code VARCHAR(20);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN IF EXISTS postal_code;

-- +goose StatementEnd
