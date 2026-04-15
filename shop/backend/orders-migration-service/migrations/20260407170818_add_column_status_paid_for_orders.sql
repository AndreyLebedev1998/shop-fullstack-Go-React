-- +goose Up
ALTER TABLE orders
ADD COLUMN status_paid text DEFAULT 'not_paid';

-- +goose Down
ALTER TABLE orders
DROP COLUMN status_paid;
