-- +migrate Up
CREATE TABLE pokemon (
    id integer NOT NULL,
    pokedex_id integer NOT NULL,
    name character varying(255) NOT NULL,
    type1 character varying(30) NOT NULL,
    type2 character varying(30),
    total integer,
    hp integer,
    speed integer,
    generation integer,
    legendary boolean NOT NULL
);

CREATE SEQUENCE pokemon_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE ONLY pokemon ALTER COLUMN id SET DEFAULT nextval('pokemon_id_seq'::regclass);

ALTER TABLE ONLY pokemon ADD CONSTRAINT pokemon_pkey PRIMARY KEY (id);

-- +migrate Down

DROP TABLE pokemon;
DROP SEQUENCE pokemon_id_seq;