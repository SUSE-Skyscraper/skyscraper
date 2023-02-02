package server

import (
	"net/http"

	"github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/application"
	db2 "github.com/suse-skyscraper/skyscraper/cli/internal/db"
	responses2 "github.com/suse-skyscraper/skyscraper/cli/internal/server/responses"

	"github.com/go-chi/render"
	"github.com/google/uuid"
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

			_ = render.Render(w, r, responses2.NewAuditLogsListResponse(logs, users))
			return
		}

		target := db2.AuditResourceType(targetType)
		id, err := uuid.Parse(targetID)
		if err != nil {
			_ = render.Render(w, r, responses.ErrNotFound)
			return
		}

		logs, users, err := app.Repository.GetAuditLogsForTarget(r.Context(), db2.GetAuditLogsForTargetParams{
			ResourceID:   id,
			ResourceType: target,
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, responses2.NewAuditLogsListResponse(logs, users))
	}
}
