package payloads

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAPIKeyPayload_Bind(t *testing.T) {
	payload := CreateAPIKeyPayload{
		Data: CreateAPIKeyPayloadData{
			Owner:       "test",
			Description: "test",
		}}
	err := payload.Bind(nil)
	assert.Nil(t, err)
}
