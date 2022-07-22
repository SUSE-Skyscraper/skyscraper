package auth

import (
	"errors"

	"github.com/google/uuid"
	"github.com/suse-skyscraper/skyscraper/internal/db"
)

type CallerType int

const (
	CallerUser CallerType = iota
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
