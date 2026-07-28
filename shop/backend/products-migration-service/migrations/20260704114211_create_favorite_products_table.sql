-- +goose Up
CREATE TABLE favorite_products (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    product_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),

    UNIQUE (user_id, product_id)
);

-- +goose Down
DROP TABLE favorite_products;
