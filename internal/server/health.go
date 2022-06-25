package server

import (
	"net/http"
)

func Health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte("{}"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
