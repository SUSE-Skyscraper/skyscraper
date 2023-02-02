package server

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/suse-skyscraper/skyscraper/cli/internal/auth"
	"github.com/suse-skyscraper/skyscraper/cli/internal/db"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/middleware"
	testhelpers2 "github.com/suse-skyscraper/skyscraper/cli/internal/testhelpers"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestV1ListCloudAccounts(t *testing.T) {
	cloudAccount := testhelpers2.FactoryCloudAccount()

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
		testApp := testhelpers2.NewTestApp()

		testApp.Repository.
			On("SearchCloudAccounts", mock.Anything, mock.Anything).
			Return([]db.CloudAccount{cloudAccount}, tc.getError)

		V1ListCloudAccounts(testApp.App)(w, req)

		_ = testhelpers2.AssertOpenAPI(t, w, req)

		result := w.Result()
		assert.Equal(t, tc.statusCode, result.StatusCode)
		_ = result.Body.Close()
	}
}

func TestV1UpdateCloudAccount(t *testing.T) {
	cloudAccount := testhelpers2.FactoryCloudAccount()

	type PubAckFuture struct {
		nats.PubAckFuture
	}

	tests := []struct {
		tags                []byte
		updateError         error
		beginError          error
		commitError         error
		createAuditLogError error
		publishError        error
		statusCode          int
		context             interface{}
	}{
		{
			context:    cloudAccount,
			tags:       []byte(`{"data": {"tags_desired": {}}}`),
			statusCode: http.StatusOK,
		},
		{
			context:    cloudAccount,
			tags:       []byte(`{}`),
			statusCode: http.StatusBadRequest,
		},
		{
			context:     cloudAccount,
			tags:        []byte(`{"data": {"tags_desired": {}}}`),
			updateError: errors.New(""),
			statusCode:  http.StatusInternalServerError,
		},
		{
			context:      cloudAccount,
			tags:         []byte(`{"data": {"tags_desired": {}}}`),
			publishError: errors.New(""),
			statusCode:   http.StatusOK,
		},
		{
			context:    cloudAccount,
			tags:       []byte(`{"data": {"tags_desired": {}}}`),
			beginError: errors.New(""),
			statusCode: http.StatusInternalServerError,
		},
		{
			context:    interface{}(nil),
			tags:       []byte(`{"data": {"tags_desired": {}}}`),
			statusCode: http.StatusInternalServerError,
		},
		{
			context:     cloudAccount,
			tags:        []byte(`{"data": {"tags_desired": {}}}`),
			commitError: errors.New(""),
			statusCode:  http.StatusInternalServerError,
		},
		{
			context:             cloudAccount,
			tags:                []byte(`{"data": {"tags_desired": {}}}`),
			createAuditLogError: errors.New(""),
			statusCode:          http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		req, _ := http.NewRequest("PUT",
			"/api/v1/cloud_accounts/12345",
			bytes.NewReader(tc.tags))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		testApp := testhelpers2.NewTestApp()

		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.ContextCloudAccount, tc.context)
		ctx = context.WithValue(ctx, middleware.ContextCaller, auth.Caller{
			ID:   uuid.New(),
			Type: auth.CallerUser,
		})
		req = req.WithContext(ctx)

		testApp.Repository.On("Begin", mock.Anything).Return(testApp.Repository, tc.beginError)
		testApp.Repository.On("Commit", mock.Anything).Return(tc.commitError)
		testApp.Repository.On("UpdateCloudAccount", mock.Anything, mock.Anything).Return(cloudAccount, tc.updateError)
		testApp.Repository.On("CreateAuditLog", mock.Anything, mock.Anything).Return(db.AuditLog{}, tc.createAuditLogError)
		testApp.JS.On("PublishAsync", mock.Anything, mock.Anything, mock.Anything).Return(PubAckFuture{}, tc.publishError)
		testApp.Repository.On("Rollback", mock.Anything).Return(nil)

		V1UpdateCloudAccount(testApp.App)(w, req)

		_ = testhelpers2.AssertOpenAPI(t, w, req)

		result := w.Result()
		assert.Equal(t, tc.statusCode, result.StatusCode)
		_ = result.Body.Close()
	}
}

func TestV1GetCloudAccount(t *testing.T) {
	cloudAccount := testhelpers2.FactoryCloudAccount()

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
		testApp := testhelpers2.NewTestApp()

		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.ContextCloudAccount, tc.context)
		req = req.WithContext(ctx)

		V1GetCloudAccount(testApp.App)(w, req)

		_ = testhelpers2.AssertOpenAPI(t, w, req)

		result := w.Result()
		assert.Equal(t, tc.status, result.StatusCode)
		_ = result.Body.Close()
	}
}
