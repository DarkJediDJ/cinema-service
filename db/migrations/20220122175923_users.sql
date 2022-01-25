-- +goose Up
CREATE TABLE IF NOT EXISTS public.users
(
    login text NOT NULL,
    password text NOT NULL,
    id integer SERIAL,
    CONSTRAINT users_pkey PRIMARY KEY (id)
)


-- +goose Down
DROP TABLE public.users;
