package responses

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloudAccountResponse_Render(t *testing.T) {
	resp := CloudAccountResponse{}
	err := resp.Render(nil, nil)
	assert.Nil(t, err)
}

func TestCloudAccountListResponse_Render(t *testing.T) {
	resp := CloudAccountListResponse{}
	err := resp.Render(nil, nil)
	assert.Nil(t, err)
}
