package db

import (
	"github.com/suse-skyscraper/skyscraper/internal/scim/filters"
)

type FindCloudAccountInput struct {
	Cloud     string
	TenantID  string
	AccountID string
	ID        string
}

type GetScimUsersInput struct {
	Filters []filters.Filter
	Offset  int32
	Limit   int32
}
