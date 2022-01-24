-- +goose Up
CREATE TABLE IF NOT EXISTS public.movies
(
    name character varying(50) NOT NULL,
    duration time without time zone NOT NULL,
    id integer NOT NULL DEFAULT nextval('movies_movie_id_seq'::regclass),
    CONSTRAINT movies_pkey PRIMARY KEY (id)
)


-- +goose Down
DROP TABLE public.movies;
