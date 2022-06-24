package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suse-skyscraper/skyscraper/internal/helpers"
)

func TestHealth(t *testing.T) {
	req, _ := http.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	Health(w, req)

	_ = helpers.AssertOpenAPI(t, w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}
