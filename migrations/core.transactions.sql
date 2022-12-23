create table transactions
(
    uuid        varchar(64) not null
        constraint transactions_pk
            primary key,
    user_id     integer     not null
        constraint transactions_users_id_fk
            references core.users,
    msg_id      integer     not null,
    amount      integer     not null,
    description varchar(64) not null,
    wallet_id   integer     not null
        constraint transactions_wallets_id_fk
            references core.wallets,
    created_at  timestamp   not null,
    category_id integer
        constraint transactions_categories_id_fk
            references core.categories
);

alter table transactions
    owner to pgroot;

create index transactions_category_id_index
    on transactions (category_id);

create unique index transactions_uuid_uindex
    on transactions (uuid);

create index transactions_wallet_id_index
    on transactions (wallet_id);

create unique index transactions_msg_id_uindex
    on transactions (msg_id);

