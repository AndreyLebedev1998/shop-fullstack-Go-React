-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT,
    lastname TEXT,
    email TEXT UNIQUE,
    phone TEXT UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    CHECK (email IS NOT NULL OR phone IS NOT NULL)
);

-- +goose Down
DROP TABLE users;

