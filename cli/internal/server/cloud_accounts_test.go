package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/suse-skyscraper/skyscraper/cli/internal/db"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/middleware"
	"github.com/suse-skyscraper/skyscraper/cli/internal/testhelpers"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestV1ListCloudAccounts(t *testing.T) {
	cloudAccount := testhelpers.FactoryCloudAccount()

	tests := []struct {
		getError   error
		statusCode int
	}{
		{
			getError:   nil,
			statusCode: http.StatusOK,
		},
		{
			getError:   errors.New("test error"),
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		req, _ := http.NewRequest("GET", "/api/v1/cloud_accounts?cloud=AWS", nil)
		w := httptest.NewRecorder()
		testApp, err := testhelpers.NewTestApp()
		if err != nil {
			t.Fatal(err)
		}

		testApp.Searcher.
			On("SearchCloudAccounts", mock.Anything, mock.Anything).
			Return([]db.CloudAccount{cloudAccount}, tc.getError)

		V1ListCloudAccounts(testApp.App)(w, req)

		_ = testhelpers.AssertOpenAPI(t, w, req)

		result := w.Result()
		assert.Equal(t, tc.statusCode, result.StatusCode)
		_ = result.Body.Close()

		testApp.Close() // Close the test app, we cannot defer in a loop
	}
}

func TestV1GetCloudAccount(t *testing.T) {
	cloudAccount := testhelpers.FactoryCloudAccount()

	tests := []struct {
		context interface{}
		status  int
	}{
		{
			context: cloudAccount,
			status:  http.StatusOK,
		},
		{
			context: interface{}(nil),
			status:  http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		req, _ := http.NewRequest("GET", "/api/v1/cloud_accounts/123456", nil)
		w := httptest.NewRecorder()
		testApp, err := testhelpers.NewTestApp()
		if err != nil {
			t.Fatal(err)
		}

		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.ContextCloudAccount, tc.context)
		req = req.WithContext(ctx)

		V1GetCloudAccount(testApp.App)(w, req)

		_ = testhelpers.AssertOpenAPI(t, w, req)

		result := w.Result()
		assert.Equal(t, tc.status, result.StatusCode)
		_ = result.Body.Close()

		testApp.Close() // Close the test app, we cannot defer in a loop
	}
}
