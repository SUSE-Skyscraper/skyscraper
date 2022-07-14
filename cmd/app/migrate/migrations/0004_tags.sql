-- +goose Up
create table tags
(
    id           uuid                  default uuid_generate_v4() not null primary key,
    display_name varchar(255) not null,
    description  varchar(255) not null,
    required     boolean      not null default false,
    key          varchar(255) not null unique,
    overrides    jsonb        not null default '{}',
    created_at   timestamp    not null default now(),
    updated_at   timestamp    not null default now()
);

-- +goose Down
drop table tags;
