package id

import (
	"strings"

	"github.com/google/uuid"
)

func New() string {
	return strings.ReplaceAll(uuid.Must(uuid.NewV7()).String(), "-", "")
}
