create table categories
(
    id             serial
        constraint categories_pk
            primary key,
    user_id        integer     not null
        constraint categories_users_id_fk
            references core.users,
    name           varchar(64) not null,
    sanitized_name varchar(64) not null,
    emoji          varchar(8)
);

alter table categories
    owner to pgroot;

create index categories_emoji_index
    on categories (emoji);

create unique index categories_id_uindex
    on categories (id);

create index categories_sanitized_name_index
    on categories (sanitized_name);