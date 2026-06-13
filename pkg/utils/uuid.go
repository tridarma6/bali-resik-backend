package utils

import (
	"github.com/google/uuid"
)

func NewUUID() uuid.UUID {
	return uuid.New()
}

func ParseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

func MustParseUUID(s string) uuid.UUID {
	return uuid.MustParse(s)
}
