package payloads

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrganizationalUnitsPayload_Bind(t *testing.T) {
	tc := []struct {
		name        string
		parentID    string
		displayName string
		expectError bool
	}{
		{
			name:        "valid parent id",
			parentID:    "b3d0c1a0-5b1e-4b5e-8f1d-8c1c1c1c1c1c",
			displayName: "test",
			expectError: false,
		},
		{
			name:        "invalid parent id",
			parentID:    "foobar",
			displayName: "test",
			expectError: true,
		},
		{
			name:        "empty parent id",
			parentID:    "",
			displayName: "test",
			expectError: false,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			payload := CreateOrganizationalUnitsPayload{
				Data: CreateOrganizationalUnitsPayloadData{
					ParentID:    tt.parentID,
					DisplayName: tt.displayName,
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
