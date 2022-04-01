-- +goose Up
CREATE TABLE IF NOT EXISTS public.test
(
    test real NOT NULL,
    id real NOT NULL
);

-- +goose Down
DROP TABLE public.test;
