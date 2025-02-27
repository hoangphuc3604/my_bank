package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandString(6)

	hashed, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashed)

	err = CheckPassword(password, hashed)
	require.NoError(t, err)

	wrongPassword := RandString(6)
	err = CheckPassword(wrongPassword, hashed)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

}