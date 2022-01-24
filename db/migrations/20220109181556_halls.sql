-- +goose Up
CREATE TABLE IF NOT EXISTS public.halls
(
    vip boolean NOT NULL,
    id integer NOT NULL DEFAULT nextval('halls_hall_id_seq'::regclass),
    CONSTRAINT halls_pkey PRIMARY KEY (id)
)
-- +goose Down
DROP TABLE public.halls;
