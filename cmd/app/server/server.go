package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/spf13/cobra"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/fga"
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

			r.Route("/api/v1", func(r chi.Router) {
				apiAuthorizer := apimiddleware.AuthorizationHandler(app)

				r.Use(apiAuthorizer)

				r.Route("/caller", func(r chi.Router) {
					r.Get("/profile", server.V1CallerProfile(app))
					r.Get("/cloud_accounts", server.V1CallerCloudAccounts(app))
				})

				r.Route("/organizational_units", func(r chi.Router) {
					viewerEnforcer := apimiddleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationOrganizationalUnitsViewer)
					editorEnforcer := apimiddleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationOrganizationalUnitsEditor)

					// read actions
					r.Group(func(r chi.Router) {
						r.Use(viewerEnforcer)
						r.Get("/", server.V1ListOrganizationalUnits(app))
					})

					// write actions
					r.Group(func(r chi.Router) {
						r.Use(editorEnforcer)
						r.Post("/", server.V1CreateOrganizationalUnit(app))
					})

					r.Route("/{id}", func(r chi.Router) {
						organizationalUnitCtx := apimiddleware.OrganizationalUnitCtx(app)

						r.Use(organizationalUnitCtx)

						// read actions
						r.Group(func(r chi.Router) {
							r.Use(viewerEnforcer)
							r.Get("/", server.V1GetOrganizationalUnit(app))
						})

						// write actions
						r.Group(func(r chi.Router) {
							r.Use(editorEnforcer)
							r.Delete("/", server.V1DeleteOrganizationalUnit(app))
						})
					})
				})

				r.Route("/audit_logs", func(r chi.Router) {
					viewerEnforcer := apimiddleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationAuditLogViewer)

					// read actions
					r.Group(func(r chi.Router) {
						r.Use(viewerEnforcer)
						r.Get("/", server.V1ListAuditLogs(app))
					})
				})

				r.Route("/api_keys", func(r chi.Router) {
					viewerEnforcer := apimiddleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationAPIKeysViewer)
					editorEnforcer := apimiddleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationAPIKeysEditor)

					// read actions
					r.Group(func(r chi.Router) {
						r.Use(viewerEnforcer)
						r.Get("/", server.V1ListAPIKeys(app))
					})

					// write actions
					r.Group(func(r chi.Router) {
						r.Use(editorEnforcer)
						r.Post("/", server.V1CreateAPIKey(app))
					})

					r.Route("/{id}", func(r chi.Router) {
						apiKeyCtx := apimiddleware.APIKeyCtx(app)

						r.Use(apiKeyCtx)

						// read actions
						r.Group(func(r chi.Router) {
							r.Use(viewerEnforcer)
							r.Get("/", server.V1GetAPIKey(app))
						})
					})
				})

				r.Route("/standard_tags", func(r chi.Router) {
					viewerEnforcer := apimiddleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationStandardTagsViewer)
					editorEnforcer := apimiddleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationStandardTagsEditor)

					// read actions
					r.Group(func(r chi.Router) {
						r.Use(viewerEnforcer)
						r.Get("/", server.V1StandardTags(app))
					})

					// write actions
					r.Group(func(r chi.Router) {
						r.Use(editorEnforcer)
						r.Post("/", server.V1CreateStandardTag(app))
					})

					r.Route("/{id}", func(r chi.Router) {
						standardTagCtx := apimiddleware.TagCtx(app)

						r.Use(standardTagCtx)

						// write actions
						r.Group(func(r chi.Router) {
							r.Use(editorEnforcer)
							r.Put("/", server.V1UpdateStandardTag(app))
						})
					})
				})

				r.Route("/users", func(r chi.Router) {
					viewerEnforcer := apimiddleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationUsersViewer)

					// read actions
					r.Group(func(r chi.Router) {
						r.Use(viewerEnforcer)
						r.Get("/", server.V1Users(app))
					})

					r.Route("/{id}", func(r chi.Router) {
						userCtx := apimiddleware.UserCtx(app)

						r.Use(userCtx)

						// read actions
						r.Group(func(r chi.Router) {
							r.Use(viewerEnforcer)
							r.Get("/", server.V1User(app))
						})
					})
				})

				r.Route("/cloud_accounts", func(r chi.Router) {
					viewerEnforcer := apimiddleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationCloudAccountsViewer)

					// read actions
					r.Group(func(r chi.Router) {
						r.Use(viewerEnforcer)
						r.Get("/", server.V1ListCloudAccounts(app))
					})

					r.Route("/{id}", func(r chi.Router) {
						viewerEnforcer := apimiddleware.EnforcerHandler(app, fga.DocumentAccount, fga.DocumentAccountRelationViewer)
						editorEnforcer := apimiddleware.EnforcerHandler(app, fga.DocumentAccount, fga.DocumentAccountRelationEditor)
						cloudAccountCtx := apimiddleware.CloudAccountCtx(app)

						r.Use(cloudAccountCtx)

						// read actions
						r.Group(func(r chi.Router) {
							r.Use(viewerEnforcer)
							r.Get("/", server.V1GetCloudAccount(app))
						})

						// write actions
						r.Group(func(r chi.Router) {
							r.Use(editorEnforcer)
							r.Put("/", server.V1UpdateCloudAccount(app))
						})

						r.Route("/organizational_unit", func(r chi.Router) {
							// write actions
							r.Group(func(r chi.Router) {
								r.Use(editorEnforcer)
								r.Post("/", server.V1AssignCloudAccountToOU(app))
							})
						})
					})
				})

				r.Route("/cloud_tenants", func(r chi.Router) {
					organizationCloudTenantsViewerEnforcer := apimiddleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationCloudTenantsViewer)

					// read actions
					r.Group(func(r chi.Router) {
						r.Use(organizationCloudTenantsViewerEnforcer)
						r.Get("/", server.V1ListCloudTenants(app))
					})
				})
			})

			r.Route("/scim/v2", func(r chi.Router) {
				scimAuthorizer := scimmiddleware.BearerAuthorizationHandler(app)

				r.Use(scimAuthorizer)

				r.Get("/Users", scim.V2ListUsers(app))
				r.Post("/Users", scim.V2CreateUser(app))
				r.Route("/Users/{id}", func(r chi.Router) {
					scimUserCtx := scimmiddleware.UserCtx(app)

					r.Use(scimUserCtx)

					r.Get("/", scim.V2GetUser(app))
					r.Put("/", scim.V2UpdateUser(app))
					r.Patch("/", scim.V2PatchUser(app))
					r.Delete("/", scim.V2DeleteUser(app))
				})

				r.Get("/Groups", scim.V2ListGroups(app))
				r.Post("/Groups", scim.V2CreateGroup(app))
				r.Route("/Groups/{id}", func(r chi.Router) {
					scimGroupCtx := scimmiddleware.GroupCtx(app)

					r.Use(scimGroupCtx)

					r.Get("/", scim.V2GetGroup(app))
					r.Patch("/", scim.V2PatchGroup(app))
					r.Delete("/", scim.V2DeleteGroup(app))
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
