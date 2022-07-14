package server

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/suse-skyscraper/skyscraper/internal/db"
	"github.com/suse-skyscraper/skyscraper/internal/server/middleware"
	"github.com/suse-skyscraper/skyscraper/internal/server/responses"
)

func V1Profile(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.User).(db.User)
	if !ok {
		_ = render.Render(w, r, responses.ErrInternalServerError)
		return
	}

	_ = render.Render(w, r, responses.NewUserResponse(user))
}
