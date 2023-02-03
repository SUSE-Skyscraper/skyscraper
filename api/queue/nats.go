package queue

import (
	"encoding/json"
	"fmt"

	"github.com/suse-skyscraper/skyscraper/cli/application"
)

type PluginAction string

const (
	PluginActionTagUpdate PluginAction = "TAG_UPDATE"
)

type PluginPayload struct {
	Cloud      string       `json:"cloud"`
	TenantID   string       `json:"tenant_id"`
	ResourceID string       `json:"resource_id"`
	Action     PluginAction `json:"action"`
	Payload    interface{}  `json:"payload"`
}

type PluginWorker interface {
	PublishMessage(cloud string, payload PluginPayload) error
}

type DefaultPluginWorker struct {
	app *application.App
}

func NewPluginWorker(app *application.App) PluginWorker {
	return &DefaultPluginWorker{app: app}
}

func (w *DefaultPluginWorker) PublishMessage(cloud string, payload PluginPayload) error {
	if payload.Cloud == "" {
		return fmt.Errorf("cloud is required")
	}

	if payload.TenantID == "" {
		return fmt.Errorf("tenant_id is required")
	}

	if payload.ResourceID == "" {
		return fmt.Errorf("resource_id is required")
	}

	if payload.Action == "" {
		return fmt.Errorf("action is required")
	}

	workerPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	stream := fmt.Sprintf("PLUGIN.%s", cloud)
	_, err = w.app.JS.PublishAsync(stream, workerPayload)
	if err != nil {
		return err
	}

	return nil
}
