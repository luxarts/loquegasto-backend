create table core.wallets
(
    id             uuid        not null
        constraint wallets_pk
            primary key,
    user_id        uuid        not null
        constraint wallets_users_id_fk
            references core.users,
    name           varchar(64) not null,
    sanitized_name varchar(64) not null,
    balance        bigint      not null,
    emoji          varchar(8),
    created_at     timestamp   not null,
    deleted bool
);

alter table core.wallets
    owner to postgres;

create index wallets_sanitized_name_index
    on core.wallets (sanitized_name);

create index wallets_user_id_index
    on core.wallets (user_id);

