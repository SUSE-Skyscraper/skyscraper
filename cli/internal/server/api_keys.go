package server

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/suse-skyscraper/skyscraper/api/payloads"
	"github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/cli/internal/auth/apikeys"
	"github.com/suse-skyscraper/skyscraper/cli/internal/db"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/auditors"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/middleware"

	"github.com/go-chi/render"
)

func V1ListAPIKeys(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKeys, err := app.Repo.GetAPIKeys(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, NewAPIKeysResponse(apiKeys))
	}
}

func V1GetAPIKey(_ *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, ok := r.Context().Value(middleware.ContextAPIKey).(db.ApiKey)
		if !ok {
			_ = render.Render(w, r, responses.ErrNotFound)
			return
		}

		_ = render.Render(w, r, NewAPIKeyResponse(apiKey, ""))
	}
}

func V1CreateAPIKey(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	apiKeyGenerator := apikeys.NewGenerator(app)

	return func(w http.ResponseWriter, r *http.Request) {
		// bind the payload
		var payload payloads.CreateAPIKeyPayload
		err := render.Bind(r, &payload)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInvalidRequest(err))
			return
		}

		// Begin a database transaction
		tx, err := app.PostgresPool.Begin(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}
		repo := app.Repo.WithTx(tx)

		// If any error occurs, rollback the transaction
		defer func(tx pgx.Tx, ctx context.Context) {
			_ = tx.Rollback(ctx)
		}(tx, r.Context())

		// create an auditor within our transaction
		auditor := auditors.NewAuditor(repo)

		// create an API key
		encodedHash, token, err := apiKeyGenerator.Generate()
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// persist the API key
		apiKey, err := app.Repo.InsertAPIKey(r.Context(), db.InsertAPIKeyParams{
			Owner:       payload.Data.Owner,
			Description: sql.NullString{String: payload.Data.Description, Valid: true},
			System:      false,
			Encodedhash: encodedHash,
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// audit the change
		err = auditor.AuditChange(r.Context(), db.AuditResourceTypeApiKey, apiKey.ID, payload)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// Commit the transaction
		err = tx.Commit(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, NewAPIKeyResponse(apiKey, token))
	}
}

func newAPIKeyItem(apiKey db.ApiKey, token string) responses.APIKeyItem {
	return responses.APIKeyItem{
		ID:   apiKey.ID.String(),
		Type: responses.ObjectResponseTypeAPIKey,
		Attributes: responses.APIKeyItemAttributes{
			Owner:       apiKey.Owner,
			Description: apiKey.Description.String,
			System:      apiKey.System,
			Token:       token,
			CreatedAt:   apiKey.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   apiKey.UpdatedAt.Format(time.RFC3339),
		},
	}
}

func NewAPIKeyResponse(apiKey db.ApiKey, token string) *responses.APIKeyResponse {
	return &responses.APIKeyResponse{
		Data: newAPIKeyItem(apiKey, token),
	}
}

func NewAPIKeysResponse(apiKeys []db.ApiKey) *responses.APIKeysResponse {
	list := make([]responses.APIKeyItem, len(apiKeys))
	for i, key := range apiKeys {
		list[i] = newAPIKeyItem(key, "")
	}

	return &responses.APIKeysResponse{
		Data: list,
	}
}
