package server

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/db"
	"github.com/suse-skyscraper/skyscraper/internal/server/responses"
)

func V1ListAuditLogs(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		targetID := r.URL.Query().Get("resource_id")
		targetType := r.URL.Query().Get("resource_type")

		if targetID == "" || targetType == "" {
			logs, users, err := app.Repository.GetAuditLogs(r.Context())
			if err != nil {
				_ = render.Render(w, r, responses.ErrInternalServerError)
				return
			}

			_ = render.Render(w, r, responses.NewAuditLogsListResponse(logs, users))
			return
		}

		target := db.AuditResourceType(targetType)
		id, err := uuid.Parse(targetID)
		if err != nil {
			_ = render.Render(w, r, responses.ErrNotFound)
			return
		}

		logs, users, err := app.Repository.GetAuditLogsForTarget(r.Context(), db.GetAuditLogsForTargetParams{
			ResourceID:   id,
			ResourceType: target,
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, responses.NewAuditLogsListResponse(logs, users))
	}
}
