package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/suse-skyscraper/skyscraper/cli/internal/auth"
	"github.com/suse-skyscraper/skyscraper/cli/internal/db"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/middleware"
	testhelpers2 "github.com/suse-skyscraper/skyscraper/cli/internal/testhelpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestV1Profile(t *testing.T) {
	callerUser := auth.Caller{
		ID:   uuid.New(),
		Type: auth.CallerUser,
	}
	callerAPIKey := auth.Caller{
		ID:   uuid.New(),
		Type: auth.CallerAPIKey,
	}
	user := db.User{
		Username: "foo@bar.com",
	}

	tests := []struct {
		caller    interface{}
		status    int
		findError error
	}{
		{
			caller: callerUser,
			status: http.StatusOK,
		},
		{
			caller: callerAPIKey,
			status: http.StatusNotFound,
		},
		{
			findError: fmt.Errorf(""),
			caller:    callerUser,
			status:    http.StatusInternalServerError,
		},
		{
			caller: interface{}(nil),
			status: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		req, _ := http.NewRequest("GET", "/api/v1/caller/profile", nil)
		req.Header.Set("Content-Type", "application/json")
		ctx := req.Context()

		ctx = context.WithValue(ctx, middleware.ContextCaller, tc.caller)
		req = req.WithContext(ctx)
		testApp := testhelpers2.NewTestApp()

		testApp.Repository.On("FindUser", mock.Anything, mock.Anything).Return(user, tc.findError)

		w := httptest.NewRecorder()
		V1CallerProfile(testApp.App)(w, req)
		_ = testhelpers2.AssertOpenAPI(t, w, req)

		result := w.Result()
		assert.Equal(t, tc.status, result.StatusCode)
		_ = result.Body.Close()
	}
}

func TestV1CallerCloudAccounts(t *testing.T) {
	cloudAccount := testhelpers2.FactoryCloudAccount()
	user := auth.Caller{
		Type: auth.CallerUser,
		ID:   uuid.New(),
	}
	apiKey := auth.Caller{
		Type: auth.CallerAPIKey,
		ID:   uuid.New(),
	}
	ou := testhelpers2.FactoryOrganizationalUnit()

	tests := []struct {
		context          interface{}
		getOUError       error
		getAccountsError error
		statusCode       int
	}{
		{
			context:    user,
			statusCode: http.StatusOK,
		},
		{
			context:    apiKey,
			statusCode: http.StatusOK,
		},
		{
			context:    apiKey,
			getOUError: errors.New("test error"),
			statusCode: http.StatusInternalServerError,
		},
		{
			context:    interface{}(nil),
			statusCode: http.StatusInternalServerError,
		},
		{
			context: auth.Caller{
				Type: auth.CallerType(0),
				ID:   uuid.New(),
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			context:          apiKey,
			getAccountsError: errors.New("test error"),
			statusCode:       http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		req, _ := http.NewRequest("GET", "/api/v1/caller/cloud_accounts", nil)
		w := httptest.NewRecorder()

		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.ContextCaller, tc.context)
		req = req.WithContext(ctx)

		testApp := testhelpers2.NewTestApp()

		testApp.Repository.
			On("GetUserOrganizationalUnits", mock.Anything, mock.Anything).
			Return([]db.OrganizationalUnit{ou}, tc.getOUError)
		testApp.Repository.
			On("GetAPIKeysOrganizationalUnits", mock.Anything, mock.Anything).
			Return([]db.OrganizationalUnit{ou}, tc.getOUError)
		testApp.Repository.
			On("OrganizationalUnitsCloudAccounts", mock.Anything, mock.Anything).
			Return([]db.CloudAccount{cloudAccount}, tc.getAccountsError)

		V1CallerCloudAccounts(testApp.App)(w, req)

		_ = testhelpers2.AssertOpenAPI(t, w, req)

		result := w.Result()
		assert.Equal(t, tc.statusCode, result.StatusCode)
		_ = result.Body.Close()
	}
}
