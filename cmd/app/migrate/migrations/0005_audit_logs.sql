-- +goose Up

create type audit_resource_type as enum ('cloud_account','tag', 'policy', 'cloud_tenant', 'user', 'group', 'scim_api_key');
create table audit_logs
(
    id          uuid                       default uuid_generate_v4() not null primary key,
    user_id     uuid              not null,
    resource_type audit_resource_type not null,
    resource_id   uuid              not null,
    message     text              not null,
    created_at  timestamp         not null default now(),
    updated_at  timestamp         not null default now(),
    FOREIGN KEY (user_id) REFERENCES users (id)
);
create index idx_audit_logs_target_type_target_id on audit_logs (resource_type, resource_id);

-- +goose Down
drop table audit_logs;
drop type audit_resource_type;
