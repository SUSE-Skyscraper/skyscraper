package server

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackc/pgtype"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/suse-skyscraper/skyscraper/internal/db"
	"github.com/suse-skyscraper/skyscraper/internal/helpers"
	"github.com/suse-skyscraper/skyscraper/internal/server/middleware"
)

var cloudAccount = db.CloudAccount{
	Cloud:             "aws",
	TenantID:          "1234",
	AccountID:         "12345",
	Name:              "test",
	Active:            true,
	TagsCurrent:       pgtype.JSONB{Bytes: []byte("{}"), Status: pgtype.Present},
	TagsDesired:       pgtype.JSONB{Bytes: []byte("{}"), Status: pgtype.Present},
	TagsDriftDetected: false,
	CreatedAt:         time.Now(),
	UpdatedAt:         time.Now(),
}

func TestV1ListCloudAccounts(t *testing.T) {
	type Test struct {
		getError   error
		statusCode int
	}

	tests := []Test{
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
		req, _ := http.NewRequest("GET", "/api/v1/cloud_tenants/cloud/aws/tenant/12345/accounts", nil)
		w := httptest.NewRecorder()
		testApp := helpers.NewTestApp()

		testApp.Repository.
			On("SearchCloudAccounts", mock.Anything, mock.Anything).
			Return([]db.CloudAccount{cloudAccount}, tc.getError)

		V1ListCloudAccounts(testApp.App)(w, req)

		_ = helpers.AssertOpenAPI(t, w, req)
		assert.Equal(t, tc.statusCode, w.Result().StatusCode)
	}
}

func TestV1UpdateCloudTenantAccount(t *testing.T) {
	type PubAckFuture struct {
		nats.PubAckFuture
	}
	type Test struct {
		tags         []byte
		updateError  error
		publishError error
		statusCode   int
	}

	tests := []Test{
		{
			tags:         []byte(`{"data": {"tags_desired": {}}}`),
			updateError:  nil,
			publishError: nil,
			statusCode:   http.StatusOK,
		},
		{
			tags:         []byte(`{}`),
			updateError:  nil,
			publishError: nil,
			statusCode:   http.StatusBadRequest,
		},
		{
			tags:         []byte(`{"data": {"tags_desired": {}}}`),
			updateError:  errors.New(""),
			publishError: nil,
			statusCode:   http.StatusInternalServerError,
		},
		{
			tags:         []byte(`{"data": {"tags_desired": {}}}`),
			updateError:  nil,
			publishError: errors.New(""),
			statusCode:   http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		req, _ := http.NewRequest("PUT",
			"/api/v1/cloud_tenants/cloud/aws/tenant/1234/accounts/12345",
			bytes.NewReader(tc.tags))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		testApp := helpers.NewTestApp()

		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.CloudAccount, cloudAccount)
		req = req.WithContext(ctx)

		testApp.Repository.On("UpdateCloudAccount", mock.Anything, mock.Anything).Return(cloudAccount, tc.updateError)
		testApp.JS.On("PublishAsync", mock.Anything, mock.Anything, mock.Anything).Return(PubAckFuture{}, tc.publishError)

		V1UpdateCloudTenantAccount(testApp.App)(w, req)

		_ = helpers.AssertOpenAPI(t, w, req)
		assert.Equal(t, tc.statusCode, w.Result().StatusCode)
	}
}

func TestV1GetCloudAccount(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/cloud_tenants/cloud/aws/tenant/12345/accounts/123456", nil)
	w := httptest.NewRecorder()
	testApp := helpers.NewTestApp()

	ctx := req.Context()
	ctx = context.WithValue(ctx, middleware.CloudAccount, cloudAccount)
	req = req.WithContext(ctx)

	V1GetCloudAccount(testApp.App)(w, req)

	_ = helpers.AssertOpenAPI(t, w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}
