package main

import (
	"log"
	"net/http"

	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/middleware"
	"github.com/suse-skyscraper/skyscraper/internal/server"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	config, err := application.Configuration()
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	oktaAuthorizer := middleware.OktaAuthorizationHandler(config)

	// common middleware
	r.Use(chimiddleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{config.Frontend.URL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/healthz", server.Health)

	// protected routes
	r.Group(func(r chi.Router) {
		if config.Okta.Enabled {
			r.Use(oktaAuthorizer)
		}

		r.Get("/api/v1/profile", server.V1Profile)
	})

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
