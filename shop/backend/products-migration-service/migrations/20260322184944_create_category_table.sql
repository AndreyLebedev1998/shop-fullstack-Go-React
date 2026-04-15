-- +goose Up
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    category_name text
);

INSERT INTO categories (category_name)
VALUES ('energy drinks');

-- +goose Down
DROP TABLE categories;
