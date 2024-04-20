create table core.categories
(
    id             uuid        not null
        constraint categories_pk
            primary key,
    user_id        uuid        not null
        constraint categories_users_id_fk
            references core.users,
    name           varchar(64) not null,
    sanitized_name varchar(64) not null,
    emoji          varchar(8),
    created_at     timestamp   not null
);

alter table core.categories
    owner to postgres;

create index categories_sanitized_name_index
    on core.categories (sanitized_name);

create index categories_user_id_index
    on core.categories (user_id);

