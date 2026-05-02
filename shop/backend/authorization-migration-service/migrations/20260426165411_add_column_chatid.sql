-- +goose Up
ALTER TABLE users
ADD COLUMN chat_id_telegram BIGINT;

-- +goose Down
ALTER TABLE users
DROP COLUMN chat_id_telegram;
