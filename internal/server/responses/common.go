package responses

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/suse-skyscraper/skyscraper/internal/db"
)

type ObjectResponseType string

const (
	ObjectResponseTypeUnknown      ObjectResponseType = ""
	ObjectResponseTypeUser         ObjectResponseType = "user"
	ObjectResponseTypeAuditLog     ObjectResponseType = "audit_log"
	ObjectResponseTypeCloudAccount ObjectResponseType = "cloud_account"
	ObjectResponseTypeCloudTenant  ObjectResponseType = "cloud_tenant"
	ObjectResponseTypeTag          ObjectResponseType = "tag"
	ObjectResponseTypeAPIKey       ObjectResponseType = "api_key"
)

func parseDBCallerType(callerType db.CallerType) ObjectResponseType {
	switch callerType {
	case db.CallerTypeUser:
		return ObjectResponseTypeUser
	case db.CallerTypeApiKey:
		return ObjectResponseTypeAPIKey
	default:
		return ObjectResponseTypeUnknown
	}
}

type RelationshipData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type Relationship struct {
	RelationshipData RelationshipData `json:"data"`
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}
var ErrInternalServerError = &ErrResponse{HTTPStatusCode: 500, StatusText: "internal error"}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}
