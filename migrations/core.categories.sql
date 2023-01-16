CREATE TABLE core.categories
(
    id             SERIAL
        CONSTRAINT categories_pk
            PRIMARY KEY,
    user_id        BIGINT     NOT NULL,
    name           VARCHAR(64) NOT NULL,
    sanitized_name VARCHAR(64) NOT NULL,
    emoji          VARCHAR(8)
);

ALTER TABLE core.categories
    OWNER TO postgres;

CREATE INDEX categories_emoji_index
    ON core.categories (emoji);

CREATE UNIQUE INDEX categories_id_uindex
    ON core.categories (id);

CREATE INDEX categories_sanitized_name_index
    ON core.categories (sanitized_name);