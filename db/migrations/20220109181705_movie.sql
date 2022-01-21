-- +goose Up
CREATE TABLE IF NOT EXISTS public.movie
(
    "movie_ID" integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    duration time without time zone NOT NULL,
    CONSTRAINT "Movie_pkey" PRIMARY KEY ("movie_ID")
)


-- +goose Down
DROP TABLE public.movie;