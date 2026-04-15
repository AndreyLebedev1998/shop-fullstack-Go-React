-- +goose Up
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT,
    product_id INT,
    quantity INT,
    price NUMERIC
);

INSERT INTO order_items (order_id, product_id, quantity, price)
VALUES
    (1, 3, 2, 25.25),
    (1, 1, 1, 50.00),
    (2, 4, 3, 50.00),
    (2, 2, 1, 100.00);

-- +goose Down
DROP TABLE order_items;
