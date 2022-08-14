-- +goose Up
-- +goose StatementBegin
INSERT INTO public.products(code, name)
SELECT i::text, i::text
FROM generate_series(1000, 2000) as t(i);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM public.products
USING generate_series(1000, 2000) as t(i)
WHERE code = i::text
-- +goose StatementEnd
