-- +goose Up
CREATE TABLE IF NOT EXISTS public.hall
(
    "hall_ID" integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    "VIP" boolean NOT NULL,
    CONSTRAINT "Hall_pkey" PRIMARY KEY ("hall_ID")
)

-- +goose Down
DROP TABLE public."hall";

