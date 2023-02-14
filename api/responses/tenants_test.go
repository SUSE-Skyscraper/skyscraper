package responses

import "testing"

func TestCloudTenantResponse_Render(t *testing.T) {
	resp := CloudTenantResponse{}
	err := resp.Render(nil, nil)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}

func TestCloudTenantsResponse_Render(t *testing.T) {
	resp := CloudTenantsResponse{}
	err := resp.Render(nil, nil)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}
