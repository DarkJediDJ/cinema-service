-- +goose Up
CREATE TABLE IF NOT EXISTS public.halls
(
    vip boolean NOT NULL,
    hall_id integer NOT NULL DEFAULT nextval('halls_hall_id_seq'::regclass),
    CONSTRAINT halls_pkey PRIMARY KEY (hall_id)
)
-- +goose Down
DROP TABLE public."hall";
