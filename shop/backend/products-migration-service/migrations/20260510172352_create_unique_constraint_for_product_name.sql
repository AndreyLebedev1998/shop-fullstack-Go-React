-- +goose Up
ALTER TABLE products ADD CONSTRAINT products_product_name_unique UNIQUE (product_name);

-- +goose Down
ALTER TABLE products DROP CONSTRAINT products_product_name_unique;
