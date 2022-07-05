package workers

import (
	"encoding/json"

	"github.com/suse-skyscraper/skyscraper/internal/application"
)

type ChangeTagsPayload struct {
	Cloud       string `json:"cloud"`
	TenantID    string `json:"tenant_id"`
	AccountID   string `json:"account_id"`
	AccountName string `json:"account_name"`
}

type NatsWorker interface {
	PublishTagChange(payload ChangeTagsPayload) error
}

type DefaultNatsWorker struct {
	app *application.App
}

func NewWorker(app *application.App) NatsWorker {
	return &DefaultNatsWorker{app: app}
}

func (w *DefaultNatsWorker) PublishTagChange(payload ChangeTagsPayload) error {
	workerPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = w.app.JS.PublishAsync("TAGS.change", workerPayload)
	if err != nil {
		return err
	}

	return nil
}
