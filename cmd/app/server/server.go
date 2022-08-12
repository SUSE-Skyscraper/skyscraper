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

			apiAuthorizer := apimiddleware.AuthorizationHandler(app)
			tagCtx := apimiddleware.TagCtx(app)
			userCtx := apimiddleware.UserCtx(app)
			apiKeyCtx := apimiddleware.APIKeyCtx(app)
			cloudAccountCtx := apimiddleware.CloudAccountCtx(app)

			organizationAuditLogViewerEnforcer := apimiddleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationAuditLogViewer)
			organizationAPIKeysViewerEnforcer := apimiddleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationAPIKeysViewer)
			organizationAPIKeysEditorEnforcer := apimiddleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationAPIKeysEditor)
			organizationStandardTagsViewerEnforcer := apimiddleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationStandardTagsViewer)
			organizationStandardTagsEditorEnforcer := apimiddleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationStandardTagsEditor)
			organizationUsersViewerEnforcer := apimiddleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationUsersViewer)
			organizationCloudAccountsViewerEnforcer := apimiddleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationCloudAccountsViewer)
			organizationCloudTenantsViewerEnforcer := apimiddleware.EnforcerHandler(app, fga.DocumentOrganization, fga.DocumentOrganizationRelationCloudTenantsViewer)

			cloudAccountViewerEnforcer := apimiddleware.EnforcerHandler(app, fga.DocumentAccount, fga.DocumentAccountRelationViewer)
			cloudAccountEditorEnforcer := apimiddleware.EnforcerHandler(app, fga.DocumentAccount, fga.DocumentAccountRelationEditor)

			r.Route("/api/v1", func(r chi.Router) {
				r.Use(apiAuthorizer)
				r.Get("/profile", server.V1Profile(app))

				r.Route("/audit_logs", func(r chi.Router) {
					// read actions
					r.Group(func(r chi.Router) {
						r.Use(organizationAuditLogViewerEnforcer)
						r.Get("/", server.V1ListAuditLogs(app))
					})
				})

				r.Route("/api_keys", func(r chi.Router) {
					// read actions
					r.Group(func(r chi.Router) {
						r.Use(organizationAPIKeysViewerEnforcer)
						r.Get("/", server.V1ListAPIKeys(app))
					})

					// write actions
					r.Group(func(r chi.Router) {
						r.Use(organizationAPIKeysEditorEnforcer)
						r.Post("/", server.V1CreateAPIKey(app))
					})

					r.Route("/{id}", func(r chi.Router) {
						r.Use(apiKeyCtx)

						// read actions
						r.Group(func(r chi.Router) {
							r.Use(organizationAPIKeysViewerEnforcer)
							r.Get("/", server.V1GetAPIKey(app))
						})
					})
				})

				r.Route("/standard_tags", func(r chi.Router) {
					// read actions
					r.Group(func(r chi.Router) {
						r.Use(organizationStandardTagsViewerEnforcer)
						r.Get("/", server.V1StandardTags(app))
					})

					// write actions
					r.Group(func(r chi.Router) {
						r.Use(organizationStandardTagsEditorEnforcer)
						r.Post("/", server.V1CreateStandardTag(app))
					})

					r.Route("/{id}", func(r chi.Router) {
						r.Use(tagCtx)

						// write actions
						r.Group(func(r chi.Router) {
							r.Use(organizationStandardTagsEditorEnforcer)
							r.Put("/", server.V1UpdateStandardTag(app))
						})
					})
				})

				r.Route("/users", func(r chi.Router) {
					// read actions
					r.Group(func(r chi.Router) {
						r.Use(organizationUsersViewerEnforcer)
						r.Get("/", server.V1Users(app))
					})

					r.Route("/{id}", func(r chi.Router) {
						r.Use(userCtx)

						// read actions
						r.Group(func(r chi.Router) {
							r.Use(organizationUsersViewerEnforcer)
							r.Get("/", server.V1User(app))
						})
					})
				})

				r.Route("/cloud_accounts", func(r chi.Router) {
					// read actions
					r.Group(func(r chi.Router) {
						r.Use(organizationCloudAccountsViewerEnforcer)
						r.Get("/", server.V1ListCloudAccounts(app))
					})

					r.Route("/{id}", func(r chi.Router) {
						r.Use(cloudAccountCtx)

						// read actions
						r.Group(func(r chi.Router) {
							r.Use(cloudAccountViewerEnforcer)
							r.Get("/", server.V1GetCloudAccount(app))
						})

						// write actions
						r.Group(func(r chi.Router) {
							r.Use(cloudAccountEditorEnforcer)
							r.Put("/", server.V1UpdateCloudAccount(app))
						})
					})
				})

				r.Route("/cloud_tenants", func(r chi.Router) {
					// read actions
					r.Group(func(r chi.Router) {
						r.Use(organizationCloudTenantsViewerEnforcer)
						r.Get("/", server.V1ListCloudTenants(app))
					})
				})
			})

			scimAuthorizer := scimmiddleware.BearerAuthorizationHandler(app)
			scimUserCtx := scimmiddleware.UserCtx(app)
			scimGroupCtx := scimmiddleware.GroupCtx(app)

			r.Route("/scim/v2", func(r chi.Router) {
				r.Use(scimAuthorizer)
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
