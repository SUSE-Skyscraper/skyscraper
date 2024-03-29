// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: queries.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

const assignAccountToOU = `-- name: AssignAccountToOU :exec
insert into organizational_units_cloud_accounts (cloud_account_id, organizational_unit_id)
values ($1, $2)
`

type AssignAccountToOUParams struct {
	CloudAccountID       uuid.UUID
	OrganizationalUnitID uuid.UUID
}

func (q *Queries) AssignAccountToOU(ctx context.Context, arg AssignAccountToOUParams) error {
	_, err := q.db.Exec(ctx, assignAccountToOU, arg.CloudAccountID, arg.OrganizationalUnitID)
	return err
}

const createAuditLog = `-- name: CreateAuditLog :one
insert into audit_logs (resource_type, resource_id, caller_id, caller_type, message, created_at, updated_at)
values ($1, $2, $3, $4, $5, now(), now())
returning id, caller_id, caller_type, resource_type, resource_id, message, created_at, updated_at
`

type CreateAuditLogParams struct {
	ResourceType AuditResourceType
	ResourceID   uuid.UUID
	CallerID     uuid.UUID
	CallerType   CallerType
	Message      string
}

func (q *Queries) CreateAuditLog(ctx context.Context, arg CreateAuditLogParams) (AuditLog, error) {
	row := q.db.QueryRow(ctx, createAuditLog,
		arg.ResourceType,
		arg.ResourceID,
		arg.CallerID,
		arg.CallerType,
		arg.Message,
	)
	var i AuditLog
	err := row.Scan(
		&i.ID,
		&i.CallerID,
		&i.CallerType,
		&i.ResourceType,
		&i.ResourceID,
		&i.Message,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createGroup = `-- name: CreateGroup :one
insert into groups (display_name, created_at, updated_at)
values ($1, now(), now())
returning id, display_name, created_at, updated_at
`

func (q *Queries) CreateGroup(ctx context.Context, displayName string) (Group, error) {
	row := q.db.QueryRow(ctx, createGroup, displayName)
	var i Group
	err := row.Scan(
		&i.ID,
		&i.DisplayName,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createMembershipForUserAndGroup = `-- name: CreateMembershipForUserAndGroup :exec
insert into group_users (user_id, group_id, created_at, updated_at)
values ($1, $2, now(), now())
on conflict (user_id, group_id) do update set updated_at = now()
`

type CreateMembershipForUserAndGroupParams struct {
	UserID  uuid.UUID
	GroupID uuid.UUID
}

func (q *Queries) CreateMembershipForUserAndGroup(ctx context.Context, arg CreateMembershipForUserAndGroupParams) error {
	_, err := q.db.Exec(ctx, createMembershipForUserAndGroup, arg.UserID, arg.GroupID)
	return err
}

const createOrUpdateCloudAccount = `-- name: CreateOrUpdateCloudAccount :one
insert into cloud_accounts (cloud, tenant_id, account_id, name, tags_current, tags_desired, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, COALESCE(nullif($6::jsonb, '{}'::jsonb), $5), now(), now())
ON CONFLICT (cloud, tenant_id, account_id)
    DO UPDATE SET name         = COALESCE(nullif($4, ''), cloud_accounts.name),
                  tags_current = COALESCE(nullif($5::jsonb, '{}'::jsonb), cloud_accounts.tags_current),
                  tags_desired = COALESCE(nullif($6::jsonb, '{}'::jsonb), cloud_accounts.tags_desired),
                  updated_at   = now()
returning id, cloud, tenant_id, account_id, name, active, tags_current, tags_desired, tags_drift_detected, created_at, updated_at
`

type CreateOrUpdateCloudAccountParams struct {
	Cloud       string
	TenantID    string
	AccountID   string
	Name        string
	TagsCurrent pgtype.JSONB
	TagsDesired pgtype.JSONB
}

func (q *Queries) CreateOrUpdateCloudAccount(ctx context.Context, arg CreateOrUpdateCloudAccountParams) (CloudAccount, error) {
	row := q.db.QueryRow(ctx, createOrUpdateCloudAccount,
		arg.Cloud,
		arg.TenantID,
		arg.AccountID,
		arg.Name,
		arg.TagsCurrent,
		arg.TagsDesired,
	)
	var i CloudAccount
	err := row.Scan(
		&i.ID,
		&i.Cloud,
		&i.TenantID,
		&i.AccountID,
		&i.Name,
		&i.Active,
		&i.TagsCurrent,
		&i.TagsDesired,
		&i.TagsDriftDetected,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createOrUpdateCloudTenant = `-- name: CreateOrUpdateCloudTenant :one

insert into cloud_tenants (cloud, tenant_id, name)
values ($1, $2, $3)
on conflict (cloud, tenant_id) do update set name       = COALESCE(nullif($3, ''), cloud_tenants.name),
                                             updated_at = now()
returning id, cloud, tenant_id, name, active, created_at, updated_at
`

type CreateOrUpdateCloudTenantParams struct {
	Cloud    string
	TenantID string
	Name     string
}

// ------------------------------------------------------------------------------------------------------------------
// Cloud Tenants
// ------------------------------------------------------------------------------------------------------------------
func (q *Queries) CreateOrUpdateCloudTenant(ctx context.Context, arg CreateOrUpdateCloudTenantParams) (CloudTenant, error) {
	row := q.db.QueryRow(ctx, createOrUpdateCloudTenant, arg.Cloud, arg.TenantID, arg.Name)
	var i CloudTenant
	err := row.Scan(
		&i.ID,
		&i.Cloud,
		&i.TenantID,
		&i.Name,
		&i.Active,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createOrganizationalUnit = `-- name: CreateOrganizationalUnit :one
insert into organizational_units (parent_id, display_name, created_at, updated_at)
values ($1, $2, now(), now())
returning id, parent_id, display_name, created_at, updated_at
`

type CreateOrganizationalUnitParams struct {
	ParentID    uuid.NullUUID
	DisplayName string
}

func (q *Queries) CreateOrganizationalUnit(ctx context.Context, arg CreateOrganizationalUnitParams) (OrganizationalUnit, error) {
	row := q.db.QueryRow(ctx, createOrganizationalUnit, arg.ParentID, arg.DisplayName)
	var i OrganizationalUnit
	err := row.Scan(
		&i.ID,
		&i.ParentID,
		&i.DisplayName,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createTag = `-- name: CreateTag :one
insert into standard_tags (display_name, key, description, created_at, updated_at)
values ($1, $2, $3, now(), now())
returning id, display_name, description, key, created_at, updated_at
`

type CreateTagParams struct {
	DisplayName string
	Key         string
	Description string
}

func (q *Queries) CreateTag(ctx context.Context, arg CreateTagParams) (StandardTag, error) {
	row := q.db.QueryRow(ctx, createTag, arg.DisplayName, arg.Key, arg.Description)
	var i StandardTag
	err := row.Scan(
		&i.ID,
		&i.DisplayName,
		&i.Description,
		&i.Key,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
insert into users (username, name, display_name, emails, active, locale, external_id, created_at, updated_at)
values ($1, $2, $3, $4, $5, $6, $7, now(), now())
returning id, username, external_id, name, display_name, locale, active, emails, created_at, updated_at
`

type CreateUserParams struct {
	Username    string
	Name        pgtype.JSONB
	DisplayName sql.NullString
	Emails      pgtype.JSONB
	Active      bool
	Locale      sql.NullString
	ExternalID  sql.NullString
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Username,
		arg.Name,
		arg.DisplayName,
		arg.Emails,
		arg.Active,
		arg.Locale,
		arg.ExternalID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.ExternalID,
		&i.Name,
		&i.DisplayName,
		&i.Locale,
		&i.Active,
		&i.Emails,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteAPIKey = `-- name: DeleteAPIKey :exec
delete
from api_keys
where id = $1
`

func (q *Queries) DeleteAPIKey(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteAPIKey, id)
	return err
}

const deleteGroup = `-- name: DeleteGroup :exec
delete
from groups
where id = $1
`

func (q *Queries) DeleteGroup(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteGroup, id)
	return err
}

const deleteOrganizationalUnit = `-- name: DeleteOrganizationalUnit :exec
delete
from organizational_units
where id = $1
`

func (q *Queries) DeleteOrganizationalUnit(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteOrganizationalUnit, id)
	return err
}

const deleteScimAPIKey = `-- name: DeleteScimAPIKey :exec
delete
from scim_api_keys
where domain = 'default'
`

func (q *Queries) DeleteScimAPIKey(ctx context.Context) error {
	_, err := q.db.Exec(ctx, deleteScimAPIKey)
	return err
}

const deleteTag = `-- name: DeleteTag :exec
delete
from standard_tags
where id = $1
`

func (q *Queries) DeleteTag(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteTag, id)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
delete
from users
where id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const dropMembershipForGroup = `-- name: DropMembershipForGroup :exec
delete
from group_users
where group_id = $1
`

func (q *Queries) DropMembershipForGroup(ctx context.Context, groupID uuid.UUID) error {
	_, err := q.db.Exec(ctx, dropMembershipForGroup, groupID)
	return err
}

const dropMembershipForUserAndGroup = `-- name: DropMembershipForUserAndGroup :exec
delete
from group_users
where user_id = $1
  and group_id = $2
`

type DropMembershipForUserAndGroupParams struct {
	UserID  uuid.UUID
	GroupID uuid.UUID
}

func (q *Queries) DropMembershipForUserAndGroup(ctx context.Context, arg DropMembershipForUserAndGroupParams) error {
	_, err := q.db.Exec(ctx, dropMembershipForUserAndGroup, arg.UserID, arg.GroupID)
	return err
}

const findAPIKey = `-- name: FindAPIKey :one
select id, encodedhash, owner, description, system, created_at, updated_at
from api_keys
where id = $1
  and system = false
`

func (q *Queries) FindAPIKey(ctx context.Context, id uuid.UUID) (ApiKey, error) {
	row := q.db.QueryRow(ctx, findAPIKey, id)
	var i ApiKey
	err := row.Scan(
		&i.ID,
		&i.Encodedhash,
		&i.Owner,
		&i.Description,
		&i.System,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findAPIKeysByID = `-- name: FindAPIKeysByID :many
select id, encodedhash, owner, description, system, created_at, updated_at
from api_keys
where id = ANY ($1::uuid[])
`

func (q *Queries) FindAPIKeysByID(ctx context.Context, id []uuid.UUID) ([]ApiKey, error) {
	rows, err := q.db.Query(ctx, findAPIKeysByID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ApiKey
	for rows.Next() {
		var i ApiKey
		if err := rows.Scan(
			&i.ID,
			&i.Encodedhash,
			&i.Owner,
			&i.Description,
			&i.System,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findCloudAccount = `-- name: FindCloudAccount :one
select id, cloud, tenant_id, account_id, name, active, tags_current, tags_desired, tags_drift_detected, created_at, updated_at
from cloud_accounts
where id = $1
`

func (q *Queries) FindCloudAccount(ctx context.Context, id uuid.UUID) (CloudAccount, error) {
	row := q.db.QueryRow(ctx, findCloudAccount, id)
	var i CloudAccount
	err := row.Scan(
		&i.ID,
		&i.Cloud,
		&i.TenantID,
		&i.AccountID,
		&i.Name,
		&i.Active,
		&i.TagsCurrent,
		&i.TagsDesired,
		&i.TagsDriftDetected,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findCloudAccountByCloudAndTenant = `-- name: FindCloudAccountByCloudAndTenant :one
select id, cloud, tenant_id, account_id, name, active, tags_current, tags_desired, tags_drift_detected, created_at, updated_at
from cloud_accounts
where cloud = $1
  and tenant_id = $2
  and account_id = $3
`

type FindCloudAccountByCloudAndTenantParams struct {
	Cloud     string
	TenantID  string
	AccountID string
}

func (q *Queries) FindCloudAccountByCloudAndTenant(ctx context.Context, arg FindCloudAccountByCloudAndTenantParams) (CloudAccount, error) {
	row := q.db.QueryRow(ctx, findCloudAccountByCloudAndTenant, arg.Cloud, arg.TenantID, arg.AccountID)
	var i CloudAccount
	err := row.Scan(
		&i.ID,
		&i.Cloud,
		&i.TenantID,
		&i.AccountID,
		&i.Name,
		&i.Active,
		&i.TagsCurrent,
		&i.TagsDesired,
		&i.TagsDriftDetected,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findOrganizationalUnit = `-- name: FindOrganizationalUnit :one
select id, parent_id, display_name, created_at, updated_at
from organizational_units
where id = $1
`

func (q *Queries) FindOrganizationalUnit(ctx context.Context, id uuid.UUID) (OrganizationalUnit, error) {
	row := q.db.QueryRow(ctx, findOrganizationalUnit, id)
	var i OrganizationalUnit
	err := row.Scan(
		&i.ID,
		&i.ParentID,
		&i.DisplayName,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findScimAPIKey = `-- name: FindScimAPIKey :one
select api_keys.id, api_keys.encodedhash, api_keys.owner, api_keys.description, api_keys.system, api_keys.created_at, api_keys.updated_at
from api_keys
         left join scim_api_keys on scim_api_keys.api_key_id = api_keys.id
where scim_api_keys.domain = 'default'
  and api_keys.system = true
`

func (q *Queries) FindScimAPIKey(ctx context.Context) (ApiKey, error) {
	row := q.db.QueryRow(ctx, findScimAPIKey)
	var i ApiKey
	err := row.Scan(
		&i.ID,
		&i.Encodedhash,
		&i.Owner,
		&i.Description,
		&i.System,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findTag = `-- name: FindTag :one
select id, display_name, description, key, created_at, updated_at
from standard_tags
where id = $1
`

func (q *Queries) FindTag(ctx context.Context, id uuid.UUID) (StandardTag, error) {
	row := q.db.QueryRow(ctx, findTag, id)
	var i StandardTag
	err := row.Scan(
		&i.ID,
		&i.DisplayName,
		&i.Description,
		&i.Key,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findUserByUsername = `-- name: FindUserByUsername :one
select id, username, external_id, name, display_name, locale, active, emails, created_at, updated_at
from users
where username = $1
`

func (q *Queries) FindUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRow(ctx, findUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.ExternalID,
		&i.Name,
		&i.DisplayName,
		&i.Locale,
		&i.Active,
		&i.Emails,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAPIKeys = `-- name: GetAPIKeys :many
select id, encodedhash, owner, description, system, created_at, updated_at
from api_keys
where system = false
`

func (q *Queries) GetAPIKeys(ctx context.Context) ([]ApiKey, error) {
	rows, err := q.db.Query(ctx, getAPIKeys)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ApiKey
	for rows.Next() {
		var i ApiKey
		if err := rows.Scan(
			&i.ID,
			&i.Encodedhash,
			&i.Owner,
			&i.Description,
			&i.System,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAPIKeysOrganizationalUnits = `-- name: GetAPIKeysOrganizationalUnits :many
select o2.id, o2.parent_id, o2.display_name, o2.created_at, o2.updated_at
from (select group_id
      from group_api_keys
      where group_api_keys.api_key_id = $1) m
         inner join (select organizational_unit_id, group_id
                     from organizational_units_groups
                              inner join organizational_units on organizational_units.id =
                                                                 organizational_units_groups.organizational_unit_id) o
                    on m.group_id = o.group_id
         inner join (select id, parent_id, display_name, created_at, updated_at
                     from organizational_units) o2 on o.organizational_unit_id = o2.id
`

func (q *Queries) GetAPIKeysOrganizationalUnits(ctx context.Context, apiKeyID uuid.UUID) ([]OrganizationalUnit, error) {
	rows, err := q.db.Query(ctx, getAPIKeysOrganizationalUnits, apiKeyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []OrganizationalUnit
	for rows.Next() {
		var i OrganizationalUnit
		if err := rows.Scan(
			&i.ID,
			&i.ParentID,
			&i.DisplayName,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAuditLogs = `-- name: GetAuditLogs :many

select id, caller_id, caller_type, resource_type, resource_id, message, created_at, updated_at
from audit_logs
order by created_at desc
`

// ------------------------------------------------------------------------------------------------------------------
// Audit Logs
// ------------------------------------------------------------------------------------------------------------------
func (q *Queries) GetAuditLogs(ctx context.Context) ([]AuditLog, error) {
	rows, err := q.db.Query(ctx, getAuditLogs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AuditLog
	for rows.Next() {
		var i AuditLog
		if err := rows.Scan(
			&i.ID,
			&i.CallerID,
			&i.CallerType,
			&i.ResourceType,
			&i.ResourceID,
			&i.Message,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAuditLogsForTarget = `-- name: GetAuditLogsForTarget :many
select audit_logs.id, audit_logs.caller_id, audit_logs.caller_type, audit_logs.resource_type, audit_logs.resource_id, audit_logs.message, audit_logs.created_at, audit_logs.updated_at
from audit_logs
where resource_id = $1
  and resource_type = $2
order by created_at desc
`

type GetAuditLogsForTargetParams struct {
	ResourceID   uuid.UUID
	ResourceType AuditResourceType
}

func (q *Queries) GetAuditLogsForTarget(ctx context.Context, arg GetAuditLogsForTargetParams) ([]AuditLog, error) {
	rows, err := q.db.Query(ctx, getAuditLogsForTarget, arg.ResourceID, arg.ResourceType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AuditLog
	for rows.Next() {
		var i AuditLog
		if err := rows.Scan(
			&i.ID,
			&i.CallerID,
			&i.CallerType,
			&i.ResourceType,
			&i.ResourceID,
			&i.Message,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCloudTenant = `-- name: GetCloudTenant :one
select id, cloud, tenant_id, name, active, created_at, updated_at
from cloud_tenants
where cloud = $1
  and tenant_id = $2
`

type GetCloudTenantParams struct {
	Cloud    string
	TenantID string
}

func (q *Queries) GetCloudTenant(ctx context.Context, arg GetCloudTenantParams) (CloudTenant, error) {
	row := q.db.QueryRow(ctx, getCloudTenant, arg.Cloud, arg.TenantID)
	var i CloudTenant
	err := row.Scan(
		&i.ID,
		&i.Cloud,
		&i.TenantID,
		&i.Name,
		&i.Active,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getCloudTenants = `-- name: GetCloudTenants :many
select id, cloud, tenant_id, name, active, created_at, updated_at
from cloud_tenants
order by cloud, tenant_id
`

func (q *Queries) GetCloudTenants(ctx context.Context) ([]CloudTenant, error) {
	rows, err := q.db.Query(ctx, getCloudTenants)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CloudTenant
	for rows.Next() {
		var i CloudTenant
		if err := rows.Scan(
			&i.ID,
			&i.Cloud,
			&i.TenantID,
			&i.Name,
			&i.Active,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGroup = `-- name: GetGroup :one
select id, display_name, created_at, updated_at
from groups
where id = $1
`

func (q *Queries) GetGroup(ctx context.Context, id uuid.UUID) (Group, error) {
	row := q.db.QueryRow(ctx, getGroup, id)
	var i Group
	err := row.Scan(
		&i.ID,
		&i.DisplayName,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getGroupCount = `-- name: GetGroupCount :one
select count(*)
from groups
`

func (q *Queries) GetGroupCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, getGroupCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getGroupMembership = `-- name: GetGroupMembership :many

select group_users.group_id, group_users.user_id, users.username as username
from group_users
         left join users on users.id = group_users.user_id
where group_users.group_id = $1
`

type GetGroupMembershipRow struct {
	GroupID  uuid.UUID
	UserID   uuid.UUID
	Username sql.NullString
}

// ------------------------------------------------------------------------------------------------------------------
// Membership
// ------------------------------------------------------------------------------------------------------------------
func (q *Queries) GetGroupMembership(ctx context.Context, groupID uuid.UUID) ([]GetGroupMembershipRow, error) {
	rows, err := q.db.Query(ctx, getGroupMembership, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetGroupMembershipRow
	for rows.Next() {
		var i GetGroupMembershipRow
		if err := rows.Scan(&i.GroupID, &i.UserID, &i.Username); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGroupMembershipForUser = `-- name: GetGroupMembershipForUser :one
select group_users.group_id, group_users.user_id, users.username as username
from group_users
         left join users on users.id = group_users.user_id
where group_users.group_id = $1
  and group_users.user_id = $2
`

type GetGroupMembershipForUserParams struct {
	GroupID uuid.UUID
	UserID  uuid.UUID
}

type GetGroupMembershipForUserRow struct {
	GroupID  uuid.UUID
	UserID   uuid.UUID
	Username sql.NullString
}

func (q *Queries) GetGroupMembershipForUser(ctx context.Context, arg GetGroupMembershipForUserParams) (GetGroupMembershipForUserRow, error) {
	row := q.db.QueryRow(ctx, getGroupMembershipForUser, arg.GroupID, arg.UserID)
	var i GetGroupMembershipForUserRow
	err := row.Scan(&i.GroupID, &i.UserID, &i.Username)
	return i, err
}

const getGroups = `-- name: GetGroups :many

select id, display_name, created_at, updated_at
from groups
order by id
LIMIT $1 OFFSET $2
`

type GetGroupsParams struct {
	Limit  int32
	Offset int32
}

// ------------------------------------------------------------------------------------------------------------------
// Groups
// ------------------------------------------------------------------------------------------------------------------
func (q *Queries) GetGroups(ctx context.Context, arg GetGroupsParams) ([]Group, error) {
	rows, err := q.db.Query(ctx, getGroups, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Group
	for rows.Next() {
		var i Group
		if err := rows.Scan(
			&i.ID,
			&i.DisplayName,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getOrganizationalUnitChildren = `-- name: GetOrganizationalUnitChildren :many

select id, parent_id, display_name, created_at, updated_at
from organizational_units
where parent_id = $1::uuid
`

// ------------------------------------------------------------------------------------------------------------------
// Organizational Units
// ------------------------------------------------------------------------------------------------------------------
func (q *Queries) GetOrganizationalUnitChildren(ctx context.Context, parentID uuid.UUID) ([]OrganizationalUnit, error) {
	rows, err := q.db.Query(ctx, getOrganizationalUnitChildren, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []OrganizationalUnit
	for rows.Next() {
		var i OrganizationalUnit
		if err := rows.Scan(
			&i.ID,
			&i.ParentID,
			&i.DisplayName,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getOrganizationalUnitCloudAccounts = `-- name: GetOrganizationalUnitCloudAccounts :many
select cloud_accounts.id, cloud_accounts.cloud, cloud_accounts.tenant_id, cloud_accounts.account_id, cloud_accounts.name, cloud_accounts.active, cloud_accounts.tags_current, cloud_accounts.tags_desired, cloud_accounts.tags_drift_detected, cloud_accounts.created_at, cloud_accounts.updated_at
from organizational_units_cloud_accounts
         join cloud_accounts on cloud_accounts.id = organizational_units_cloud_accounts.cloud_account_id
where organizational_unit_id = $1
`

func (q *Queries) GetOrganizationalUnitCloudAccounts(ctx context.Context, organizationalUnitID uuid.UUID) ([]CloudAccount, error) {
	rows, err := q.db.Query(ctx, getOrganizationalUnitCloudAccounts, organizationalUnitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CloudAccount
	for rows.Next() {
		var i CloudAccount
		if err := rows.Scan(
			&i.ID,
			&i.Cloud,
			&i.TenantID,
			&i.AccountID,
			&i.Name,
			&i.Active,
			&i.TagsCurrent,
			&i.TagsDesired,
			&i.TagsDriftDetected,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getOrganizationalUnits = `-- name: GetOrganizationalUnits :many
select id, parent_id, display_name, created_at, updated_at
from organizational_units
`

func (q *Queries) GetOrganizationalUnits(ctx context.Context) ([]OrganizationalUnit, error) {
	rows, err := q.db.Query(ctx, getOrganizationalUnits)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []OrganizationalUnit
	for rows.Next() {
		var i OrganizationalUnit
		if err := rows.Scan(
			&i.ID,
			&i.ParentID,
			&i.DisplayName,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTags = `-- name: GetTags :many

select id, display_name, description, key, created_at, updated_at
from standard_tags
order by key
`

// ------------------------------------------------------------------------------------------------------------------
// Tags
// ------------------------------------------------------------------------------------------------------------------
func (q *Queries) GetTags(ctx context.Context) ([]StandardTag, error) {
	rows, err := q.db.Query(ctx, getTags)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []StandardTag
	for rows.Next() {
		var i StandardTag
		if err := rows.Scan(
			&i.ID,
			&i.DisplayName,
			&i.Description,
			&i.Key,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUser = `-- name: GetUser :one
select id, username, external_id, name, display_name, locale, active, emails, created_at, updated_at
from users
where id = $1
`

func (q *Queries) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.ExternalID,
		&i.Name,
		&i.DisplayName,
		&i.Locale,
		&i.Active,
		&i.Emails,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserCount = `-- name: GetUserCount :one
select count(*)
from users
`

func (q *Queries) GetUserCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, getUserCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getUserOrganizationalUnits = `-- name: GetUserOrganizationalUnits :many
select o2.id, o2.parent_id, o2.display_name, o2.created_at, o2.updated_at
from (select group_id
      from group_users
      where group_users.user_id = $1) m
         inner join (select organizational_unit_id, group_id
                     from organizational_units_groups
                              inner join organizational_units on organizational_units.id =
                                                                 organizational_units_groups.organizational_unit_id) o
                    on m.group_id = o.group_id
         inner join (select id, parent_id, display_name, created_at, updated_at
                     from organizational_units) o2 on o.organizational_unit_id = o2.id
`

func (q *Queries) GetUserOrganizationalUnits(ctx context.Context, userID uuid.UUID) ([]OrganizationalUnit, error) {
	rows, err := q.db.Query(ctx, getUserOrganizationalUnits, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []OrganizationalUnit
	for rows.Next() {
		var i OrganizationalUnit
		if err := rows.Scan(
			&i.ID,
			&i.ParentID,
			&i.DisplayName,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUsers = `-- name: GetUsers :many

select id, username, external_id, name, display_name, locale, active, emails, created_at, updated_at
from users
order by created_at
LIMIT $1 OFFSET $2
`

type GetUsersParams struct {
	Limit  int32
	Offset int32
}

// ------------------------------------------------------------------------------------------------------------------
// Users
// ------------------------------------------------------------------------------------------------------------------
func (q *Queries) GetUsers(ctx context.Context, arg GetUsersParams) ([]User, error) {
	rows, err := q.db.Query(ctx, getUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.ExternalID,
			&i.Name,
			&i.DisplayName,
			&i.Locale,
			&i.Active,
			&i.Emails,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUsersByID = `-- name: GetUsersByID :many
select id, username, external_id, name, display_name, locale, active, emails, created_at, updated_at
from users
where id = ANY ($1::uuid[])
order by display_name
`

func (q *Queries) GetUsersByID(ctx context.Context, userIds []uuid.UUID) ([]User, error) {
	rows, err := q.db.Query(ctx, getUsersByID, userIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.ExternalID,
			&i.Name,
			&i.DisplayName,
			&i.Locale,
			&i.Active,
			&i.Emails,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertAPIKey = `-- name: InsertAPIKey :one

insert into api_keys (encodedhash, system, owner, description, created_at, updated_at)
values ($1, $2, $3, $4, now(), now())
returning id, encodedhash, owner, description, system, created_at, updated_at
`

type InsertAPIKeyParams struct {
	Encodedhash string
	System      bool
	Owner       string
	Description sql.NullString
}

// ------------------------------------------------------------------------------------------------------------------
// SCIM API Key
// ------------------------------------------------------------------------------------------------------------------
func (q *Queries) InsertAPIKey(ctx context.Context, arg InsertAPIKeyParams) (ApiKey, error) {
	row := q.db.QueryRow(ctx, insertAPIKey,
		arg.Encodedhash,
		arg.System,
		arg.Owner,
		arg.Description,
	)
	var i ApiKey
	err := row.Scan(
		&i.ID,
		&i.Encodedhash,
		&i.Owner,
		&i.Description,
		&i.System,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const insertScimAPIKey = `-- name: InsertScimAPIKey :one
insert into scim_api_keys (api_key_id, domain, created_at, updated_at)
values ($1, 'default', now(), now())
returning id, domain, api_key_id, created_at, updated_at
`

func (q *Queries) InsertScimAPIKey(ctx context.Context, apiKeyID uuid.UUID) (ScimApiKey, error) {
	row := q.db.QueryRow(ctx, insertScimAPIKey, apiKeyID)
	var i ScimApiKey
	err := row.Scan(
		&i.ID,
		&i.Domain,
		&i.ApiKeyID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const organizationalUnitsCloudAccounts = `-- name: OrganizationalUnitsCloudAccounts :many
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
SELECT cloud_accounts.id, cloud_accounts.cloud, cloud_accounts.tenant_id, cloud_accounts.account_id, cloud_accounts.name, cloud_accounts.active, cloud_accounts.tags_current, cloud_accounts.tags_desired, cloud_accounts.tags_drift_detected, cloud_accounts.created_at, cloud_accounts.updated_at
FROM organizational_unit_ids
         join organizational_units_cloud_accounts
              on organizational_unit_ids.id = organizational_units_cloud_accounts.organizational_unit_id
         join cloud_accounts on organizational_units_cloud_accounts.cloud_account_id = cloud_accounts.id
`

func (q *Queries) OrganizationalUnitsCloudAccounts(ctx context.Context, id []uuid.UUID) ([]CloudAccount, error) {
	rows, err := q.db.Query(ctx, organizationalUnitsCloudAccounts, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CloudAccount
	for rows.Next() {
		var i CloudAccount
		if err := rows.Scan(
			&i.ID,
			&i.Cloud,
			&i.TenantID,
			&i.AccountID,
			&i.Name,
			&i.Active,
			&i.TagsCurrent,
			&i.TagsDesired,
			&i.TagsDriftDetected,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const patchGroupDisplayName = `-- name: PatchGroupDisplayName :exec
update groups
set display_name = $2,
    updated_at   = now()
where id = $1
`

type PatchGroupDisplayNameParams struct {
	ID          uuid.UUID
	DisplayName string
}

func (q *Queries) PatchGroupDisplayName(ctx context.Context, arg PatchGroupDisplayNameParams) error {
	_, err := q.db.Exec(ctx, patchGroupDisplayName, arg.ID, arg.DisplayName)
	return err
}

const patchUser = `-- name: PatchUser :exec
update users
set active     = $2,
    updated_at = now()
where id = $1
`

type PatchUserParams struct {
	ID     uuid.UUID
	Active bool
}

func (q *Queries) PatchUser(ctx context.Context, arg PatchUserParams) error {
	_, err := q.db.Exec(ctx, patchUser, arg.ID, arg.Active)
	return err
}

const searchTag = `-- name: SearchTag :many

select id, cloud, tenant_id, account_id, name, active, tags_current, tags_desired, tags_drift_detected, created_at, updated_at
from cloud_accounts
where cloud = $1
  and tenant_id = $2
  and tags_current ->> $3 = sqcl.arg(tag_value)
`

type SearchTagParams struct {
	Cloud    string
	TenantID string
	TagKey   pgtype.JSONB
}

// ------------------------------------------------------------------------------------------------------------------
// Cloud Accounts
// ------------------------------------------------------------------------------------------------------------------
func (q *Queries) SearchTag(ctx context.Context, arg SearchTagParams) ([]CloudAccount, error) {
	rows, err := q.db.Query(ctx, searchTag, arg.Cloud, arg.TenantID, arg.TagKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CloudAccount
	for rows.Next() {
		var i CloudAccount
		if err := rows.Scan(
			&i.ID,
			&i.Cloud,
			&i.TenantID,
			&i.AccountID,
			&i.Name,
			&i.Active,
			&i.TagsCurrent,
			&i.TagsDesired,
			&i.TagsDriftDetected,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const unAssignAccountFromOUs = `-- name: UnAssignAccountFromOUs :exec
delete
from organizational_units_cloud_accounts
where cloud_account_id = $1
`

func (q *Queries) UnAssignAccountFromOUs(ctx context.Context, cloudAccountID uuid.UUID) error {
	_, err := q.db.Exec(ctx, unAssignAccountFromOUs, cloudAccountID)
	return err
}

const updateCloudAccount = `-- name: UpdateCloudAccount :exec
update cloud_accounts
set tags_desired = $2,
    updated_at   = now()
where id = $1
`

type UpdateCloudAccountParams struct {
	ID          uuid.UUID
	TagsDesired pgtype.JSONB
}

func (q *Queries) UpdateCloudAccount(ctx context.Context, arg UpdateCloudAccountParams) error {
	_, err := q.db.Exec(ctx, updateCloudAccount, arg.ID, arg.TagsDesired)
	return err
}

const updateCloudAccountTagsDriftDetected = `-- name: UpdateCloudAccountTagsDriftDetected :exec
update cloud_accounts
set tags_drift_detected = $1,
    updated_at          = now()
where cloud = $2
  and tenant_id = $3
  and account_id = $4
`

type UpdateCloudAccountTagsDriftDetectedParams struct {
	TagsDriftDetected bool
	Cloud             string
	TenantID          string
	AccountID         string
}

func (q *Queries) UpdateCloudAccountTagsDriftDetected(ctx context.Context, arg UpdateCloudAccountTagsDriftDetectedParams) error {
	_, err := q.db.Exec(ctx, updateCloudAccountTagsDriftDetected,
		arg.TagsDriftDetected,
		arg.Cloud,
		arg.TenantID,
		arg.AccountID,
	)
	return err
}

const updateTag = `-- name: UpdateTag :one
update standard_tags
set display_name = $2,
    key          = $3,
    description  = $4,
    updated_at   = now()
where id = $1
returning id, display_name, description, key, created_at, updated_at
`

type UpdateTagParams struct {
	ID          uuid.UUID
	DisplayName string
	Key         string
	Description string
}

func (q *Queries) UpdateTag(ctx context.Context, arg UpdateTagParams) (StandardTag, error) {
	row := q.db.QueryRow(ctx, updateTag,
		arg.ID,
		arg.DisplayName,
		arg.Key,
		arg.Description,
	)
	var i StandardTag
	err := row.Scan(
		&i.ID,
		&i.DisplayName,
		&i.Description,
		&i.Key,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
update users
set username     =$2,
    name         = $3,
    display_name = $4,
    emails       = $5,
    active       = $6,
    external_id  = $7,
    locale       = $8,
    updated_at   = now()
where id = $1
returning id, username, external_id, name, display_name, locale, active, emails, created_at, updated_at
`

type UpdateUserParams struct {
	ID          uuid.UUID
	Username    string
	Name        pgtype.JSONB
	DisplayName sql.NullString
	Emails      pgtype.JSONB
	Active      bool
	ExternalID  sql.NullString
	Locale      sql.NullString
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.ID,
		arg.Username,
		arg.Name,
		arg.DisplayName,
		arg.Emails,
		arg.Active,
		arg.ExternalID,
		arg.Locale,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.ExternalID,
		&i.Name,
		&i.DisplayName,
		&i.Locale,
		&i.Active,
		&i.Emails,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
