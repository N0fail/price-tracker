-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.price_history (
       code varchar(255) NOT NULL,
       CONSTRAINT fk_product
            FOREIGN KEY(code)
                REFERENCES products(code),
       price float8 NOT NULL,
       date DATE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.price_history;
-- +goose StatementEnd
