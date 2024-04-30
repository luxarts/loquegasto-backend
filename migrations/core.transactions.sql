create table core.transactions
(
    uuid        uuid not null
        constraint transactions_pk
            primary key,
    user_id     uuid        not null
        constraint transactions_users_id_fk
            references core.users,
    msg_id      bigint,
    amount      bigint      not null,
    description varchar(64) not null,
    wallet_id   uuid     not null,
    category_id uuid not null,
    created_at  timestamp   not null
);

ALTER TABLE core.transactions
    OWNER TO postgres;

create index transactions_user_id_index
    on core.transactions (user_id);

create index transactions_wallet_id_index
    on core.transactions (wallet_id);

create index transactions_category_id_index
    on core.transactions (category_id);