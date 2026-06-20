-- +goose Up
CREATE TABLE admins (
    id SERIAL PRIMARY KEY,
    nick TEXT UNIQUE NOT NULL,
    password TEXT
);

-- +goose Down
DROP TABLE admins;