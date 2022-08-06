-- +goose Up
-- +goose StatementBegin
INSERT INTO public.price_history(code, price, date)
SELECT i::text, random()*1000, timestamp '2014-01-10 20:00:00' +
                               random() * (timestamp '2014-01-20 20:00:00' -
                                           timestamp '2014-01-10 10:00:00')
FROM generate_series(1000, 2000) as t(i), generate_series(1, 10);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM public.price_history
    USING generate_series(1000, 2000) as t(i)
WHERE code = i::text
-- +goose StatementEnd
