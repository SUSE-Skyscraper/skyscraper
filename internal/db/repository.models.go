package db

import (
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/filters"
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
