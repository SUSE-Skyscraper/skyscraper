package server

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/suse-skyscraper/skyscraper/cli/internal/auth"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/middleware"
	testhelpers2 "github.com/suse-skyscraper/skyscraper/cli/internal/testhelpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestV1AssignCloudAccountToOU(t *testing.T) {
	cloudAccount := testhelpers2.FactoryCloudAccount()

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
		testApp.Repository.On("UnAssignCloudAccountFromOrganizationalUnits", mock.Anything, mock.Anything).Return(tc.unAssignError)
		testApp.Repository.On("AssignCloudAccountToOrganizationalUnit", mock.Anything, mock.Anything, mock.Anything).Return(tc.assignError)
		testApp.Repository.On("Rollback", mock.Anything).Return(nil)

		V1AssignCloudAccountToOU(testApp.App)(w, req)

		_ = testhelpers2.AssertOpenAPI(t, w, req)

		result := w.Result()
		assert.Equal(t, tc.statusCode, result.StatusCode)
		_ = result.Body.Close()
	}
}
