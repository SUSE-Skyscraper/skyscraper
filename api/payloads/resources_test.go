package payloads

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrUpdateResourcePayload_Bind(t *testing.T) {
	tc := []struct {
		name                string
		accountName         string
		tagsCurrent         map[string]string
		tagsCurrentExpected map[string]string
		tagsDesired         map[string]string
		tagsDesiredExpected map[string]string
	}{
		{
			name:                "valid tags",
			accountName:         "test",
			tagsCurrent:         map[string]string{"foo": "bar"},
			tagsCurrentExpected: map[string]string{"foo": "bar"},
			tagsDesired:         map[string]string{"bar": "foo"},
			tagsDesiredExpected: map[string]string{"bar": "foo"},
		},
		{
			name:                "nil tags",
			accountName:         "test",
			tagsCurrent:         nil,
			tagsCurrentExpected: map[string]string{},
			tagsDesired:         nil,
			tagsDesiredExpected: map[string]string{},
		},
		{
			name:                "empty tags",
			accountName:         "test",
			tagsCurrent:         map[string]string{},
			tagsCurrentExpected: map[string]string{},
			tagsDesired:         map[string]string{},
			tagsDesiredExpected: map[string]string{},
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			payload := CreateOrUpdateResourcePayload{
				Data: CreateOrUpdateResourcePayloadData{
					AccountName: tt.accountName,
					TagsCurrent: tt.tagsCurrent,
					TagsDesired: tt.tagsDesired,
				}}
			err := payload.Bind(nil)
			assert.Nil(t, err)
			assert.Equal(t, tt.tagsCurrentExpected, payload.Data.TagsCurrent)
			assert.Equal(t, tt.tagsDesiredExpected, payload.Data.TagsDesired)
		})
	}
}
