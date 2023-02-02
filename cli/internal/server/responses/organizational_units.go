package responses

import (
	resp "github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/internal/db"
)

func newOrganizationalUnitItem(organizationalUnit db.OrganizationalUnit) resp.OrganizationalUnitItem {
	parentID := ""
	if organizationalUnit.ParentID.Valid {
		parentID = organizationalUnit.ParentID.UUID.String()
	}

	return resp.OrganizationalUnitItem{
		ID:   organizationalUnit.ID.String(),
		Type: resp.ObjectResponseTypeOrganizationalUnit,
		Attributes: resp.OrganizationalUnitAttributes{
			ParentID:    parentID,
			DisplayName: organizationalUnit.DisplayName,
		},
	}
}

func NewOrganizationalUnitResponse(organizationalUnit db.OrganizationalUnit) *resp.OrganizationalUnitResponse {
	return &resp.OrganizationalUnitResponse{
		Data: newOrganizationalUnitItem(organizationalUnit),
	}
}

func NewOrganizationalUnitsResponse(organizationalUnits []db.OrganizationalUnit) *resp.OrganizationalUnitsResponse {
	list := make([]resp.OrganizationalUnitItem, len(organizationalUnits))
	for i, ou := range organizationalUnits {
		list[i] = newOrganizationalUnitItem(ou)
	}

	return &resp.OrganizationalUnitsResponse{
		Data: list,
	}
}
