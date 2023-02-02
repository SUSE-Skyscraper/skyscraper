package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewRepository(pool *pgxpool.Pool, db Querier) *Repository {
	search := newSearch(pool)
	return &Repository{
		db:           db,
		postgresPool: pool,
		search:       search,
	}
}

type RepositoryQueries interface {
	Begin(ctx context.Context) (RepositoryQueries, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	CreateOrUpdateCloudAccount(ctx context.Context, input CreateOrUpdateCloudAccountParams) (CloudAccount, error)
	FindCloudAccount(ctx context.Context, input FindCloudAccountInput) (CloudAccount, error)
	UpdateCloudAccount(ctx context.Context, input UpdateCloudAccountParams) (CloudAccount, error)
	SearchCloudAccounts(ctx context.Context, input SearchCloudAccountsInput) ([]CloudAccount, error)
	AssignCloudAccountToOrganizationalUnit(ctx context.Context, id, organizationalUnitID uuid.UUID) error
	UnAssignCloudAccountFromOrganizationalUnits(ctx context.Context, id uuid.UUID) error
	OrganizationalUnitsCloudAccounts(ctx context.Context, id []uuid.UUID) ([]CloudAccount, error)

	CreateTag(ctx context.Context, input CreateTagParams) (StandardTag, error)
	UpdateTag(ctx context.Context, input UpdateTagParams) (StandardTag, error)
	FindTag(ctx context.Context, id uuid.UUID) (StandardTag, error)
	GetTags(ctx context.Context) ([]StandardTag, error)

	GetCloudTenant(ctx context.Context, params GetCloudTenantParams) (CloudTenant, error)
	GetCloudTenants(ctx context.Context) ([]CloudTenant, error)
	CreateOrUpdateCloudTenant(ctx context.Context, input CreateOrUpdateCloudTenantParams) (CloudTenant, error)

	FindGroup(ctx context.Context, id string) (Group, error)
	CreateGroup(ctx context.Context, displayName string) (Group, error)
	DeleteGroup(ctx context.Context, id string) error
	UpdateGroup(ctx context.Context, input PatchGroupDisplayNameParams) (Group, error)
	RemoveUserFromGroup(ctx context.Context, userID, groupID uuid.UUID) error
	AddUserToGroup(ctx context.Context, userID, groupID uuid.UUID) error
	ReplaceUsersInGroup(ctx context.Context, groupID uuid.UUID, members []uuid.UUID) error
	AddUsersToGroup(ctx context.Context, groupID uuid.UUID, members []uuid.UUID) error
	GetGroupMembership(ctx context.Context, idString string) ([]GetGroupMembershipRow, error)
	GetGroups(ctx context.Context, params GetGroupsParams) (int64, []Group, error)

	GetUserOrganizationalUnits(ctx context.Context, id uuid.UUID) ([]OrganizationalUnit, error)
	FindUser(ctx context.Context, id string) (User, error)
	FindUserByUsername(ctx context.Context, username string) (User, error)
	GetUsers(ctx context.Context, input GetUsersParams) ([]User, error)
	GetScimUsers(ctx context.Context, input GetScimUsersInput) (int64, []User, error)
	CreateUser(ctx context.Context, input CreateUserParams) (User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	UpdateUser(ctx context.Context, id uuid.UUID, input UpdateUserParams) (User, error)
	ScimPatchUser(ctx context.Context, input PatchUserParams) error

	GetAPIKeysOrganizationalUnits(ctx context.Context, id uuid.UUID) ([]OrganizationalUnit, error)
	InsertScimAPIKey(ctx context.Context, encodedHash string) (ApiKey, error)
	DeleteScimAPIKey(ctx context.Context) error
	FindAPIKey(ctx context.Context, id uuid.UUID) (ApiKey, error)
	FindScimAPIKey(ctx context.Context) (ApiKey, error)
	GetAPIKeys(ctx context.Context) ([]ApiKey, error)
	CreateAPIKey(ctx context.Context, input InsertAPIKeyParams) (ApiKey, error)

	GetAuditLogs(ctx context.Context) ([]AuditLog, []any, error)
	GetAuditLogsForTarget(ctx context.Context, input GetAuditLogsForTargetParams) ([]AuditLog, []any, error)
	CreateAuditLog(ctx context.Context, input CreateAuditLogParams) (AuditLog, error)

	CreateOrganizationalUnit(ctx context.Context, input CreateOrganizationalUnitParams) (OrganizationalUnit, error)
	GetOrganizationalUnits(ctx context.Context) ([]OrganizationalUnit, error)
	FindOrganizationalUnit(ctx context.Context, id uuid.UUID) (OrganizationalUnit, error)
	GetOrganizationalUnitChildren(ctx context.Context, id uuid.UUID) ([]OrganizationalUnit, error)
	GetOrganizationalUnitCloudAccounts(ctx context.Context, id uuid.UUID) ([]CloudAccount, error)
	DeleteOrganizationalUnit(ctx context.Context, id uuid.UUID) error
}
