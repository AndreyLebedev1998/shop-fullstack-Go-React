-- +goose Up
CREATE TABLE product_stats (
    product_id INT PRIMARY KEY REFERENCES products(id) ON DELETE CASCADE,
    purchase_count INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE product_stats;
