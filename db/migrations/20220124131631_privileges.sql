-- +goose Up
CREATE TABLE IF NOT EXISTS public.privileges
(
    privilege_id integer NOT NULL DEFAULT nextval('privileges_privilege_id_seq'::regclass),
    name text NOT NULL,
    CONSTRAINT privileges_pkey PRIMARY KEY (privilege_id)
)

-- +goose Down
DROP TABLE public.privileges;
