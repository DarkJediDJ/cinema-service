-- +goose Up
CREATE TABLE IF NOT EXISTS public.privileges
(
    id integer NOT NULL DEFAULT nextval('privileges_privilege_id_seq'::regclass),
    name text NOT NULL,
    CONSTRAINT privileges_pkey PRIMARY KEY (id)
)

-- +goose Down
DROP TABLE public.privileges;
