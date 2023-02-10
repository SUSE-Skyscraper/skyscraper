package payloads

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrUpdateTenantPayload_Bind(t *testing.T) {
	tt := []struct {
		name        string
		tenantName  string
		tenantCloud string
		expectError bool
	}{
		{
			name:        "valid tenant name",
			tenantName:  "test",
			tenantCloud: "aws",
			expectError: false,
		},
		{
			name:        "invalid tenant name",
			tenantName:  "",
			tenantCloud: "aws",
			expectError: true,
		},
		{
			name:        "invalid tenant cloud",
			tenantName:  "test",
			tenantCloud: "",
			expectError: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			payload := CreateOrUpdateTenantPayload{
				Data: CreateOrUpdateTenantPayloadData{
					Name:  tc.tenantName,
					Cloud: tc.tenantCloud,
				}}
			err := payload.Bind(nil)
			if tc.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
