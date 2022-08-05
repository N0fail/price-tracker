-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.products (
    code varchar(255) NOT NULL PRIMARY KEY,
    name varchar(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.products;
-- +goose StatementEnd
