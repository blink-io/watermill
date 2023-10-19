package watermill

import (
	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
)

// NewUUID returns a new UUID Version 4.
func NewUUID() string {
	return uuid.New().String()
}

// NewShortUUID returns a new short UUID.
func NewShortUUID() string {
	return shortuuid.New()
}
