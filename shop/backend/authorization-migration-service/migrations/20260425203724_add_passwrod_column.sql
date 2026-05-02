-- +goose Up
ALTER TABLE users
ADD COLUMN password TEXT;

-- +goose Down
ALTER TABLE users
DROP COLUMN password;
