package payloads

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAssignCloudAccountToOUPayload_Bind(t *testing.T) {
	tc := []struct {
		name        string
		ouID        string
		expectError bool
	}{
		{
			name:        "valid ou id",
			ouID:        uuid.New().String(),
			expectError: false,
		},
		{
			name:        "invalid ou id",
			ouID:        "foobar",
			expectError: true,
		},
		{
			name:        "empty ou id",
			ouID:        "",
			expectError: true,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			payload := AssignCloudAccountToOUPayload{
				Data: AssignCloudAccountToOUPayloadData{
					OrganizationalUnitID: tt.ouID,
				}}
			err := payload.Bind(nil)
			if tt.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
