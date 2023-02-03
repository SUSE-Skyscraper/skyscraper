package responses

import "net/http"

type OrganizationalUnitAttributes struct {
	ParentID    string `json:"parent_id"`
	DisplayName string `json:"display_name"`
}

type OrganizationalUnitItem struct {
	ID         string                       `json:"id"`
	Type       ObjectResponseType           `json:"type"`
	Attributes OrganizationalUnitAttributes `json:"attributes"`
}

type OrganizationalUnitResponse struct {
	Data OrganizationalUnitItem `json:"data"`
}

func (rd *OrganizationalUnitResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

type OrganizationalUnitsResponse struct {
	Data []OrganizationalUnitItem `json:"data"`
}

func (rd *OrganizationalUnitsResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
