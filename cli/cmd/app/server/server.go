package server

import (
	"net/http"
	"time"

	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/cli/internal/fga"
	"github.com/suse-skyscraper/skyscraper/cli/internal/scimbridgedb"
	server2 "github.com/suse-skyscraper/skyscraper/cli/internal/server"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/spf13/cobra"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/bridge"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/router"
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

			r.Get("/healthz", server2.Health)

			r.Route("/api/v1", func(r chi.Router) {
				apiAuthorizer := middleware.AuthorizationHandler(app)

				r.Use(apiAuthorizer)

				r.Route("/caller", func(r chi.Router) {
					r.Get("/profile", server2.V1CallerProfile(app))
					r.Get("/cloud_accounts", server2.V1CallerCloudAccounts(app))
				})

				r.Route("/organizational_units", func(r chi.Router) {
					viewerEnforcer := middleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationOrganizationalUnitsViewer)
					editorEnforcer := middleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationOrganizationalUnitsEditor)

					// read actions
					r.Group(func(r chi.Router) {
						r.Use(viewerEnforcer)
						r.Get("/", server2.V1ListOrganizationalUnits(app))
					})

					// write actions
					r.Group(func(r chi.Router) {
						r.Use(editorEnforcer)
						r.Post("/", server2.V1CreateOrganizationalUnit(app))
					})

					r.Route("/{id}", func(r chi.Router) {
						organizationalUnitCtx := middleware.OrganizationalUnitCtx(app)

						r.Use(organizationalUnitCtx)

						// read actions
						r.Group(func(r chi.Router) {
							r.Use(viewerEnforcer)
							r.Get("/", server2.V1GetOrganizationalUnit(app))
						})

						// write actions
						r.Group(func(r chi.Router) {
							r.Use(editorEnforcer)
							r.Delete("/", server2.V1DeleteOrganizationalUnit(app))
						})
					})
				})

				r.Route("/audit_logs", func(r chi.Router) {
					viewerEnforcer := middleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationAuditLogViewer)

					// read actions
					r.Group(func(r chi.Router) {
						r.Use(viewerEnforcer)
						r.Get("/", server2.V1ListAuditLogs(app))
					})
				})

				r.Route("/api_keys", func(r chi.Router) {
					viewerEnforcer := middleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationAPIKeysViewer)
					editorEnforcer := middleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationAPIKeysEditor)

					// read actions
					r.Group(func(r chi.Router) {
						r.Use(viewerEnforcer)
						r.Get("/", server2.V1ListAPIKeys(app))
					})

					// write actions
					r.Group(func(r chi.Router) {
						r.Use(editorEnforcer)
						r.Post("/", server2.V1CreateAPIKey(app))
					})

					r.Route("/{id}", func(r chi.Router) {
						apiKeyCtx := middleware.APIKeyCtx(app)

						r.Use(apiKeyCtx)

						// read actions
						r.Group(func(r chi.Router) {
							r.Use(viewerEnforcer)
							r.Get("/", server2.V1GetAPIKey(app))
						})
					})
				})

				r.Route("/standard_tags", func(r chi.Router) {
					viewerEnforcer := middleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationStandardTagsViewer)
					editorEnforcer := middleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationStandardTagsEditor)

					// read actions
					r.Group(func(r chi.Router) {
						r.Use(viewerEnforcer)
						r.Get("/", server2.V1StandardTags(app))
					})

					// write actions
					r.Group(func(r chi.Router) {
						r.Use(editorEnforcer)
						r.Post("/", server2.V1CreateStandardTag(app))
					})

					r.Route("/{id}", func(r chi.Router) {
						standardTagCtx := middleware.TagCtx(app)

						r.Use(standardTagCtx)

						// write actions
						r.Group(func(r chi.Router) {
							r.Use(editorEnforcer)
							r.Put("/", server2.V1UpdateStandardTag(app))
						})
					})
				})

				r.Route("/users", func(r chi.Router) {
					viewerEnforcer := middleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationUsersViewer)

					// read actions
					r.Group(func(r chi.Router) {
						r.Use(viewerEnforcer)
						r.Get("/", server2.V1Users(app))
					})

					r.Route("/{id}", func(r chi.Router) {
						userCtx := middleware.UserCtx(app)

						r.Use(userCtx)

						// read actions
						r.Group(func(r chi.Router) {
							r.Use(viewerEnforcer)
							r.Get("/", server2.V1User(app))
						})
					})
				})

				r.Route("/groups/{group}/tenants", func(r chi.Router) {
					r.Route("/{tenant_id}", func(r chi.Router) {
						// organizationCloudTenantsViewerEnforcer := middleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationCloudTenantsViewer)
						organizationCloudTenantsEditorEnforcer := middleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationCloudTenantsEditor)
						tenantCtx := middleware.TenantCtx(app)

						// write actions
						r.Group(func(r chi.Router) {
							r.Use(organizationCloudTenantsEditorEnforcer)

							r.Put("/", server2.V1CreateOrUpdateTenants(app))
						})

						r.Group(func(r chi.Router) {
							r.Use(tenantCtx)

							r.Route("/resources", func(r chi.Router) {
								viewerEnforcer := middleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationCloudAccountsViewer)

								// read actions
								r.Group(func(r chi.Router) {
									r.Use(viewerEnforcer)
									r.Get("/", server2.V1ListResources(app))
								})

								r.Route("/{resource_id}", func(r chi.Router) {
									viewerEnforcer := middleware.EnforcerHandler(app, fga.DocumentAccount, fga.DocumentAccountRelationViewer)
									editorEnforcer := middleware.EnforcerHandler(app, fga.DocumentAccount, fga.DocumentAccountRelationEditor)
									resourceCtx := middleware.ResourceCtx(app)

									// write actions
									r.Group(func(r chi.Router) {
										r.Use(editorEnforcer)
										r.Put("/", server2.V1CreateOrUpdateResource(app))
									})

									// read actions
									r.Group(func(r chi.Router) {
										r.Use(resourceCtx)
										r.Use(viewerEnforcer)
										r.Get("/", server2.V1GetResource(app))
									})

									r.Route("/organizational_unit", func(r chi.Router) {
										r.Use(resourceCtx)
										r.Use(editorEnforcer)

										// write actions
										r.Group(func(r chi.Router) {
											r.Post("/", server2.V1AssignCloudAccountToOU(app))
										})
									})
								})
							})
						})
					})
				})

				r.Route("/cloud_tenants", func(r chi.Router) {
					organizationCloudTenantsViewerEnforcer := middleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationCloudTenantsViewer)

					// read actions
					r.Group(func(r chi.Router) {
						r.Use(organizationCloudTenantsViewerEnforcer)
						r.Get("/", server2.V1ListCloudTenants(app))
					})
				})
			})

			scimAuthorizer := middleware.BearerAuthorizationHandler(app)
			db := scimbridgedb.New(app)
			b := bridge.New(&db, app.Config.ServerConfig.BaseURL)
			router.Hook(r, &b, scimAuthorizer)

			s := &http.Server{
				Addr:         ":8080",
				Handler:      r,
				ReadTimeout:  2 * time.Second,
				WriteTimeout: 2 * time.Second,
			}
			err := s.ListenAndServe()
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
