package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/suse-skyscraper/skyscraper/internal/auth"
	"github.com/suse-skyscraper/skyscraper/internal/db"
	"github.com/suse-skyscraper/skyscraper/internal/helpers"
	"github.com/suse-skyscraper/skyscraper/internal/server/middleware"
	"github.com/suse-skyscraper/skyscraper/internal/server/responses"
)

func TestV1Profile(t *testing.T) {
	expectedUsername := "foo@bar.com"
	caller := auth.Caller{
		ID:   uuid.New(),
		Type: auth.CallerUser,
	}
	req, _ := http.NewRequest("GET", "/api/v1/profile", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx := req.Context()

	ctx = context.WithValue(ctx, middleware.ContextCaller, caller)
	req = req.WithContext(ctx)
	testApp := helpers.NewTestApp()

	user := db.User{
		Username: expectedUsername,
	}

	testApp.Repository.On("FindUser", mock.Anything, caller.ID.String()).Return(user, nil)

	w := httptest.NewRecorder()
	V1Profile(testApp.App)(w, req)
	body := helpers.AssertOpenAPI(t, w, req)

	var resp responses.UserResponse
	err := json.Unmarshal(body, &resp)
	assert.Nil(t, err)
	assert.Equal(t, expectedUsername, resp.Data.Attributes.Username)
}

func TestV1ProfileNoContext(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/profile", nil)
	req.Header.Set("Content-Type", "application/json")
	testApp := helpers.NewTestApp()

	w := httptest.NewRecorder()
	V1Profile(testApp.App)(w, req)
	_ = helpers.AssertOpenAPI(t, w, req)

	result := w.Result()
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	_ = result.Body.Close()
}
