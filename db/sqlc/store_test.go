package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	_, accountFrom, _ := CreateRandomAccount()
	_, accountTo, _ := CreateRandomAccount()

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: accountFrom.ID,
				ToAccountID:   accountTo.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, accountFrom.ID, transfer.FromAccountID)
		require.Equal(t, accountTo.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.NotZero(t, fromEntry.ID)
		require.Equal(t, accountFrom.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotEmpty(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.NotZero(t, toEntry.ID)
		require.Equal(t, accountTo.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotEmpty(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts' balance
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, accountFrom.ID, fromAccount.ID)
		require.Equal(t, accountFrom.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, accountTo.ID, toAccount.ID)
		require.Equal(t, accountTo.ID, toAccount.ID)

		// check account balance
		diffFromBalance := accountFrom.Balance - fromAccount.Balance
		diffToBalance := accountTo.Balance - toAccount.Balance
		require.Equal(t, diffFromBalance, -diffToBalance)
		require.True(t, diffFromBalance > 0)

		require.True(t, diffFromBalance%amount == 0)
		require.True(t, diffToBalance%amount == 0)

		k := int(diffFromBalance / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	transferAmount := int64(n) * amount

	updatedAccountFrom, err := store.GetAccount(context.Background(), accountFrom.ID)
	require.NoError(t, err)
	require.Equal(t, accountFrom.Balance-transferAmount, updatedAccountFrom.Balance)

	updatedAccountTo, err := store.GetAccount(context.Background(), accountTo.ID)
	require.NoError(t, err)
	require.Equal(t, accountTo.Balance+transferAmount, updatedAccountTo.Balance)

	fmt.Printf(">>> After Tx: from %v | to %v \n", updatedAccountFrom.Balance, updatedAccountTo.Balance)
}
