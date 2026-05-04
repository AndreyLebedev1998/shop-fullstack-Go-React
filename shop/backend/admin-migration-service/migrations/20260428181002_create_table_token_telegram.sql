-- +goose Up
CREATE TABLE telegram_tokens (
    id SERIAL PRIMARY KEY,
    token text
);

-- +goose Down
DROP TABLE telegram_tokens;
