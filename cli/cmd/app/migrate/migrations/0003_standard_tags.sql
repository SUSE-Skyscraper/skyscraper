-- +goose Up
create table standard_tags
(
    id           uuid                  default uuid_generate_v4() not null primary key,
    display_name varchar(255) not null,
    description  varchar(255) not null,
    key          varchar(255) not null unique,
    created_at   timestamp    not null default now(),
    updated_at   timestamp    not null default now()
);

-- +goose Down
drop table standard_tags;
