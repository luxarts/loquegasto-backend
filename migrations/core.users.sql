CREATE TYPE states AS ENUM ('wallet_selection');

CREATE TABLE core.users
(
    id         INTEGER NOT NULL
        CONSTRAINT users_pk
            PRIMARY KEY,
    chat_id    INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    state states
);

ALTER TABLE core.users
    owner TO postgres;

CREATE UNIQUE INDEX users_chat_id_uindex
    ON core.users (chat_id);

CREATE UNIQUE INDEX users_id_uindex
    ON core.users (id);