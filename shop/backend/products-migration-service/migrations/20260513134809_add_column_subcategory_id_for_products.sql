-- +goose Up
ALTER TABLE products
ADD COLUMN subcategory_id INT;

-- +goose Down
ALTER TABLE products
DROP COLUMN subcategory_id;
