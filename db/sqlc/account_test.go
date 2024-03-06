package db

import (
	"context"
	"database/sql"
	"github.com/rusrom/simplebank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func CreateRandomAccount() (CreateAccountParams, Account, error) {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	return arg, account, err
}

func TestCreateAccount(t *testing.T) {
	arg, account, err := CreateRandomAccount()

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}

func TestGetAccount(t *testing.T) {
	_, accountExp, err := CreateRandomAccount()
	accountAct, err := testQueries.GetAccount(context.Background(), accountExp.ID)

	require.NoError(t, err)
	require.NotEmpty(t, accountAct)
	require.Equal(t, accountExp.ID, accountAct.ID)
	require.Equal(t, accountExp.Owner, accountAct.Owner)
	require.Equal(t, accountExp.Balance, accountAct.Balance)
	require.Equal(t, accountExp.Currency, accountAct.Currency)
	require.WithinDuration(t, accountExp.CreatedAt, accountAct.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	_, accountAct, _ := CreateRandomAccount()

	arg := UpdateAccountParams{
		ID:      accountAct.ID,
		Balance: util.RandomBalance(),
	}

	accountExp, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, accountExp)
	require.Equal(t, accountExp.ID, accountAct.ID)
	require.Equal(t, accountExp.Owner, accountAct.Owner)
	require.Equal(t, accountExp.Balance, arg.Balance)
	require.Equal(t, accountExp.Currency, accountAct.Currency)
	require.WithinDuration(t, accountExp.CreatedAt, accountAct.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	_, accountAct, _ := CreateRandomAccount()

	err := testQueries.DeleteAccount(context.Background(), accountAct.ID)

	require.NoError(t, err)

	account, err := testQueries.GetAccount(context.Background(), accountAct.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		_, _, _ = CreateRandomAccount()
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
