-- +goose Up
CREATE TABLE IF NOT EXISTS public.session
(
    "session_ID" integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    "hall_ID" integer NOT NULL,
    "movie_ID" integer NOT NULL,
    "time" time with time zone NOT NULL,
    CONSTRAINT "Session_pkey" PRIMARY KEY ("session_ID"),
    CONSTRAINT "FK_<Session>_<Hall>" FOREIGN KEY ("hall_ID")
        REFERENCES public.hall ("hall_ID") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT "FK_<Session>_<Movie>" FOREIGN KEY ("movie_ID")
        REFERENCES public.movie ("movie_ID") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

-- +goose Down
DROP TABLE public."session";