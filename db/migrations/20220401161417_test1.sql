-- +goose Up
CREATE TABLE IF NOT EXISTS public.test1
(
    test1 real NOT NULL,
    id1 real NOT NULL
);

-- +goose Down
DROP TABLE public.test1;
