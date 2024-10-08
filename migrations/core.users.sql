create table core.users
(
    id              uuid      not null
        constraint users_pk
            primary key,
    chat_id         bigint    not null
        constraint chat_id
            unique,
    timezone_offset integer   not null,
    created_at      timestamp not null
);

alter table core.users
    owner to postgres;

