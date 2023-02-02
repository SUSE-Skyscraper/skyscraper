--------------------------------------------------------------------------------------------------------------------
-- Cloud Tenants
--------------------------------------------------------------------------------------------------------------------

-- name: CreateOrUpdateCloudTenant :one
insert into cloud_tenants (cloud, tenant_id, name)
values ($1, $2, $3)
on conflict (cloud, tenant_id) do update set name       = COALESCE(nullif($3, ''), cloud_tenants.name),
                                             updated_at = now()
returning *;

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
-- Cloud Accounts
--------------------------------------------------------------------------------------------------------------------

-- name: SearchTag :many
select *
from cloud_accounts
where cloud = $1
  and tenant_id = $2
  and tags_current ->> sqlc.arg(tag_key) = sqcl.arg(tag_value);

-- name: CreateOrUpdateCloudAccount :one
insert into cloud_accounts (cloud, tenant_id, account_id, name, tags_current, tags_desired, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, COALESCE(nullif(@tags_desired::jsonb, '{}'::jsonb), $5), now(), now())
ON CONFLICT (cloud, tenant_id, account_id)
    DO UPDATE SET name         = COALESCE(nullif($4, ''), cloud_accounts.name),
                  tags_current = COALESCE(nullif($5::jsonb, '{}'::jsonb), cloud_accounts.tags_current),
                  tags_desired = COALESCE(nullif(@tags_desired::jsonb, '{}'::jsonb), cloud_accounts.tags_desired),
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
select group_users.*, users.username as username
from group_users
         left join users on users.id = group_users.user_id
where group_users.group_id = $1;

-- name: GetGroupMembershipForUser :one
select group_users.*, users.username as username
from group_users
         left join users on users.id = group_users.user_id
where group_users.group_id = $1
  and group_users.user_id = $2;

-- name: DropMembershipForGroup :exec
delete
from group_users
where group_id = $1;

-- name: DropMembershipForUserAndGroup :exec
delete
from group_users
where user_id = $1
  and group_id = $2;

-- name: CreateMembershipForUserAndGroup :exec
insert into group_users (user_id, group_id, created_at, updated_at)
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
where id = $1
  and system = false;

-- name: FindAPIKeysById :many
select *
from api_keys
where id = ANY ($1::uuid[]);

-- name: FindScimAPIKey :one
select api_keys.*
from api_keys
         left join scim_api_keys on scim_api_keys.api_key_id = api_keys.id
where scim_api_keys.domain = 'default'
  and api_keys.system = true;

-- name: GetAPIKeys :many
select *
from api_keys
where system = false;

--------------------------------------------------------------------------------------------------------------------
-- Tags
--------------------------------------------------------------------------------------------------------------------

-- name: GetTags :many
select *
from standard_tags
order by key;

-- name: CreateTag :one
insert into standard_tags (display_name, key, description, created_at, updated_at)
values ($1, $2, $3, now(), now())
returning *;

-- name: FindTag :one
select *
from standard_tags
where id = $1;

-- name: UpdateTag :exec
update standard_tags
set display_name = $2,
    key          = $3,
    description  = $4,
    updated_at   = now()
where id = $1;

-- name: DeleteTag :exec
delete
from standard_tags
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

--------------------------------------------------------------------------------------------------------------------
-- Organizational Units
--------------------------------------------------------------------------------------------------------------------

-- name: GetOrganizationalUnitChildren :many
select *
from organizational_units
where parent_id = $1;

-- name: GetOrganizationalUnitCloudAccounts :many
select cloud_accounts.*
from organizational_units_cloud_accounts
         join cloud_accounts on cloud_accounts.id = organizational_units_cloud_accounts.cloud_account_id
where organizational_unit_id = $1;

-- name: DeleteOrganizationalUnit :exec
delete
from organizational_units
where id = $1;

-- name: FindOrganizationalUnit :one
select *
from organizational_units
where id = $1;

-- name: GetOrganizationalUnits :many
select *
from organizational_units;

-- name: UnAssignAccountFromOUs :exec
delete
from organizational_units_cloud_accounts
where cloud_account_id = $1;

-- name: AssignAccountToOU :exec
insert into organizational_units_cloud_accounts (cloud_account_id, organizational_unit_id)
values ($1, $2);

-- name: CreateOrganizationalUnit :one
insert into organizational_units (parent_id, display_name, created_at, updated_at)
values ($1, $2, now(), now())
returning *;

-- name: GetAPIKeysOrganizationalUnits :many
select o2.*
from (select group_id
      from group_api_keys
      where group_api_keys.api_key_id = $1) m
         inner join (select organizational_unit_id, group_id
                     from organizational_units_groups
                              inner join organizational_units on organizational_units.id =
                                                                 organizational_units_groups.organizational_unit_id) o
                    on m.group_id = o.group_id
         inner join (select *
                     from organizational_units) o2 on o.organizational_unit_id = o2.id;

-- name: GetUserOrganizationalUnits :many
select o2.*
from (select group_id
      from group_users
      where group_users.user_id = $1) m
         inner join (select organizational_unit_id, group_id
                     from organizational_units_groups
                              inner join organizational_units on organizational_units.id =
                                                                 organizational_units_groups.organizational_unit_id) o
                    on m.group_id = o.group_id
         inner join (select *
                     from organizational_units) o2 on o.organizational_unit_id = o2.id;

-- name: OrganizationalUnitsCloudAccounts :many
WITH RECURSIVE organizational_unit_ids(id, parent_id, display_name) AS (SELECT id,
                                                                               parent_id,
                                                                               display_name
                                                                        FROM organizational_units
                                                                        WHERE id = ANY ($1::uuid[])
                                                                        UNION ALL
                                                                        SELECT o2.id,
                                                                               o2.parent_id,
                                                                               o2.display_name
                                                                        FROM organizational_units o2
                                                                                 INNER JOIN organizational_unit_ids o ON o.id = o2.parent_id)
SELECT cloud_accounts.*
FROM organizational_unit_ids
         join organizational_units_cloud_accounts
              on organizational_unit_ids.id = organizational_units_cloud_accounts.organizational_unit_id
         join cloud_accounts on organizational_units_cloud_accounts.cloud_account_id = cloud_accounts.id;
