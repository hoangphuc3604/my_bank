package token

import (
	"testing"
	"time"

	"github.com/hoangphuc3064/MyBank/util"
	"github.com/stretchr/testify/require"
)

func TestPaseto(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandString(32))
	require.NoError(t, err)

	username := util.RandString(13)
	role := util.RandString(10)
	duration := time.Minute

	issuedAt := time.Now()
	expiratedAt := time.Now().Add(duration)

	token, payload, err := maker.CreateToken(username, role, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)
	require.Equal(t, username, payload.Username)
	require.Equal(t, role, payload.Role)
	require.NotZero(t, payload.ID)
	require.NotZero(t, payload.IssuedAt)
	require.NotZero(t, payload.ExpiredAt)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
	require.WithinDuration(t, payload.ExpiredAt, expiratedAt, time.Second)

	payload2, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload2)

	require.NotEmpty(t, payload2.ID)
	require.Equal(t, payload.ID, payload2.ID)
	require.Equal(t, payload.Username, payload2.Username)
	require.Equal(t, payload.Role, payload2.Role)
	require.Equal(t, payload.IssuedAt, payload2.IssuedAt)
	require.Equal(t, payload.ExpiredAt, payload2.ExpiredAt)
	require.WithinDuration(t, payload2.IssuedAt, issuedAt, time.Second)
	require.WithinDuration(t, payload2.ExpiredAt, expiratedAt, time.Second)
}