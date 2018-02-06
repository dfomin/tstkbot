CREATE DATABASE tstkbot OWNER tstkbot;

CREATE SCHEMA tstkbot AUTHORIZATION tstkbot;

-- Table: tstkbot.user

-- DROP TABLE tstkbot.user;

CREATE TABLE tstkbot.user
(
    id serial primary key,
    telegram_id integer NOT NULL
)
TABLESPACE pg_default;

ALTER TABLE tstkbot.user OWNER to tstkbot;

-- Table: tstkbot.judge

-- DROP TABLE tstkbot.judge;

CREATE TABLE tstkbot.judge
(
    id serial primary key,
    phrase character varying(255) unique NOT NULL,
    author_id integer NOT NULL,
    CONSTRAINT judge_user_id_fk FOREIGN KEY (author_id)
        REFERENCES tstkbot.user (id))
TABLESPACE pg_default;

ALTER TABLE tstkbot.judge OWNER to tstkbot;

-- Table: tstkbot.vote

-- DROP TABLE tstkbot.vote;

CREATE TABLE tstkbot.vote
(
    id serial primary key,
    judge_id integer NOT NULL,
    user_id integer NOT NULL,
    CONSTRAINT vote_judge_id_fk FOREIGN KEY (judge_id)
        REFERENCES tstkbot.judge (id),
    CONSTRAINT vote_user_id_fk FOREIGN KEY (user_id)
        REFERENCES tstkbot.user (id))
TABLESPACE pg_default;

ALTER TABLE tstkbot.vote OWNER to tstkbot;