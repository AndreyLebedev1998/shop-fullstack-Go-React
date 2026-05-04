-- +goose Up
ALTER TABLE users
ADD COLUMN code_recovery TEXT,
ADD COLUMN code_recovery_expires_at TIMESTAMP;

-- +goose Down
ALTER TABLE users
DROP COLUMN code_recovery,
DROP COLUMN code_recovery_expires_at;
