package responses

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIKeyResponse_Render(t *testing.T) {
	resp := APIKeyResponse{}
	err := resp.Render(nil, nil)
	assert.Nil(t, err)
}

func TestAPIKeysResponse_Render(t *testing.T) {
	resp := APIKeysResponse{}
	err := resp.Render(nil, nil)
	assert.Nil(t, err)
}
