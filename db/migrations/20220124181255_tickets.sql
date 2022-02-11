-- +goose Up
CREATE TABLE IF NOT EXISTS public.tickets
(
    "time" time with time zone NOT NULL,
    price real NOT NULL,
    user_id integer NOT NULL,
    id SERIAL,
    session_id integer NOT NULL,
    CONSTRAINT tickets_pkey PRIMARY KEY (id),
    CONSTRAINT "FK_tickets_to_session" FOREIGN KEY (session_id)
        REFERENCES public.sessions (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT "FK_tickets_to_users" FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);


-- +goose Down
DROP TABLE public.tickets;
