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

-- +goose Down
DROP TABLE orders;
