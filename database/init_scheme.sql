CREATE DATABASE tstkbot OWNER tstkbot;

CREATE SCHEMA tstkbot AUTHORIZATION tstkbot;

CREATE SEQUENCE tstkbot.user_id_seq
    INCREMENT 1
    START 3
    MINVALUE 1
    MAXVALUE 9223372036854775807
    CACHE 1;

ALTER SEQUENCE tstkbot.user_id_seq
    OWNER TO tstkbot;

CREATE SEQUENCE tstkbot.judge_id_seq
    INCREMENT 1
    START 2
    MINVALUE 1
    MAXVALUE 9223372036854775807
    CACHE 1;

ALTER SEQUENCE tstkbot.judge_id_seq
    OWNER TO tstkbot;

CREATE SEQUENCE tstkbot.vote_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    CACHE 1;

ALTER SEQUENCE tstkbot.vote_id_seq
    OWNER TO tstkbot;

-- Table: tstkbot.user

-- DROP TABLE tstkbot.user;

CREATE TABLE tstkbot.user
(
    id integer NOT NULL DEFAULT nextval('tstkbot.user_id_seq'::regclass),
    telegram_id integer NOT NULL,
    CONSTRAINT user_pkey PRIMARY KEY (id)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE tstkbot.user
    OWNER to tstkbot;

-- Table: tstkbot.judge

-- DROP TABLE tstkbot.judge;

CREATE TABLE tstkbot.judge
(
    id integer NOT NULL DEFAULT nextval('tstkbot.judge_id_seq'::regclass),
    phrase character varying(255) NOT NULL,
    author_id integer NOT NULL,
    CONSTRAINT judge_pkey PRIMARY KEY (id),
    CONSTRAINT phrase_unique UNIQUE (phrase),
    CONSTRAINT judge_user_id_fk FOREIGN KEY (author_id)
        REFERENCES tstkbot.user (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE tstkbot.judge
    OWNER to tstkbot;

-- Table: tstkbot.vote

-- DROP TABLE tstkbot.vote;

CREATE TABLE tstkbot.vote
(
    id integer NOT NULL DEFAULT nextval('tstkbot.vote_id_seq'::regclass),
    judge_id integer NOT NULL,
    user_id integer NOT NULL,
    CONSTRAINT vote_pkey PRIMARY KEY (id),
    CONSTRAINT vote_judge_id_fk FOREIGN KEY (judge_id)
        REFERENCES tstkbot.judge (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT vote_user_id_fk FOREIGN KEY (user_id)
        REFERENCES tstkbot.user (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE tstkbot.vote
    OWNER to tstkbot;