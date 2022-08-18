package server

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/suse-skyscraper/skyscraper/internal/auth"
	"github.com/suse-skyscraper/skyscraper/internal/db"
	"github.com/suse-skyscraper/skyscraper/internal/server/middleware"
	"github.com/suse-skyscraper/skyscraper/internal/testhelpers"
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
		testApp := testhelpers.NewTestApp()

		testApp.Repository.
			On("SearchCloudAccounts", mock.Anything, mock.Anything).
			Return([]db.CloudAccount{cloudAccount}, tc.getError)

		V1ListCloudAccounts(testApp.App)(w, req)

		_ = testhelpers.AssertOpenAPI(t, w, req)

		result := w.Result()
		assert.Equal(t, tc.statusCode, result.StatusCode)
		_ = result.Body.Close()
	}
}

func TestV1AssignCloudAccountToOU(t *testing.T) {
	cloudAccount := testhelpers.FactoryCloudAccount()

	tests := []struct {
		payload       []byte
		statusCode    int
		context       interface{}
		beginError    error
		commitError   error
		assignError   error
		unAssignError error
	}{
		{
			context:    interface{}(nil),
			payload:    []byte(`{"data": {"organizational_unit_id": "e0f99cfa-7906-4b5b-ac05-e312abf4785e"}}`),
			statusCode: http.StatusInternalServerError,
		},
		{
			context:    cloudAccount,
			payload:    []byte(`{"data": {"organizational_unit_id": "e0f99cfa-7906-4b5b-ac05-e312abf4785e"}}`),
			statusCode: http.StatusNoContent,
		},
		{
			context:    cloudAccount,
			payload:    []byte(`{"data": {"organizational_unit_id": "e0f99cfa-7906-4b5b-ac05-e312abf4785e"}}`),
			beginError: fmt.Errorf(""),
			statusCode: http.StatusInternalServerError,
		},
		{
			context:     cloudAccount,
			payload:     []byte(`{"data": {"organizational_unit_id": "e0f99cfa-7906-4b5b-ac05-e312abf4785e"}}`),
			commitError: fmt.Errorf(""),
			statusCode:  http.StatusInternalServerError,
		},
		{
			context:     cloudAccount,
			payload:     []byte(`{"data": {"organizational_unit_id": "e0f99cfa-7906-4b5b-ac05-e312abf4785e"}}`),
			assignError: fmt.Errorf(""),
			statusCode:  http.StatusInternalServerError,
		},
		{
			context:       cloudAccount,
			payload:       []byte(`{"data": {"organizational_unit_id": "e0f99cfa-7906-4b5b-ac05-e312abf4785e"}}`),
			unAssignError: fmt.Errorf(""),
			statusCode:    http.StatusInternalServerError,
		},
		{
			context:    cloudAccount,
			payload:    []byte(`{"data": {"organizational_unit_id": "123456"}}`),
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		req, _ := http.NewRequest("POST", "/api/v1/cloud_accounts/12345/organizational_unit", bytes.NewReader(tc.payload))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		testApp := testhelpers.NewTestApp()

		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.ContextCloudAccount, tc.context)
		ctx = context.WithValue(ctx, middleware.ContextCaller, auth.Caller{
			ID:   uuid.New(),
			Type: auth.CallerUser,
		})
		req = req.WithContext(ctx)

		testApp.Repository.On("Begin", mock.Anything).Return(testApp.Repository, tc.beginError)
		testApp.Repository.On("Commit", mock.Anything).Return(tc.commitError)
		testApp.Repository.On("UnAssignCloudAccountFromOrganizationalUnits", mock.Anything, mock.Anything).Return(tc.unAssignError)
		testApp.Repository.On("AssignCloudAccountToOrganizationalUnit", mock.Anything, mock.Anything, mock.Anything).Return(tc.assignError)
		testApp.Repository.On("Rollback", mock.Anything).Return(nil)

		V1AssignCloudAccountToOU(testApp.App)(w, req)

		_ = testhelpers.AssertOpenAPI(t, w, req)

		result := w.Result()
		assert.Equal(t, tc.statusCode, result.StatusCode)
		_ = result.Body.Close()
	}
}

func TestV1UpdateCloudAccount(t *testing.T) {
	cloudAccount := testhelpers.FactoryCloudAccount()

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
		testApp := testhelpers.NewTestApp()

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

		_ = testhelpers.AssertOpenAPI(t, w, req)

		result := w.Result()
		assert.Equal(t, tc.statusCode, result.StatusCode)
		_ = result.Body.Close()
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
		testApp := testhelpers.NewTestApp()

		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.ContextCloudAccount, tc.context)
		req = req.WithContext(ctx)

		V1GetCloudAccount(testApp.App)(w, req)

		_ = testhelpers.AssertOpenAPI(t, w, req)

		result := w.Result()
		assert.Equal(t, tc.status, result.StatusCode)
		_ = result.Body.Close()
	}
}
