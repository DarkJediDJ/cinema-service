-- +goose Up
CREATE TABLE IF NOT EXISTS public.movies
(
    name character varying(50) NOT NULL,
    duration INTERVAL,
    id SERIAL,
    CONSTRAINT movies_pkey PRIMARY KEY (id)
);


-- +goose Down
DROP TABLE public.movies;
