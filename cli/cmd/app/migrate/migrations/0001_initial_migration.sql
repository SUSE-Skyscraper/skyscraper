-- +goose Up
create extension "uuid-ossp";

create table cloud_tenants
(
    id         uuid                  default uuid_generate_v4() not null primary key,
    cloud      varchar(255) not null,
    tenant_id  varchar(255) not null,
    name       varchar(255) not null,
    active     boolean      not null default true,
    created_at timestamp    not null default current_timestamp,
    updated_at timestamp    not null default current_timestamp,
    unique (cloud, tenant_id)
);

create table cloud_accounts
(
    id                  uuid                  default uuid_generate_v4() not null primary key,
    cloud               varchar(255) not null,
    tenant_id           varchar(255) not null,
    account_id          varchar(255) not null,
    name                varchar(255) not null,
    active              boolean      not null default true,
    tags_current        jsonb        not null default '{}',
    tags_desired        jsonb        not null default '{}',
    tags_drift_detected boolean      not null default false,
    created_at          timestamp    not null default current_timestamp,
    updated_at          timestamp    not null default current_timestamp,
    unique (cloud, tenant_id, account_id),
    foreign key (cloud, tenant_id) references cloud_tenants (cloud, tenant_id) on delete cascade
);

-- +goose Down
drop table cloud_accounts;
drop table cloud_tenants;

drop extension "uuid-ossp";
