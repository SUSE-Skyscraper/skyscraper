package server

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/server/responses"
)

func V1CloudTenants(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cloudTenants, err := app.DB.GetCloudTenants(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		err = render.RenderList(w, r, responses.NewCloudTenantListResponse(cloudTenants))
		if err != nil {
			_ = render.Render(w, r, responses.ErrRender(err))
			return
		}
	}
}
