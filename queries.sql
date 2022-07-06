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
set tags_desired = $4,
    updated_at   = now()
where cloud = $1
  and tenant_id = $2
  and account_id = $3;

-- name: GetCloudAllAccounts :many
select *
from cloud_accounts
order by cloud, tenant_id, account_id;

-- name: GetCloudAllAccountsForCloud :many
select *
from cloud_accounts
where cloud = $1
order by tenant_id, account_id;

-- name: GetCloudAllAccountsForCloudAndTenant :many
select *
from cloud_accounts
where cloud = $1
  and tenant_id = $2
order by account_id;

-- name: GetCloudAccount :one
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
insert into users (username, name, emails, active, created_at, updated_at)
values ($1, $2, $3, $4, now(), now())
returning *;

-- name: UpdateUser :exec
update users
set username   =$2,
    name       = $3,
    emails     = $4,
    active     = $5,
    updated_at = now()
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
-- Users
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

--------------------------------------------------------------------------------------------------------------------
-- User Membership
--------------------------------------------------------------------------------------------------------------------
