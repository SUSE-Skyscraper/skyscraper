package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/spf13/cobra"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/middleware"
	"github.com/suse-skyscraper/skyscraper/internal/server"
)

func NewCmd(app *application.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			r := chi.NewRouter()

			oktaAuthorizer := middleware.OktaAuthorizationHandler(app.Config)

			// common middleware
			r.Use(chimiddleware.Logger)
			r.Use(cors.Handler(cors.Options{
				AllowedOrigins:   []string{app.Config.Frontend.URL},
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
				AllowCredentials: true,
				MaxAge:           300,
			}))

			r.Get("/healthz", server.Health)

			// protected routes
			r.Group(func(r chi.Router) {
				if app.Config.Okta.Enabled {
					r.Use(oktaAuthorizer)
				}

				r.Get("/api/v1/profile", server.V1Profile)
				r.Get("/api/v1/cloud_tenants", server.V1CloudTenants(app))
				r.Get("/api/v1/cloud_tenants/cloud/{cloud}/tenant/{tenant_id}/accounts", server.V1CloudTenantAccounts(app))
				r.Get("/api/v1/cloud_tenants/cloud/{cloud}/tenant/{tenant_id}/accounts/{id}", server.V1CloudTenantAccount(app))
			})

			err := http.ListenAndServe(":8080", r)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
