package db

import (
	"context"
	"fmt"
	"strings"
)

const baseSearchCloudAccounts = `
select cloud, tenant_id, account_id, name, active, tags_current, tags_desired,tags_drift_detected, created_at,
updated_at
from cloud_accounts
where
`

type SearchCloudAccountsInput struct {
	Filters map[string]interface{}
}

func (s *Searches) SearchCloudAccounts(ctx context.Context, input SearchCloudAccountsInput) ([]CloudAccount, error) {
	b := strings.Builder{}
	var accounts []CloudAccount

	_, err := b.Write([]byte(baseSearchCloudAccounts))
	if err != nil {
		return accounts, err
	}

	var values []interface{}
	if input.Filters == nil || len(input.Filters) == 0 {
		_, err := b.Write([]byte(`1=1`))
		if err != nil {
			return accounts, err
		}
	} else {
		i := 0
		for key, value := range input.Filters {
			var query string
			var prefix string
			if i != 0 {
				prefix = "and"
			}
			switch key {
			case "tenant_id", "cloud", "account_id", "name", "active", "tags_drift_detected":
				// we know the columns here, so we don't need to worry about sql injection
				query = fmt.Sprintf(" %s %s = $%d ", prefix, key, i+1)
				values = append(values, value)
				i++
			default:
				query = fmt.Sprintf(" %s tags_current ->> $%d like $%d", prefix, i+1, i+2)
				values = append(values, key, value)
				i += 2
			}
			_, err := b.Write([]byte(query))
			if err != nil {
				return accounts, err
			}
		}
	}

	query := b.String()

	rows, err := s.pool.Query(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var account CloudAccount
		err := rows.Scan(
			&account.Cloud,
			&account.TenantID,
			&account.AccountID,
			&account.Name,
			&account.Active,
			&account.TagsCurrent,
			&account.TagsDesired,
			&account.TagsDriftDetected,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}
