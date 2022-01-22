-- +goose Up
CREATE TABLE IF NOT EXISTS public.users
(
    login text NOT NULL,
    password text NOT NULL,
    add_halls boolean NOT NULL,
    add_movies boolean NOT NULL,
    add_sessions boolean NOT NULL,
    user_id integer NOT NULL DEFAULT nextval('users_user_id_seq'::regclass),
    CONSTRAINT users_pkey PRIMARY KEY (user_id)
)

-- +goose Down
DROP TABLE public."user";
