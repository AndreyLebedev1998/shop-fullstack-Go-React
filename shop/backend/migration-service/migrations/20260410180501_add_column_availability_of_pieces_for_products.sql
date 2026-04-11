-- +goose Up
ALTER TABLE products
ADD COLUMN availability_of_pieces INT DEFAULT 0;

-- +goose Down
ALTER TABLE products
DROP COLUMN availability_of_pieces;
