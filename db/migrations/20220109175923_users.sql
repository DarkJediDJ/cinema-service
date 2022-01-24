-- +goose Up
CREATE TABLE IF NOT EXISTS public.users
(
    login text NOT NULL,
    password text NOT NULL,
    id integer NOT NULL DEFAULT nextval('users_user_id_seq'::regclass),
    CONSTRAINT users_pkey PRIMARY KEY (id)
)


-- +goose Down
DROP TABLE public.users;
