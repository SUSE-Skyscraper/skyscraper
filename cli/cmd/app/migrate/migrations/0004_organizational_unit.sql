-- +goose Up
create table organizational_units
(
    id           uuid                                                                 default uuid_generate_v4() not null primary key,
    parent_id    uuid references organizational_units (id) on delete cascade null,
    display_name varchar(255)                                                not null,
    created_at   timestamp                                                   not null default now(),
    updated_at   timestamp                                                   not null default now()
);

create table organizational_units_cloud_accounts
(
    cloud_account_id uuid references cloud_accounts (id) on delete cascade not null unique,
    organizational_unit_id uuid references organizational_units (id) on delete cascade not null
);

create table organizational_units_groups
(
    group_id uuid references groups (id) on delete cascade not null unique,
    organizational_unit_id uuid references organizational_units (id) on delete cascade not null
);

-- +goose Down
drop table organizational_units_groups;
drop table organizational_units_cloud_accounts;
drop table organizational_units;
