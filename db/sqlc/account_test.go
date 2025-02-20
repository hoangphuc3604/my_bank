package sqlc

import (
	"context"
	"database/sql"
	"testing"

	"github.com/hoangphuc3064/MyBank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandOwner(),
		Balance:  util.RandBalance(1, 10000),
		Currency: util.RandCurrency(),
	}

	rs, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, rs)

	id, err := rs.LastInsertId()
	require.NoError(t, err)
	require.NotZero(t, id)

	return getAccountByID(t, id, arg)
}

func getAccountByID(t *testing.T, id int64, arg CreateAccountParams) Account {
	account, err := testQueries.GetAccount(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.Equal(t, account1, account2)

	require.NotEmpty(t, account2.Owner)
	require.NotEmpty(t, account2.Currency)
	require.NotZero(t, account2.Balance)
	require.NotZero(t, account2.ID)

	require.NotEmpty(t, account2.CreatedAt)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandBalance(0, 10000),
	}

	err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.CreatedAt, account2.CreatedAt)

	require.NotEmpty(t, account2.CreatedAt)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}
