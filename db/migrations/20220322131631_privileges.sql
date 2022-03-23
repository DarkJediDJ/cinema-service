-- +goose Up
CREATE TABLE IF NOT EXISTS public.privileges
(
    id SERIAL,
    name text NOT NULL,
    CONSTRAINT privileges_pkey PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE public.privileges;
