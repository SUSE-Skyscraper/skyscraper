package scimbridgedb

import (
	"encoding/json"

	"github.com/jackc/pgtype"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/database"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/payloads"
	"github.com/suse-skyscraper/skyscraper/internal/db"
)

func toScimGroup(group db.Group) database.Group {
	scimGroup := database.Group{
		ID:          group.ID,
		DisplayName: group.DisplayName,
		CreatedAt:   group.CreatedAt,
		UpdatedAt:   group.UpdatedAt,
	}

	return scimGroup
}

func toScimUser(user db.User) (database.User, error) {
	var name map[string]string
	if user.Name.Bytes != nil {
		err := json.Unmarshal(user.Name.Bytes, &name)
		if err != nil {
			return database.User{}, err
		}
	}

	var emails []payloads.UserEmail
	if user.Emails.Bytes != nil {
		err := json.Unmarshal(user.Emails.Bytes, &emails)
		if err != nil {
			return database.User{}, err
		}
	}

	scimUser := database.User{
		ID:          user.ID,
		Username:    user.Username,
		ExternalID:  user.ExternalID,
		Name:        name,
		DisplayName: user.DisplayName,
		Locale:      user.Locale,
		Active:      user.Active,
		Emails:      emails,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	return scimUser, nil
}

func parseJSONB(arg interface{}) (pgtype.JSONB, error) {
	if arg == nil {
		return pgtype.JSONB{Bytes: nil, Status: pgtype.Null}, nil
	}

	name := pgtype.JSONB{}
	err := name.Set(arg)
	if err != nil {
		return pgtype.JSONB{}, err
	}

	return name, nil
}
