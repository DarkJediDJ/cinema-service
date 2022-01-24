-- +goose Up
CREATE TABLE IF NOT EXISTS public.users
(
    login text COLLATE NOT NULL,
    password text COLLATE NOT NULL,
    user_id integer NOT NULL DEFAULT nextval('users_user_id_seq'::regclass),
    CONSTRAINT users_pkey PRIMARY KEY (user_id)
)

-- +goose Down
DROP TABLE public.users;
