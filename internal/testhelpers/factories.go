package testhelpers

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/suse-skyscraper/skyscraper/internal/db"
)

func FactoryOrganizationalUnit() db.OrganizationalUnit {
	return db.OrganizationalUnit{
		ID: uuid.New(),
		ParentID: uuid.NullUUID{
			UUID:  uuid.New(),
			Valid: true,
		},
		DisplayName: "test",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func FactoryCloudAccount() db.CloudAccount {
	return db.CloudAccount{
		Cloud:             "aws",
		TenantID:          "1234",
		AccountID:         "12345",
		Name:              "test",
		Active:            true,
		TagsCurrent:       pgtype.JSONB{Bytes: []byte("{}"), Status: pgtype.Present},
		TagsDesired:       pgtype.JSONB{Bytes: []byte("{}"), Status: pgtype.Present},
		TagsDriftDetected: false,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}
