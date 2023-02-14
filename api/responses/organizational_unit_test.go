package responses

import "testing"

func TestOrganizationalUnitResponse_Render(t *testing.T) {
	resp := OrganizationalUnitResponse{}
	err := resp.Render(nil, nil)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}

func TestOrganizationalUnitsResponse_Render(t *testing.T) {
	resp := OrganizationalUnitsResponse{}
	err := resp.Render(nil, nil)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}
