-- +goose Up
CREATE TABLE IF NOT EXISTS public.halls
(
    vip boolean NOT NULL,
    id SERIAL,
    seats integer,
    CONSTRAINT halls_pkey PRIMARY KEY (id)
);
-- +goose Down
DROP TABLE public.halls;
