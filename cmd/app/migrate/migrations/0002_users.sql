-- +goose Up
create table users
(
    id           uuid                  default uuid_generate_v4() not null primary key,
    username     varchar(255) not null unique,
    external_id  varchar(255) null     default null,
    name         jsonb        null     default null,
    display_name varchar(255) null     default null,
    --nickname varchar(255) null default null,
    --profile_url varchar(255) null default null,
    --title varchar(255) null default null,
    --user_type varchar(255) null default null,
    --preferred_language varchar(255) null default null,
    locale       varchar(255) null     default null,
    --timezone varchar(255) null default null,
    active       boolean      not null default true,
    emails       jsonb        null     default null,
    --phone_numbers jsonb not null default '[]',
    --ims jsonb not null default '[]',
    --addresses jsonb not null default '[]',
    --photos jsonb not null default '[]',
    --roles jsonb not null default '[]',
    --entitlements jsonb not null default '[]',
    --x509_certificates jsonb not null default '[]',
    created_at   timestamp    not null default now(),
    updated_at   timestamp    not null default now()
);

create table groups
(
    id           uuid                  default uuid_generate_v4() not null primary key,
    display_name varchar(255) not null,
    created_at   timestamp    not null default now(),
    updated_at   timestamp    not null default now()
);

create table group_members
(
    id         serial primary key,
    group_id   uuid      not null references groups (id) on delete cascade,
    user_id    uuid      not null references users (id) on delete cascade,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    unique (group_id, user_id)
);

create table api_keys
(
    id          uuid                  default uuid_generate_v4() not null primary key,
    encodedHash varchar(255) not null unique,
    created_at  timestamp    not null default now(),
    updated_at  timestamp    not null default now()
);

create table scim_api_keys
(
    id         uuid                  default uuid_generate_v4() not null primary key,
    domain     varchar(255) not null unique,
    api_key_id uuid         not null,
    created_at timestamp    not null default now(),
    updated_at timestamp    not null default now(),
    FOREIGN KEY (api_key_id) references api_keys (id) on delete cascade
);

-- +goose Down
drop table scim_api_keys;
drop table api_keys;
drop table group_members;
drop table groups;
drop table users;
