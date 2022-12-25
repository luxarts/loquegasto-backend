CREATE TABLE core.wallets
(
    id             SERIAL
        CONSTRAINT wallets_pk
            PRIMARY KEY,
    user_id        INTEGER     NOT NULL,
    name           VARCHAR(64) NOT NULL,
    sanitized_name VARCHAR(64) NOT NULL,
    balance        INTEGER     NOT NULL,
    created_at     TIMESTAMP   NOT NULL
);

ALTER TABLE core.wallets
    OWNER TO postgres;

CREATE UNIQUE INDEX wallets_id_uindex
    ON core.wallets (id);

CREATE  INDEX wallets_sanitized_name_index
    ON core.wallets (sanitized_name);

CREATE INDEX wallets_user_id_index
    ON core.wallets (user_id);

