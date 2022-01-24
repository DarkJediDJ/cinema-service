-- +goose Up
CREATE TABLE IF NOT EXISTS public.sessions
(
    schedule time with time zone NOT NULL,
    hall_id integer NOT NULL,
    movie_id integer NOT NULL,
    id integer SERIAL,
    CONSTRAINT sessions_pkey PRIMARY KEY (id),
    CONSTRAINT "FK_sessions_to_halls" FOREIGN KEY (hall_id)
        REFERENCES public.halls (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT "FK_sessions_to_movie" FOREIGN KEY (movie_id)
        REFERENCES public.movies (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
)

-- +goose Down
DROP TABLE public.sessions;
