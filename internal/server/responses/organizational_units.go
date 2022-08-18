package responses

import (
	"net/http"

	"github.com/suse-skyscraper/skyscraper/internal/db"
)

type OrganizationalUnitAttributes struct {
	ParentID    string `json:"parent_id"`
	DisplayName string `json:"display_name"`
}

type OrganizationalUnitItem struct {
	ID         string                       `json:"id"`
	Type       ObjectResponseType           `json:"type"`
	Attributes OrganizationalUnitAttributes `json:"attributes"`
}

func newOrganizationalUnitItem(organizationalUnit db.OrganizationalUnit) OrganizationalUnitItem {
	parentID := ""
	if organizationalUnit.ParentID.Valid {
		parentID = organizationalUnit.ParentID.UUID.String()
	}

	return OrganizationalUnitItem{
		ID:   organizationalUnit.ID.String(),
		Type: ObjectResponseTypeOrganizationalUnit,
		Attributes: OrganizationalUnitAttributes{
			ParentID:    parentID,
			DisplayName: organizationalUnit.DisplayName,
		},
	}
}

type OrganizationalUnitResponse struct {
	Data OrganizationalUnitItem `json:"data"`
}

func (rd *OrganizationalUnitResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewOrganizationalUnitResponse(organizationalUnit db.OrganizationalUnit) *OrganizationalUnitResponse {
	return &OrganizationalUnitResponse{
		Data: newOrganizationalUnitItem(organizationalUnit),
	}
}

type OrganizationalUnitsResponse struct {
	Data []OrganizationalUnitItem `json:"data"`
}

func (rd *OrganizationalUnitsResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewOrganizationalUnitsResponse(organizationalUnits []db.OrganizationalUnit) *OrganizationalUnitsResponse {
	list := make([]OrganizationalUnitItem, len(organizationalUnits))
	for i, ou := range organizationalUnits {
		list[i] = newOrganizationalUnitItem(ou)
	}

	return &OrganizationalUnitsResponse{
		Data: list,
	}
}
