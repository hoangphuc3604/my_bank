package sqlc

import (
	"context"
	"testing"
	"time"

	"github.com/hoangphuc3064/MyBank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username: util.RandOwner(),
		Password: "password",
		Email:    util.RandEmail(),
		Fullname: util.RandOwner(),
	}
	
	user, err := testQueries.CreateUser(context.Background(), arg)
	if err != nil {
		t.Fatalf("cannot create user: %v", err)
	}

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Fullname, user.Fullname)
	require.Equal(t, arg.Password, user.Password)
	
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Fullname, user2.Fullname)
	require.Equal(t, user1.Password, user2.Password)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}