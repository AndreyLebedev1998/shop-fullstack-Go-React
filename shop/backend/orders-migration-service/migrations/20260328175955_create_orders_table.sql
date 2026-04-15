-- +goose Up
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT,
    email TEXT,
    phone TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    total_price NUMERIC(10, 2) DEFAULT 0,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    CHECK (user_id IS NOT NULL OR email IS NOT NULL OR phone IS NOT NULL)
);

INSERT INTO orders (user_id, email, phone, status, total_price)
VALUES 
    (1, NULL, NULL, 'pending', 100.50),
    (NULL, 'test@example.com', NULL, 'pending', 250.00);

-- +goose Down
DROP TABLE orders;
