-- +goose Up
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    category_name text
);

-- +goose Down
DROP TABLE categories;
