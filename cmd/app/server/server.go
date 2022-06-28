package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/spf13/cobra"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/server"
	middleware2 "github.com/suse-skyscraper/skyscraper/internal/server/middleware"
)

func NewCmd(app *application.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			r := chi.NewRouter()

			oktaAuthorizer := middleware2.OktaAuthorizationHandler(app.Config)
			cloudAccountCtx := middleware2.CloudAccountCtx(app)

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
				r.Route("/api/v1", func(r chi.Router) {
					r.Get("/profile", server.V1Profile)

					r.Route("/cloud_tenants", func(r chi.Router) {
						r.Get("/", server.V1CloudTenants(app))
						r.Route("/cloud/{cloud}/tenant/{tenant_id}/accounts", func(r chi.Router) {
							r.Get("/", server.V1ListCloudAccounts(app))

							r.Route("/{id}", func(r chi.Router) {
								r.Use(cloudAccountCtx)
								r.Get("/", server.V1GetCloudAccount(app))
								r.Put("/", server.V1UpdateCloudTenantAccount(app))
							})
						})
					})
				})
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
