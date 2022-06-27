package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/db"
)

type cloudAccountDecorator struct {
	CloudProvider     string      `json:"cloud_provider"`
	TenantID          string      `json:"tenant_id"`
	AccountID         string      `json:"account_id"`
	Name              string      `json:"name"`
	Active            bool        `json:"active"`
	TagsCurrent       interface{} `json:"tags_current"`
	TagsDesired       interface{} `json:"tags_desired"`
	TagsDriftDetected bool        `json:"tags_drift_detected"`
	CreatedAt         string      `json:"created_at"`
	UpdatedAt         string      `json:"updated_at"`
}

func V1CloudTenantAccounts(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cloudTenantID := chi.URLParam(r, "tenant_id")
		cloudProvider := chi.URLParam(r, "cloud")

		var cloudTenantAccounts, err = app.DB.GetCloudAllAccountMetadataForCloudAndTenant(
			r.Context(),
			db.GetCloudAllAccountMetadataForCloudAndTenantParams{
				Cloud:    cloudProvider,
				TenantID: cloudTenantID,
			})
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		var cloudAccountList []cloudAccountDecorator
		for _, account := range cloudTenantAccounts {
			decoratedCloudAccount := decorateCloudAccount(account)
			cloudAccountList = append(cloudAccountList, decoratedCloudAccount)
		}

		cloudAccountsJSON, err := json.Marshal(&cloudAccountList)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(cloudAccountsJSON)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
	}
}

func V1CloudTenantAccount(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cloudTenantID := chi.URLParam(r, "tenant_id")
		cloudProvider := chi.URLParam(r, "cloud")
		id := chi.URLParam(r, "id")

		var cloudTenantAccount, err = app.DB.GetCloudAccountMetadata(r.Context(), db.GetCloudAccountMetadataParams{
			Cloud:     cloudProvider,
			TenantID:  cloudTenantID,
			AccountID: id,
		})
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		decoratedCloudAccount := decorateCloudAccount(cloudTenantAccount)
		cloudAccountsJSON, err := json.Marshal(&decoratedCloudAccount)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(cloudAccountsJSON)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
	}
}

func decorateCloudAccount(account db.CloudAccountMetadatum) cloudAccountDecorator {
	return cloudAccountDecorator{
		CloudProvider:     account.Cloud,
		TenantID:          account.TenantID,
		AccountID:         account.AccountID,
		Name:              account.Name,
		Active:            account.Active,
		TagsCurrent:       account.TagsCurrent.Get(),
		TagsDesired:       account.TagsDesired.Get(),
		TagsDriftDetected: account.TagsDriftDetected,
		CreatedAt:         account.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:         account.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
