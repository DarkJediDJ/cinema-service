-- +goose Up
CREATE TABLE IF NOT EXISTS public.tickets
(
    "time" time with time zone NOT NULL,
    price real NOT NULL,
    user_id integer NOT NULL,
    ticket_id integer NOT NULL DEFAULT nextval('tickets_ticket_id_seq'::regclass),
    session_id integer NOT NULL,
    CONSTRAINT tickets_pkey PRIMARY KEY (ticket_id),
    CONSTRAINT "FK_tickets_to_session" FOREIGN KEY (session_id)
        REFERENCES public.sessions (session_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT "FK_tickets_to_users" FOREIGN KEY (user_id)
        REFERENCES public.users (user_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
)


-- +goose Down
DROP TABLE public."ticket";