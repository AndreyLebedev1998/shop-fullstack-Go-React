-- +goose Up
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    product_name text,
    category_id INT NOT NULL,
    price NUMERIC,
    image_url TEXT
);

-- +goose Down
DROP TABLE products;