-- +goose Up
ALTER TABLE products
ADD COLUMN description TEXT DEFAULT '',
ADD COLUMN tags TEXT[] NOT NULL DEFAULT '{}';

-- +goose Down
ALTER TABLE products
DROP COLUMN description,
DROP COLUMN tags;
