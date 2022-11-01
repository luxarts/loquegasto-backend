create table core.users
(
    id         integer   not null
        constraint users_pk
            primary key,
    chat_id    integer   not null,
    created_at timestamp not null
);

alter table core.users
    owner to pgroot;

create unique index users_chat_id_uindex
    on core.users (chat_id);

create unique index users_id_uindex
    on core.users (id);