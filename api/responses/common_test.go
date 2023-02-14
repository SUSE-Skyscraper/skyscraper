package responses

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrResponse_Render(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	resp := ErrResponse{}
	err := resp.Render(nil, req)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}

func TestErrInvalidRequest(t *testing.T) {
	requestError := fmt.Errorf("request error")
	resp := ErrInvalidRequest(requestError)
	assert.NotNil(t, resp)
}
