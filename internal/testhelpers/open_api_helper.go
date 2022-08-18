package testhelpers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getkin/kin-openapi/routers/gorillamux"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	oarouters "github.com/getkin/kin-openapi/routers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var router oarouters.Router

func init() {
	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}

	doc, err := loader.LoadFromFile("./../../api/skyscraper.yaml")
	if err != nil {
		panic(err)
	}
	// Our test requests are relative, so the server URL doesn't get found.
	doc.Servers = nil
	err = doc.Validate(ctx)
	if err != nil {
		panic(err)
	}

	router, err = gorillamux.NewRouter(doc)
	if err != nil {
		panic(err)
	}
}

func AssertOpenAPI(t *testing.T, rr *httptest.ResponseRecorder, req *http.Request) []byte {
	t.Helper()
	ctx := context.Background()

	// The request body exists and was already read once when the request
	// was sent. Replay it to allow ValidateRequest() to read it again.
	if req.Body != nil && req.Body != http.NoBody {
		req.Body, _ = req.GetBody()
	}

	// Validate request.
	route, pathParams, err := router.FindRoute(req)
	require.NoError(t, err, "could not find route")
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    req,
		PathParams: pathParams,
		Route:      route,
		Options: &openapi3filter.Options{
			AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
		},
	}

	// Don't validate the request body if the response indicates an error,
	// to allow testing server-side validation using known bad values.
	result := rr.Result()

	if result.StatusCode >= http.StatusBadRequest {
		requestValidationInput.Options.ExcludeRequestBody = true
	}

	err = openapi3filter.ValidateRequest(ctx, requestValidationInput)
	assert.NoError(t, err, "http request is not valid")

	// Validate response.
	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput,
		Status:                 result.StatusCode,
		//Header:                 result.Header,
		Header: http.Header{"Content-Type": []string{rr.Header().Get("Content-Type")}},
		Options: &openapi3filter.Options{
			IncludeResponseStatus: true,
		},
	}
	bodyBytes := rr.Body.Bytes()
	responseValidationInput.SetBodyBytes(bodyBytes)
	err = openapi3filter.ValidateResponse(ctx, responseValidationInput)
	assert.NoError(t, err, "http response is not valid")

	_ = result.Body.Close()

	return bodyBytes
}
