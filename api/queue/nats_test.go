package queue

import (
	"fmt"
	"testing"

	"github.com/suse-skyscraper/skyscraper/test/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/suse-skyscraper/skyscraper/test/helpers"
)

func TestDefaultPluginWorker_PublishMessage(t *testing.T) {
	tt := []struct {
		name          string
		cloud         string
		payload       PluginPayload
		errorExpected bool
		publishError  error
	}{
		{
			name:  "correct",
			cloud: "aws",
			payload: PluginPayload{
				Cloud:      "aws",
				TenantID:   "tenant",
				ResourceID: "resource",
				Action:     PluginActionTagUpdate,
				Payload:    map[string]interface{}{},
			},
			errorExpected: false,
			publishError:  nil,
		},
		{
			name:  "missing payload",
			cloud: "aws",
			payload: PluginPayload{
				Cloud:      "aws",
				TenantID:   "tenant",
				ResourceID: "resource",
				Action:     PluginActionTagUpdate,
				Payload:    nil,
			},
			errorExpected: false,
			publishError:  nil,
		},
		{
			name:  "missing cloud in call",
			cloud: "",
			payload: PluginPayload{
				Cloud:      "aws",
				TenantID:   "tenant",
				ResourceID: "resource",
				Action:     PluginActionTagUpdate,
				Payload:    map[string]interface{}{},
			},
			errorExpected: true,
			publishError:  nil,
		},
		{
			name:  "missing cloud in payload",
			cloud: "aws",
			payload: PluginPayload{
				Cloud:      "",
				TenantID:   "tenant",
				ResourceID: "resource",
				Action:     PluginActionTagUpdate,
				Payload:    map[string]interface{}{},
			},
			errorExpected: true,
			publishError:  nil,
		},
		{
			name:  "missing tenant_id in payload",
			cloud: "aws",
			payload: PluginPayload{
				Cloud:      "aws",
				TenantID:   "",
				ResourceID: "resource",
				Action:     PluginActionTagUpdate,
				Payload:    map[string]interface{}{},
			},
			errorExpected: true,
			publishError:  nil,
		},
		{
			name:  "missing tenant_id in payload",
			cloud: "aws",
			payload: PluginPayload{
				Cloud:      "aws",
				TenantID:   "",
				ResourceID: "resource",
				Action:     PluginActionTagUpdate,
				Payload:    map[string]interface{}{},
			},
			errorExpected: true,
			publishError:  nil,
		},
		{
			name:  "missing resource_id in payload",
			cloud: "aws",
			payload: PluginPayload{
				Cloud:      "aws",
				TenantID:   "tenant",
				ResourceID: "",
				Action:     PluginActionTagUpdate,
				Payload:    map[string]interface{}{},
			},
			errorExpected: true,
			publishError:  nil,
		},
		{
			name:  "missing action in payload",
			cloud: "aws",
			payload: PluginPayload{
				Cloud:      "aws",
				TenantID:   "tenant",
				ResourceID: "resource",
				Action:     "",
				Payload:    map[string]interface{}{},
			},
			errorExpected: true,
			publishError:  nil,
		},
		{
			name:  "publish error",
			cloud: "aws",
			payload: PluginPayload{
				Cloud:      "aws",
				TenantID:   "tenant",
				ResourceID: "resource",
				Action:     PluginActionTagUpdate,
				Payload:    map[string]interface{}{},
			},
			errorExpected: true,
			publishError:  fmt.Errorf("publish error"),
		},
		{
			name:  "payload error",
			cloud: "aws",
			payload: PluginPayload{
				Cloud:      "aws",
				TenantID:   "tenant",
				ResourceID: "resource",
				Action:     PluginActionTagUpdate,
				Payload:    make(chan int),
			},
			errorExpected: true,
			publishError:  nil,
		},
	}

	for _, tc := range tt {
		testApp, err := helpers.NewMockedApp()
		if err != nil {
			t.Fatal(err)
		}

		future := new(mocks.MockPubAckFuture)
		testApp.JS.On("PublishAsync", mock.Anything, mock.Anything, mock.Anything).Return(future, tc.publishError)

		worker := NewPluginWorker(testApp.App)
		err = worker.PublishMessage(tc.cloud, tc.payload)

		if tc.errorExpected {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}
