package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/suse-skyscraper/skyscraper/internal/scim/payloads"
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

	UpdateCloudAccountTagsDriftDetected(ctx context.Context, input UpdateCloudAccountTagsDriftDetectedParams) error
	CreateOrInsertCloudAccount(ctx context.Context, input CreateOrInsertCloudAccountParams) (CloudAccount, error)
	FindCloudAccount(ctx context.Context, input FindCloudAccountInput) (CloudAccount, error)
	UpdateCloudAccount(ctx context.Context, input UpdateCloudAccountParams) (CloudAccount, error)
	SearchCloudAccounts(ctx context.Context, input SearchCloudAccountsInput) ([]CloudAccount, error)

	CreateTag(ctx context.Context, input CreateTagParams) (Tag, error)
	UpdateTag(ctx context.Context, input UpdateTagParams) (Tag, error)
	FindTag(ctx context.Context, id uuid.UUID) (Tag, error)
	GetTags(ctx context.Context) ([]Tag, error)

	GetCloudTenants(ctx context.Context) ([]CloudTenant, error)
	CreateCloudTenant(ctx context.Context, input CreateCloudTenantParams) error

	FindGroup(ctx context.Context, id string) (Group, error)
	CreateGroup(ctx context.Context, displayName string) (Group, error)
	DeleteGroup(ctx context.Context, id string) error
	UpdateGroup(ctx context.Context, input PatchGroupDisplayNameParams) (Group, error)
	RemoveUserFromGroup(ctx context.Context, userID, groupID uuid.UUID) error
	AddUserToGroup(ctx context.Context, userID, groupID uuid.UUID) error
	ReplaceUsersInGroup(ctx context.Context, groupID uuid.UUID, members []payloads.MemberPatch) error
	AddUsersToGroup(ctx context.Context, groupID uuid.UUID, members []payloads.MemberPatch) error
	GetGroupMembership(ctx context.Context, idString string) ([]GetGroupMembershipRow, error)
	GetGroups(ctx context.Context, params GetGroupsParams) (int64, []Group, error)

	FindUser(ctx context.Context, id string) (User, error)
	FindUserByUsername(ctx context.Context, username string) (User, error)
	GetScimUsers(ctx context.Context, input GetScimUsersInput) (int64, []User, error)
	CreateUser(ctx context.Context, input CreateUserParams) (User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	UpdateUser(ctx context.Context, id uuid.UUID, input UpdateUserParams) (User, error)
	ScimPatchUser(ctx context.Context, input PatchUserParams) error

	GetPolicies(ctx context.Context) ([]Policy, error)
	TruncatePolicies(ctx context.Context) error
	CreatePolicy(ctx context.Context, input AddPolicyParams) error
	RemovePolicy(ctx context.Context, input RemovePolicyParams) error

	InsertAPIKey(ctx context.Context, token string) (ScimApiKey, error)
	FindAPIKey(ctx context.Context, token string) (ScimApiKey, error)

	GetAuditLogs(ctx context.Context) ([]AuditLog, []User, error)
	GetAuditLogsForTarget(ctx context.Context, input GetAuditLogsForTargetParams) ([]AuditLog, []User, error)
	CreateAuditLog(ctx context.Context, input CreateAuditLogParams) (AuditLog, error)
}
