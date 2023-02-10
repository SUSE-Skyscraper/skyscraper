package payloads

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTagPayload_Bind(t *testing.T) {
	payload := CreateTagPayload{
		Data: CreateTagPayloadData{
			DisplayName: "test",
			Required:    true,
			Description: "test",
			Key:         "test",
		},
	}

	err := payload.Bind(nil)
	assert.Nil(t, err)
}

func TestUpdateTagPayload_Bind(t *testing.T) {
	payload := UpdateTagPayload{
		Data: UpdateTagPayloadData{
			DisplayName: "test",
			Required:    true,
			Description: "test",
			Key:         "test",
		},
	}

	err := payload.Bind(nil)
	assert.Nil(t, err)
}
