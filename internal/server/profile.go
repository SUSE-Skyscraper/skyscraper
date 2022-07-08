package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/suse-skyscraper/skyscraper/internal/db"
	"github.com/suse-skyscraper/skyscraper/internal/server/middleware"
	"github.com/suse-skyscraper/skyscraper/internal/server/responses"
)

type userProfile struct {
	Email string `json:"email"`
}

func V1Profile(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.User).(db.User)
	if !ok {
		_ = render.Render(w, r, responses.ErrInternalServerError)
		return
	}

	profile := userProfile{
		Email: user.Username,
	}
	profileJSON, err := json.Marshal(&profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(profileJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
