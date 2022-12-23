create table wallets
(
    id             serial
        constraint wallets_pk
            primary key,
    user_id        integer     not null
        constraint wallets_users_id_fk
            references core.users,
    name           varchar(64) not null,
    sanitized_name varchar(64) not null,
    balance        integer     not null,
    created_at     timestamp   not null
);

alter table wallets
    owner to pgroot;

create unique index wallets_id_uindex
    on wallets (id);

create index wallets_sanitized_name_index
    on wallets (sanitized_name);

create index wallets_user_id_index
    on wallets (user_id);

