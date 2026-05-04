-- +goose Up
ALTER TABLE users
ADD COLUMN token_recovery_password TEXT,
ADD COLUMN token_recovery_password_expires_at TIMESTAMP;

-- +goose Down
ALTER TABLE users
DROP COLUMN token_recovery_password,
DROP COLUMN token_recovery_password_expires_at;
