// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	AddPolicy(ctx context.Context, arg AddPolicyParams) error
	CreateAuditLog(ctx context.Context, arg CreateAuditLogParams) (AuditLog, error)
	//------------------------------------------------------------------------------------------------------------------
	// Cloud Tenants
	//------------------------------------------------------------------------------------------------------------------
	CreateCloudTenant(ctx context.Context, arg CreateCloudTenantParams) error
	CreateGroup(ctx context.Context, displayName string) (Group, error)
	CreateMembershipForUserAndGroup(ctx context.Context, arg CreateMembershipForUserAndGroupParams) error
	CreateOrInsertCloudAccount(ctx context.Context, arg CreateOrInsertCloudAccountParams) (CloudAccount, error)
	CreateTag(ctx context.Context, arg CreateTagParams) (Tag, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteAPIKey(ctx context.Context, id uuid.UUID) error
	DeleteGroup(ctx context.Context, id uuid.UUID) error
	DeleteScimAPIKey(ctx context.Context) error
	DeleteTag(ctx context.Context, id uuid.UUID) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	DropMembershipForGroup(ctx context.Context, groupID uuid.UUID) error
	DropMembershipForUserAndGroup(ctx context.Context, arg DropMembershipForUserAndGroupParams) error
	FindByUsername(ctx context.Context, username string) (User, error)
	FindCloudAccount(ctx context.Context, id uuid.UUID) (CloudAccount, error)
	FindCloudAccountByCloudAndTenant(ctx context.Context, arg FindCloudAccountByCloudAndTenantParams) (CloudAccount, error)
	FindScimAPIKey(ctx context.Context) (ApiKey, error)
	FindTag(ctx context.Context, id uuid.UUID) (Tag, error)
	//------------------------------------------------------------------------------------------------------------------
	// Audit Logs
	//------------------------------------------------------------------------------------------------------------------
	GetAuditLogs(ctx context.Context) ([]AuditLog, error)
	GetAuditLogsForTarget(ctx context.Context, arg GetAuditLogsForTargetParams) ([]AuditLog, error)
	GetCloudTenant(ctx context.Context, arg GetCloudTenantParams) (CloudTenant, error)
	GetCloudTenants(ctx context.Context) ([]CloudTenant, error)
	GetGroup(ctx context.Context, id uuid.UUID) (Group, error)
	GetGroupCount(ctx context.Context) (int64, error)
	//------------------------------------------------------------------------------------------------------------------
	// Membership
	//------------------------------------------------------------------------------------------------------------------
	GetGroupMembership(ctx context.Context, groupID uuid.UUID) ([]GetGroupMembershipRow, error)
	GetGroupMembershipForUser(ctx context.Context, arg GetGroupMembershipForUserParams) (GetGroupMembershipForUserRow, error)
	//------------------------------------------------------------------------------------------------------------------
	// Groups
	//------------------------------------------------------------------------------------------------------------------
	GetGroups(ctx context.Context, arg GetGroupsParams) ([]Group, error)
	//------------------------------------------------------------------------------------------------------------------
	// Policies
	//
	// 6ba7b812-9dad-11d1-80b4-00c04fd430c8 is NameSpace_OID as specified in rfc4122 (https://tools.ietf.org/html/rfc4122)
	// we use uuid v5 so we can calculate the id from a collection of values
	//------------------------------------------------------------------------------------------------------------------
	GetPolicies(ctx context.Context) ([]Policy, error)
	//------------------------------------------------------------------------------------------------------------------
	// Tags
	//------------------------------------------------------------------------------------------------------------------
	GetTags(ctx context.Context) ([]Tag, error)
	GetUser(ctx context.Context, id uuid.UUID) (User, error)
	GetUserCount(ctx context.Context) (int64, error)
	//------------------------------------------------------------------------------------------------------------------
	// Users
	//------------------------------------------------------------------------------------------------------------------
	GetUsers(ctx context.Context, arg GetUsersParams) ([]User, error)
	GetUsersById(ctx context.Context, dollar_1 []uuid.UUID) ([]User, error)
	//------------------------------------------------------------------------------------------------------------------
	// SCIM API Key
	//------------------------------------------------------------------------------------------------------------------
	InsertAPIKey(ctx context.Context, encodedhash string) (ApiKey, error)
	InsertScimAPIKey(ctx context.Context, apiKeyID uuid.UUID) (ScimApiKey, error)
	PatchGroupDisplayName(ctx context.Context, arg PatchGroupDisplayNameParams) error
	PatchUser(ctx context.Context, arg PatchUserParams) error
	RemovePoliciesForGroup(ctx context.Context, v1 string) error
	RemovePolicy(ctx context.Context, arg RemovePolicyParams) error
	//------------------------------------------------------------------------------------------------------------------
	// Cloud Account Metadata
	//------------------------------------------------------------------------------------------------------------------
	SearchTag(ctx context.Context, arg SearchTagParams) ([]CloudAccount, error)
	TruncatePolicies(ctx context.Context) error
	UpdateCloudAccount(ctx context.Context, arg UpdateCloudAccountParams) error
	UpdateCloudAccountTagsDriftDetected(ctx context.Context, arg UpdateCloudAccountTagsDriftDetectedParams) error
	UpdateTag(ctx context.Context, arg UpdateTagParams) error
	UpdateUser(ctx context.Context, arg UpdateUserParams) error
}

var _ Querier = (*Queries)(nil)
