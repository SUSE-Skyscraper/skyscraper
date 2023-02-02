package server

import (
	"bytes"
	"context"
	"errors"
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

func TestV1ListOrganizationalUnits(t *testing.T) {
	organizationalUnit := testhelpers2.FactoryOrganizationalUnit()

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
		req, _ := http.NewRequest("GET", "/api/v1/organizational_units?cloud=AWS", nil)
		w := httptest.NewRecorder()
		testApp := testhelpers2.NewTestApp()

		testApp.Repository.
			On("GetOrganizationalUnits", mock.Anything).
			Return([]db.OrganizationalUnit{organizationalUnit}, tc.getError)

		V1ListOrganizationalUnits(testApp.App)(w, req)

		_ = testhelpers2.AssertOpenAPI(t, w, req)

		result := w.Result()
		assert.Equal(t, tc.statusCode, result.StatusCode)
		_ = result.Body.Close()
	}
}

func TestV1GetOrganizationalUnit(t *testing.T) {
	organizationalUnit := testhelpers2.FactoryOrganizationalUnit()

	tests := []struct {
		context    interface{}
		statusCode int
	}{
		{
			context:    interface{}(nil),
			statusCode: http.StatusInternalServerError,
		},
		{
			context:    organizationalUnit,
			statusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		req, _ := http.NewRequest("GET", "/api/v1/organizational_units/123456", nil)
		w := httptest.NewRecorder()
		testApp := testhelpers2.NewTestApp()

		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.ContextOrganizationalUnit, test.context)
		req = req.WithContext(ctx)

		V1GetOrganizationalUnit(testApp.App)(w, req)

		_ = testhelpers2.AssertOpenAPI(t, w, req)

		result := w.Result()
		assert.Equal(t, result.StatusCode, test.statusCode)
		_ = result.Body.Close()
	}
}

func TestV1CreateOrganizationalUnit(t *testing.T) {
	organizationalUnit := testhelpers2.FactoryOrganizationalUnit()

	tests := []struct {
		payload             []byte
		statusCode          int
		beginError          error
		commitError         error
		createAuditLogError error
		createError         error
		fgaError            error
	}{
		{
			payload:    []byte(`{"data": {"parent_id": "foobar", "display_name": "test"}}`),
			statusCode: http.StatusBadRequest,
		},
		{
			payload:    []byte(`{"data": {"parent_id": "07bdd057-6a50-42a7-b928-8af5eb59549f", "display_name": "test"}}`),
			statusCode: http.StatusCreated,
		},
		{
			payload:     []byte(`{"data": {"parent_id": "07bdd057-6a50-42a7-b928-8af5eb59549f", "display_name": "test"}}`),
			createError: errors.New(""),
			statusCode:  http.StatusInternalServerError,
		},
		{
			payload:    []byte(`{"data": {"parent_id": "07bdd057-6a50-42a7-b928-8af5eb59549f", "display_name": "test"}}`),
			beginError: errors.New(""),
			statusCode: http.StatusInternalServerError,
		},
		{
			payload:             []byte(`{"data": {"parent_id": "07bdd057-6a50-42a7-b928-8af5eb59549f", "display_name": "test"}}`),
			createAuditLogError: errors.New(""),
			statusCode:          http.StatusInternalServerError,
		},
		{
			payload:    []byte(`{"data": {"parent_id": "07bdd057-6a50-42a7-b928-8af5eb59549f", "display_name": "test"}}`),
			fgaError:   errors.New(""),
			statusCode: http.StatusInternalServerError,
		},
		{
			payload:     []byte(`{"data": {"parent_id": "07bdd057-6a50-42a7-b928-8af5eb59549f", "display_name": "test"}}`),
			commitError: errors.New(""),
			statusCode:  http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		req, _ := http.NewRequest("POST", "/api/v1/organizational_units", bytes.NewReader(tc.payload))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		testApp := testhelpers2.NewTestApp()

		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.ContextOrganizationalUnit, organizationalUnit)
		ctx = context.WithValue(ctx, middleware.ContextCaller, auth.Caller{
			ID:   uuid.New(),
			Type: auth.CallerUser,
		})
		req = req.WithContext(ctx)

		testApp.FGAClient.On("AddOrganizationalUnit", mock.Anything, mock.Anything, mock.Anything).Return(tc.fgaError)

		testApp.Repository.On("Begin", mock.Anything).Return(testApp.Repository, tc.beginError)
		testApp.Repository.On("Commit", mock.Anything).Return(tc.commitError)
		testApp.Repository.On("CreateOrganizationalUnit", mock.Anything, mock.Anything).Return(organizationalUnit, tc.createError)
		testApp.Repository.On("CreateAuditLog", mock.Anything, mock.Anything).Return(db.AuditLog{}, tc.createAuditLogError)
		testApp.Repository.On("Rollback", mock.Anything).Return(nil)

		V1CreateOrganizationalUnit(testApp.App)(w, req)

		_ = testhelpers2.AssertOpenAPI(t, w, req)

		result := w.Result()
		assert.Equal(t, tc.statusCode, result.StatusCode)
		_ = result.Body.Close()
	}
}

func TestV1DeleteOrganizationalUnit(t *testing.T) {
	organizationalUnit := testhelpers2.FactoryOrganizationalUnit()
	cloudAccount := testhelpers2.FactoryCloudAccount()

	tests := []struct {
		statusCode            int
		beginError            error
		commitError           error
		createAuditLogError   error
		deleteError           error
		getChildrenError      error
		getCloudAccountsError error
		fgaError              error
		children              []db.OrganizationalUnit
		accounts              []db.CloudAccount
		context               interface{}
	}{
		{
			context:    interface{}(nil),
			statusCode: http.StatusInternalServerError,
		},
		{
			context:    organizationalUnit,
			statusCode: http.StatusNoContent,
		},
		{
			context:    organizationalUnit,
			children:   []db.OrganizationalUnit{organizationalUnit},
			statusCode: http.StatusBadRequest,
		},
		{
			context:    organizationalUnit,
			accounts:   []db.CloudAccount{cloudAccount},
			statusCode: http.StatusBadRequest,
		},
		{
			context:     organizationalUnit,
			deleteError: errors.New(""),
			statusCode:  http.StatusInternalServerError,
		},
		{
			context:    organizationalUnit,
			beginError: errors.New(""),
			statusCode: http.StatusInternalServerError,
		},
		{
			context:          organizationalUnit,
			getChildrenError: errors.New(""),
			statusCode:       http.StatusInternalServerError,
		},
		{
			context:               organizationalUnit,
			getCloudAccountsError: errors.New(""),
			statusCode:            http.StatusInternalServerError,
		},
		{
			context:             organizationalUnit,
			createAuditLogError: errors.New(""),
			statusCode:          http.StatusInternalServerError,
		},
		{
			context:    organizationalUnit,
			fgaError:   errors.New(""),
			statusCode: http.StatusInternalServerError,
		},
		{
			context:     organizationalUnit,
			commitError: errors.New(""),
			statusCode:  http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		req, _ := http.NewRequest("DELETE", "/api/v1/organizational_units/123456", nil)
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		testApp := testhelpers2.NewTestApp()

		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.ContextOrganizationalUnit, tc.context)
		ctx = context.WithValue(ctx, middleware.ContextCaller, auth.Caller{
			ID:   uuid.New(),
			Type: auth.CallerUser,
		})
		req = req.WithContext(ctx)

		testApp.FGAClient.On("RemoveOrganizationalUnitRelationships", mock.Anything, mock.Anything, mock.Anything).Return(tc.fgaError)

		testApp.Repository.On("Begin", mock.Anything).Return(testApp.Repository, tc.beginError)
		testApp.Repository.On("Commit", mock.Anything).Return(tc.commitError)
		testApp.Repository.On("GetOrganizationalUnitChildren", mock.Anything, mock.Anything).Return(tc.children, tc.getChildrenError)
		testApp.Repository.On("GetOrganizationalUnitCloudAccounts", mock.Anything, mock.Anything).Return(tc.accounts, tc.getCloudAccountsError)
		testApp.Repository.On("DeleteOrganizationalUnit", mock.Anything, mock.Anything).Return(tc.deleteError)
		testApp.Repository.On("CreateAuditLog", mock.Anything, mock.Anything).Return(db.AuditLog{}, tc.createAuditLogError)
		testApp.Repository.On("Rollback", mock.Anything).Return(nil)

		V1DeleteOrganizationalUnit(testApp.App)(w, req)

		_ = testhelpers2.AssertOpenAPI(t, w, req)

		result := w.Result()
		assert.Equal(t, tc.statusCode, result.StatusCode)
		_ = result.Body.Close()
	}
}
