package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/suse-skyscraper/skyscraper/cli/internal/auth"
	"github.com/suse-skyscraper/skyscraper/cli/internal/db"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/middleware"
	"github.com/suse-skyscraper/skyscraper/cli/internal/testhelpers"
)

func TestV1CreateOrUpdateResource(t *testing.T) {
	cloudAccount := testhelpers.FactoryCloudAccount()
	tenant := testhelpers.FactoryTenant()

	type PubAckFuture struct {
		nats.PubAckFuture
	}

	tests := []struct {
		description         string
		resourceID          string
		body                []byte
		updateError         error
		beginError          error
		commitError         error
		createAuditLogError error
		publishToNatsError  error
		fgaError            error
		statusCode          int
		account             interface{}
		ctxTenant           interface{}
	}{
		{
			description: "success",
			resourceID:  "12345",
			account:     cloudAccount,
			ctxTenant:   tenant,
			body:        []byte(`{"data": {"tags_desired": {}}}`),
			statusCode:  http.StatusOK,
		},
		{
			description: "success with nats message",
			resourceID:  "12345",
			account: db.CloudAccount{
				TagsCurrent: pgtype.JSONB{Bytes: []byte("{\"bar\": \"foo\"}"), Status: pgtype.Present},
				TagsDesired: pgtype.JSONB{Bytes: []byte("{\"foo\": \"bar\"}"), Status: pgtype.Present},
			},
			ctxTenant:  tenant,
			body:       []byte(`{"data": {"tags_desired": {"foo": "bar"}, "tags_current": {"bar": "foo"}}}`),
			statusCode: http.StatusOK,
		},
		{
			description: "bad request",
			resourceID:  "12345",
			account:     cloudAccount,
			ctxTenant:   tenant,
			body:        []byte(`asd1!!{]`),
			statusCode:  http.StatusBadRequest,
		},
		{
			description: "no resource id",
			statusCode:  http.StatusBadRequest,
		},
		{
			description: "update error",
			resourceID:  "12345",
			account:     cloudAccount,
			ctxTenant:   tenant,
			body:        []byte(`{"data": {"tags_desired": {}}}`),
			updateError: errors.New(""),
			statusCode:  http.StatusInternalServerError,
		},
		{
			description: "fga error",
			resourceID:  "12345",
			account:     cloudAccount,
			ctxTenant:   tenant,
			body:        []byte(`{"data": {"tags_desired": {}}}`),
			fgaError:    errors.New(""),
			statusCode:  http.StatusInternalServerError,
		},
		{
			description:        "publish error",
			resourceID:         "12345",
			account:            cloudAccount,
			ctxTenant:          tenant,
			body:               []byte(`{"data": {"tags_desired": {}}}`),
			publishToNatsError: errors.New(""),
			statusCode:         http.StatusOK,
		},
		{
			description: "tx begin error",
			resourceID:  "12345",
			account:     cloudAccount,
			ctxTenant:   tenant,
			body:        []byte(`{"data": {"tags_desired": {}}}`),
			beginError:  errors.New(""),
			statusCode:  http.StatusInternalServerError,
		},
		{
			description: "invalid account context",
			resourceID:  "12345",
			ctxTenant:   interface{}(nil),
			body:        []byte(`{"data": {"tags_desired": {}}}`),
			statusCode:  http.StatusInternalServerError,
		},
		{
			description: "tx commit error",
			resourceID:  "12345",
			ctxTenant:   tenant,
			account:     cloudAccount,
			body:        []byte(`{"data": {"tags_desired": {}}}`),
			commitError: errors.New(""),
			statusCode:  http.StatusInternalServerError,
		},
		{
			description:         "audit log error",
			resourceID:          "12345",
			ctxTenant:           tenant,
			account:             cloudAccount,
			body:                []byte(`{"data": {"tags_desired": {}}}`),
			createAuditLogError: errors.New(""),
			statusCode:          http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		req, _ := http.NewRequest("PUT",
			"/api/v1/groups/AWS/tenants/tenant1234/resources/12345",
			bytes.NewReader(tc.body))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		testApp, err := testhelpers.NewTestApp()
		if err != nil {
			t.Fatal(err)
		}

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("resource_id", tc.resourceID)

		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.ContextCloudAccount, tc.account)
		ctx = context.WithValue(ctx, middleware.ContextTenant, tc.ctxTenant)
		ctx = context.WithValue(ctx, middleware.ContextCaller, auth.Caller{
			ID:   uuid.New(),
			Type: auth.CallerUser,
		})
		req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

		testApp.PostgresPool.ExpectBegin().WillReturnError(tc.beginError)
		if tc.beginError == nil && tc.createAuditLogError == nil && tc.updateError == nil {
			testApp.PostgresPool.ExpectCommit().WillReturnError(tc.commitError)
		}
		testApp.Repo.On("WithTx", mock.Anything).Return(testApp.Repo)
		testApp.Repo.On("CreateOrUpdateCloudAccount", mock.Anything, mock.Anything).Return(tc.account, tc.updateError)
		testApp.Repo.On("CreateAuditLog", mock.Anything, mock.Anything).Return(db.AuditLog{}, tc.createAuditLogError)
		testApp.FGAClient.On("AddAccountToOrganization", mock.Anything, mock.Anything).Return(tc.fgaError)
		testApp.JS.On("PublishAsync", mock.Anything, mock.Anything, mock.Anything).Return(PubAckFuture{}, tc.publishToNatsError)
		testApp.PostgresPool.ExpectRollback()

		V1CreateOrUpdateResource(testApp.App)(w, req)

		_ = testhelpers.AssertOpenAPI(t, w, req)

		result := w.Result()
		assert.Equal(t, tc.statusCode, result.StatusCode, fmt.Sprintf("status should match for %s", tc.description))
		_ = result.Body.Close()

		testApp.Close() // Close the test app, we cannot defer in a loop
	}
}
