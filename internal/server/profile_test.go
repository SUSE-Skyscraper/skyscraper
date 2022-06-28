package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suse-skyscraper/skyscraper/internal/helpers"
	"github.com/suse-skyscraper/skyscraper/internal/server/middleware"
)

func TestV1Profile(t *testing.T) {
	expectedEmail := "foo@bar.com"
	req, _ := http.NewRequest("GET", "/api/v1/profile", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx := req.Context()
	ctx = context.WithValue(ctx, middleware.UserEmail, "foo@bar.com")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	V1Profile(w, req)
	body := helpers.AssertOpenAPI(t, w, req)

	var userProfile userProfile
	err := json.Unmarshal(body, &userProfile)
	assert.Nil(t, err)
	assert.Equal(t, expectedEmail, userProfile.Email)
}

func TestV1ProfileNoContext(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/profile", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	V1Profile(w, req)
	_ = helpers.AssertOpenAPI(t, w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}
