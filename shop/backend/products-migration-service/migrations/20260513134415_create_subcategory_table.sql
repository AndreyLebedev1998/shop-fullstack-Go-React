-- +goose Up
CREATE TABLE subcategories (
    id SERIAL PRIMARY KEY,
    category_id int,
    category_name text,

    CONSTRAINT FK_subcategories_categories_id FOREIGN KEY(category_id) REFERENCES categories(id)
);

-- +goose Down
DROP TABLE subcategories;
