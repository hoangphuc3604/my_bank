package sqlc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	arg := CreateAccountParams{
		Owner: "Phuc",
		Balance: 1000,
		Currency: "USD",
	}

	rs, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, rs)

	id, err := rs.LastInsertId()
	require.NoError(t, err)
	require.NotZero(t, id)

	account, err := testQueries.GetAccount(context.Background(), id)
	require.NoError(t, err)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
}