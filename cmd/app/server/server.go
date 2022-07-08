package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/spf13/cobra"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/scim"
	scimmiddleware "github.com/suse-skyscraper/skyscraper/internal/scim/middleware"
	"github.com/suse-skyscraper/skyscraper/internal/server"
	apimiddleware "github.com/suse-skyscraper/skyscraper/internal/server/middleware"
)

func NewCmd(app *application.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			r := chi.NewRouter()

			err := app.StartEnforcer()
			if err != nil {
				return err
			}

			enforcerHandler := apimiddleware.EnforcerHandler(app)
			oktaAuthorizer := apimiddleware.OktaAuthorizationHandler(app)
			scimAuthorizer := scimmiddleware.BearerAuthorizationHandler(app)
			cloudAccountCtx := apimiddleware.CloudAccountCtx(app)
			scimUserCtx := scimmiddleware.UserCtx(app)
			scimGroupCtc := scimmiddleware.GroupCtx(app)

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
				r.Use(enforcerHandler)
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

			r.Group(func(r chi.Router) {
				r.Use(scimAuthorizer)
				r.Route("/scim/v2", func(r chi.Router) {
					r.Get("/Users", scim.V2ListUsers(app))
					r.Post("/Users", scim.V2CreateUser(app))
					r.Route("/Users/{id}", func(r chi.Router) {
						r.Use(scimUserCtx)
						r.Get("/", scim.V2GetUser(app))
						r.Put("/", scim.V2UpdateUser(app))
						r.Patch("/", scim.V2PatchUser(app))
						r.Delete("/", scim.V2DeleteUser(app))
					})

					r.Get("/Groups", scim.V2ListGroups(app))
					r.Post("/Groups", scim.V2CreateGroup(app))
					r.Route("/Groups/{id}", func(r chi.Router) {
						r.Use(scimGroupCtc)
						r.Get("/", scim.V2GetGroup(app))
						// r.Put("/", scim.V2UpdateUser(app))
						r.Patch("/", scim.V2PatchGroup(app))
						r.Delete("/", scim.V2DeleteGroup(app))
					})
					/*
						r.Get("/ServiceProviderConfig", server.V2GetServiceProviderConfig(app))
						r.Get("/Schemas", server.V2ListSchemas(app))
						r.Get("/Schemas/{id}", server.V2GetSchema(app))
						r.Get("/ResourceTypes", server.V2ListResourceTypes(app))
						r.Get("/ResourceTypes/{id}", server.V2GetResourceType(app))
					*/
				})
			})

			err = http.ListenAndServe(":8080", r)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
