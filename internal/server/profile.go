package server

import (
	"encoding/json"
	"net/http"

	"github.com/suse-skyscraper/skyscraper/internal/server/middleware"
)

type userProfile struct {
	Email string `json:"email"`
}

func V1Profile(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(middleware.UserEmail)

	// This should only happen if we didn't enable an authorizer such as Okta
	if email == nil {
		http.Error(w, "user not found", http.StatusInternalServerError)
		return
	}

	profile := userProfile{
		Email: email.(string),
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
