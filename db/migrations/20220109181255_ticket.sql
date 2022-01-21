-- +goose Up
CREATE TABLE IF NOT EXISTS public.ticket
(
    "ticket_ID" integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    "time" time with time zone NOT NULL,
    "user_ID" integer NOT NULL,
    "session_ID" integer NOT NULL,
    price real NOT NULL,
    CONSTRAINT "Ticket_pkey" PRIMARY KEY ("ticket_ID"),
    CONSTRAINT "FK_<Ticket>_<Session>" FOREIGN KEY ("session_ID")
        REFERENCES public.session ("session_ID") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT "FK_<Ticket>_<User>" FOREIGN KEY ("user_ID")
        REFERENCES public."user" ("user_ID") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

-- +goose Down
DROP TABLE public."ticket";
