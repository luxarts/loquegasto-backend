CREATE TABLE core.transactions
(
    uuid        VARCHAR(64) NOT NULL
        CONSTRAINT transactions_pk
            PRIMARY KEY,
    user_id     INTEGER     NOT NULL,
    msg_id      INTEGER     NOT NULL,
    amount      INTEGER     NOT NULL,
    description VARCHAR(64) NOT NULL,
    wallet_id   INTEGER     NOT NULL,
    created_at  TIMESTAMP   NOT NULL,
    category_id INTEGER
);

ALTER TABLE core.transactions
    OWNER TO postgres;

CREATE INDEX transactions_category_id_index
    ON core.transactions (category_id);

CREATE UNIQUE INDEX transactions_uuid_uindex
    ON core.transactions (uuid);

CREATE INDEX transactions_wallet_id_index
    ON core.transactions (wallet_id);

CREATE UNIQUE INDEX transactions_msg_id_uindex
    ON core.transactions (msg_id);