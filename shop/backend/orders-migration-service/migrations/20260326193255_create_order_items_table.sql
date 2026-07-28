-- +goose Up
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT,
    product_id INT,
    quantity INT,
    price NUMERIC
);

-- +goose Down
DROP TABLE order_items;
