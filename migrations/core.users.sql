CREATE TABLE core.users
(
    id         BIGINT NOT NULL
        CONSTRAINT users_pk
            PRIMARY KEY,
    chat_id    BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    timezone_offset INTEGER
);

ALTER TABLE core.users
    owner TO postgres;

CREATE UNIQUE INDEX users_chat_id_uindex
    ON core.users (chat_id);

CREATE UNIQUE INDEX users_id_uindex
    ON core.users (id);