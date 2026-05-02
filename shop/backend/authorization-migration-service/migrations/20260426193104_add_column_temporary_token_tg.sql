-- +goose Up
ALTER TABLE users
ADD COLUMN temporary_token_tg TEXT,
ADD COLUMN telegram_token_expires_at TIMESTAMP;

-- +goose Down
ALTER TABLE users
DROP COLUMN temporary_token_tg,
DROP COLUMN telegram_token_expires_at;
