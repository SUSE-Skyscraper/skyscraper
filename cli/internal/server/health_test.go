package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/suse-skyscraper/skyscraper/test/helpers"

	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	req, _ := http.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	Health(w, req)

	_ = helpers.AssertOpenAPI(t, w, req)

	result := w.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)
	_ = result.Body.Close()
}
