package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthHash(t *testing.T) {
	password := "helloPass"
	hashOne, err := HashPassword(password)
	require.NoError(t, err)
	require.NotNil(t, hashOne)

	isValid, err := ValidatePassword(password, hashOne)
	require.NoError(t, err)
	assert.True(t, isValid)
}
