-- +goose Up
CREATE TABLE IF NOT EXISTS public."User"
(
    "user_ID" integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    "login" text COLLATE pg_catalog."default" NOT NULL,
    "password" text COLLATE pg_catalog."default" NOT NULL,
    "add_halls" boolean NOT NULL,
    "add_movies" boolean NOT NULL,
    "add_sessions" boolean NOT NULL,
    CONSTRAINT "user_pkey" PRIMARY KEY ("user_ID")
)

-- +goose Down
DROP TABLE public."user";
