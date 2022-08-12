--------------------------------------------------------------------------------------------------------------------
-- Cloud Tenants
--------------------------------------------------------------------------------------------------------------------

-- name: CreateCloudTenant :exec
insert into cloud_tenants (cloud, tenant_id, name)
values ($1, $2, $3)
on conflict (cloud, tenant_id) do update set name       = $3,
                                             updated_at = now();

-- name: GetCloudTenants :many
select *
from cloud_tenants
order by cloud, tenant_id;

-- name: GetCloudTenant :one
select *
from cloud_tenants
where cloud = $1
  and tenant_id = $2;

--------------------------------------------------------------------------------------------------------------------
-- Cloud Account Metadata
--------------------------------------------------------------------------------------------------------------------

-- name: SearchTag :many
select *
from cloud_accounts
where cloud = $1
  and tenant_id = $2
  and tags_current ->> sqlc.arg(tag_key) = sqcl.arg(tag_value);

-- name: CreateOrInsertCloudAccount :one
insert into cloud_accounts (cloud, tenant_id, account_id, name, tags_current, tags_desired)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (cloud, tenant_id, account_id)
    DO UPDATE SET name         = $4,
                  tags_current = $5,
                  updated_at   = now()
returning *;

-- name: UpdateCloudAccountTagsDriftDetected :exec
update cloud_accounts
set tags_drift_detected = $1,
    updated_at          = now()
where cloud = $2
  and tenant_id = $3
  and account_id = $4;

-- name: UpdateCloudAccount :exec
update cloud_accounts
set tags_desired = $2,
    updated_at   = now()
where id = $1;

-- name: FindCloudAccount :one
select *
from cloud_accounts
where id = $1;

-- name: FindCloudAccountByCloudAndTenant :one
select *
from cloud_accounts
where cloud = $1
  and tenant_id = $2
  and account_id = $3;

--------------------------------------------------------------------------------------------------------------------
-- Users
--------------------------------------------------------------------------------------------------------------------

-- name: GetUsers :many
select *
from users
order by created_at
LIMIT $1 OFFSET $2;

-- name: GetUsersById :many
select *
from users
where id = ANY ($1::uuid[])
order by display_name;

-- name: GetUser :one
select *
from users
where id = $1;

-- name: FindByUsername :one
select *
from users
where username = $1;

-- name: CreateUser :one
insert into users (username, name, display_name, emails, active, locale, external_id, created_at, updated_at)
values ($1, $2, $3, $4, $5, $6, $7, now(), now())
returning *;

-- name: UpdateUser :exec
update users
set username     =$2,
    name         = $3,
    display_name = $4,
    emails       = $5,
    active       = $6,
    external_id  = $7,
    locale       = $8,
    updated_at   = now()
where id = $1;

-- name: PatchUser :exec
update users
set active     = $2,
    updated_at = now()
where id = $1;

-- name: DeleteUser :exec
delete
from users
where id = $1;

-- name: GetUserCount :one
select count(*)
from users;

--------------------------------------------------------------------------------------------------------------------
-- Groups
--------------------------------------------------------------------------------------------------------------------

-- name: GetGroups :many
select *
from groups
order by id
LIMIT $1 OFFSET $2;

-- name: GetGroup :one
select *
from groups
where id = $1;

-- name: CreateGroup :one
insert into groups (display_name, created_at, updated_at)
values ($1, now(), now())
returning *;

-- name: DeleteGroup :exec
delete
from groups
where id = $1;

-- name: GetGroupCount :one
select count(*)
from groups;

-- name: PatchGroupDisplayName :exec
update groups
set display_name = $2,
    updated_at   = now()
where id = $1;

--------------------------------------------------------------------------------------------------------------------
-- Membership
--------------------------------------------------------------------------------------------------------------------

-- name: GetGroupMembership :many
select group_members.*, users.username as username
from group_members
         left join users on users.id = group_members.user_id
where group_members.group_id = $1;

-- name: GetGroupMembershipForUser :one
select group_members.*, users.username as username
from group_members
         left join users on users.id = group_members.user_id
where group_members.group_id = $1
  and group_members.user_id = $2;

-- name: DropMembershipForGroup :exec
delete
from group_members
where group_id = $1;

-- name: DropMembershipForUserAndGroup :exec
delete
from group_members
where user_id = $1
  and group_id = $2;

-- name: CreateMembershipForUserAndGroup :exec
insert into group_members (user_id, group_id, created_at, updated_at)
values ($1, $2, now(), now())
on conflict (user_id, group_id) do update set updated_at = now();;

--------------------------------------------------------------------------------------------------------------------
-- SCIM API Key
--------------------------------------------------------------------------------------------------------------------

-- name: InsertAPIKey :one
insert into api_keys (encodedhash, system, owner, description, created_at, updated_at)
values ($1, $2, $3, $4, now(), now())
returning *;

-- name: InsertScimAPIKey :one
insert into scim_api_keys (api_key_id, domain, created_at, updated_at)
values ($1, 'default', now(), now())
returning *;

-- name: DeleteAPIKey :exec
delete
from api_keys
where id = $1;

-- name: DeleteScimAPIKey :exec
delete
from scim_api_keys
where domain = 'default';

-- name: FindAPIKey :one
select *
from api_keys
where id = $1 and system = false;

-- name: FindAPIKeysById :many
select *
from api_keys
where id = ANY ($1::uuid[]);

-- name: FindScimAPIKey :one
select api_keys.*
from api_keys
         left join scim_api_keys on scim_api_keys.api_key_id = api_keys.id
where scim_api_keys.domain = 'default' and api_keys.system = true;

-- name: GetAPIKeys :many
select *
from api_keys
where system = false;

--------------------------------------------------------------------------------------------------------------------
-- Tags
--------------------------------------------------------------------------------------------------------------------

-- name: GetTags :many
select *
from tags
order by key;

-- name: CreateTag :one
insert into tags (display_name, key, description, created_at, updated_at)
values ($1, $2, $3, now(), now())
returning *;

-- name: FindTag :one
select *
from tags
where id = $1;

-- name: UpdateTag :exec
update tags
set display_name = $2,
    key          = $3,
    description  = $4,
    updated_at   = now()
where id = $1;

-- name: DeleteTag :exec
delete
from tags
where id = $1;

--------------------------------------------------------------------------------------------------------------------
-- Audit Logs
--------------------------------------------------------------------------------------------------------------------

-- name: GetAuditLogs :many
select *
from audit_logs
order by created_at desc;

-- name: GetAuditLogsForTarget :many
select audit_logs.*
from audit_logs
where resource_id = $1
  and resource_type = $2
order by created_at desc;

-- name: CreateAuditLog :one
insert into audit_logs (resource_type, resource_id, caller_id, caller_type, message, created_at, updated_at)
values ($1, $2, $3, $4, $5, now(), now())
returning *;
