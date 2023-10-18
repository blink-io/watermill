package watermill

import (
	"github.com/google/uuid"
)

// NewUUID returns a new UUID Version 4.
func NewUUID() string {
	return uuid.New().String()
}
