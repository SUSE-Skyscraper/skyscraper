package server

import (
	"encoding/json"
	"net/http"

	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/db"
)

type cloudTenantDecorator struct {
	CloudProvider string `json:"cloud_provider"`
	TenantID      string `json:"tenant_id"`
	Name          string `json:"name"`
	Active        bool   `json:"active"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

func V1CloudTenants(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cloudTenants, err := app.DB.GetCloudTenants(r.Context())
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		var cloudTenantList []cloudTenantDecorator
		for _, tenant := range cloudTenants {
			decoratedCloudTenant := decorateCloudTenant(tenant)
			cloudTenantList = append(cloudTenantList, decoratedCloudTenant)
		}

		cloudTenantsJSON, err := json.Marshal(&cloudTenantList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(cloudTenantsJSON)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func decorateCloudTenant(tenant db.CloudTenant) cloudTenantDecorator {
	return cloudTenantDecorator{
		CloudProvider: tenant.Cloud,
		TenantID:      tenant.TenantID,
		Name:          tenant.Name,
		Active:        tenant.Active,
		CreatedAt:     tenant.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     tenant.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
