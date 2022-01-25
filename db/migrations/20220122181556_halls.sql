-- +goose Up
CREATE TABLE IF NOT EXISTS public.halls
(
    vip boolean NOT NULL,
    id integer SERIAL,
    seats integer,
    CONSTRAINT halls_pkey PRIMARY KEY (id)
)
-- +goose Down
DROP TABLE public.halls;
