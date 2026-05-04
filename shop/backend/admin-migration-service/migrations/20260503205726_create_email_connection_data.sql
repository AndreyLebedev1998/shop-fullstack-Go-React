-- +goose Up
CREATE TABLE email_conn (
    id SERIAL PRIMARY KEY,
    sender_email text,
    password text
);

-- +goose Down
DROP TABLE email_conn;

