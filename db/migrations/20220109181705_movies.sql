-- +goose Up
CREATE TABLE IF NOT EXISTS public.movies
(
    name character varying(50) NOT NULL,
    duration time without time zone NOT NULL,
    movie_id integer NOT NULL DEFAULT nextval('movies_movie_id_seq'::regclass),
    CONSTRAINT movies_pkey PRIMARY KEY (movie_id)
)


-- +goose Down
DROP TABLE public.movie;
