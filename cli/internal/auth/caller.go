package auth

import (
	"errors"

	"github.com/suse-skyscraper/skyscraper/cli/db"

	"github.com/google/uuid"
)

type CallerType int

const (
	_ CallerType = iota
	CallerUser
	CallerAPIKey
)

type Caller struct {
	ID   uuid.UUID
	Type CallerType
}

func (c *Caller) GetDBType() (db.CallerType, error) {
	switch c.Type {
	case CallerUser:
		return db.CallerTypeUser, nil
	case CallerAPIKey:
		return db.CallerTypeApiKey, nil
	}

	return "", errors.New("unknown caller type")
}
