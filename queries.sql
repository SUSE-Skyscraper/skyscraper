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
order by id
LIMIT $1 OFFSET $2;

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
insert into scim_api_keys (token, created_at, updated_at)
values ($1, now(), now())
returning *;

-- name: FindAPIKey :one
select *
from scim_api_keys
where token = $1;

--------------------------------------------------------------------------------------------------------------------
-- Policies
--
-- 6ba7b812-9dad-11d1-80b4-00c04fd430c8 is NameSpace_OID as specified in rfc4122 (https://tools.ietf.org/html/rfc4122)
-- we use uuid v5 so we can calculate the id from a collection of values
--------------------------------------------------------------------------------------------------------------------

-- name: GetPolicies :many
select *
from policies
order by id;

-- name: AddPolicy :exec
insert into policies (id, ptype, v0, v1, v2, v3, v4, v5)
values (uuid_generate_v5('6ba7b812-9dad-11d1-80b4-00c04fd430c8',
                         concat(sqlc.arg(ptype)::text, sqlc.arg(v0)::text, sqlc.arg(v1)::text, sqlc.arg(v2)::text,
                                sqlc.arg(v3)::text, sqlc.arg(v4)::text, sqlc.arg(v5)::text)), sqlc.arg(ptype)::text,
        sqlc.arg(v0)::text, sqlc.arg(v1)::text, sqlc.arg(v2)::text, sqlc.arg(v3)::text, sqlc.arg(v4)::text,
        sqlc.arg(v5)::text)
on conflict do nothing;

-- name: RemovePolicy :exec
delete
from policies
where id = uuid_generate_v5('6ba7b812-9dad-11d1-80b4-00c04fd430c8',
                            concat(sqlc.arg(ptype), sqlc.arg(v0), sqlc.arg(v1), sqlc.arg(v2), sqlc.arg(v3),
                                   sqlc.arg(v4), sqlc.arg(v5)));

-- name: RemovePoliciesForGroup :exec
delete
from policies
where ptype = 'g'
  and v1 = $1;

-- name: TruncatePolicies :exec
truncate policies;

--------------------------------------------------------------------------------------------------------------------
-- Tags
--------------------------------------------------------------------------------------------------------------------

-- name: GetTags :many
select *
from tags
order by key;

-- name: CreateTag :one
insert into tags (display_name, key, required, description, created_at, updated_at)
values ($1, $2, $3, $4, now(), now())
returning *;

-- name: FindTag :one
select *
from tags
where id = $1;

-- name: UpdateTag :exec
update tags
set display_name = $2,
    key          = $3,
    required     = $4,
    description  = $5,
    updated_at   = now()
where id = $1;

-- name: DeleteTag :exec
delete
from tags
where id = $1;
