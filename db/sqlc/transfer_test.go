package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func CreateRandomTransfer() (CreateTransferParams, Transfer, error) {
	_, accountFrom, _ := CreateRandomAccount()
	_, accountTo, _ := CreateRandomAccount()

	arg := CreateTransferParams{
		FromAccountID: accountFrom.ID,
		ToAccountID:   accountTo.ID,
		Amount:        int64(20),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	return arg, transfer, err
}

func TestCreateTransfer(t *testing.T) {
	arg, transfer, err := CreateRandomTransfer()

	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
}

func TestGetTransfer(t *testing.T) {
	_, transferExp, _ := CreateRandomTransfer()

	transferAct, err := testQueries.GetTransfer(context.Background(), transferExp.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transferAct)
	require.Equal(t, transferExp.ID, transferAct.ID)
	require.Equal(t, transferExp.FromAccountID, transferAct.FromAccountID)
	require.Equal(t, transferExp.ToAccountID, transferAct.ToAccountID)
	require.Equal(t, transferExp.Amount, transferAct.Amount)
	require.WithinDuration(t, transferExp.CreatedAt, transferAct.CreatedAt, time.Second)
}

func TestListListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		_, _, _ = CreateRandomTransfer()
	}

	arg := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, entry := range transfers {
		require.NotEmpty(t, entry)
	}
}
