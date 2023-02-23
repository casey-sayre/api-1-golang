-- Table: public.albums

-- DROP TABLE IF EXISTS public.albums;

CREATE TABLE IF NOT EXISTS public.albums
(
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    title character varying(50) COLLATE pg_catalog."default" NOT NULL,
    artist character varying(100) COLLATE pg_catalog."default" NOT NULL,
    price numeric NOT NULL,
    CONSTRAINT albums_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.albums
    OWNER to postgres;
    