package responses

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuditLogsResponse_Render(t *testing.T) {
	resp := AuditLogsResponse{}
	err := resp.Render(nil, nil)
	assert.Nil(t, err)
}
