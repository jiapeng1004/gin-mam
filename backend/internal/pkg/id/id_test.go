package id_test

import (
	"testing"

	"github.com/gin-mam/backend/internal/pkg/id"
	"github.com/stretchr/testify/require"
)

func TestNew_Returns32HexChars(t *testing.T) {
	got := id.New()
	require.Len(t, got, 32)
	require.NotContains(t, got, "-")
	require.Regexp(t, `^[0-9a-f]{32}$`, got)
}

func TestNew_IsUnique(t *testing.T) {
	a, b := id.New(), id.New()
	require.NotEqual(t, a, b)
}
