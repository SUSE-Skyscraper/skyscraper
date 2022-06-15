package main

import (
	"log"
	"net/http"

	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/middleware"
	"github.com/suse-skyscraper/skyscraper/internal/server"
)

func main() {
	config, err := application.Configuration()
	if err != nil {
		log.Fatal(err)
	}

	authorizer := middleware.AuthorizationHandler(config)
	cors := middleware.CorsHandler(config)

	http.Handle("/healthz", http.HandlerFunc(server.Health))
	http.Handle("/api/v1/profile", cors(authorizer(http.HandlerFunc(server.V1Profile))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
