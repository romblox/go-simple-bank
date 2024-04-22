package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func CreateRandomEntry() (CreateEntryParams, Entry, error) {
	_, account, _ := CreateRandomAccount()
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    int64(15),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	return arg, entry, err
}

func TestCreateEntry(t *testing.T) {

	arg, entry, err := CreateRandomEntry()

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
}

func TestGetEntry(t *testing.T) {
	_, entryExp, err := CreateRandomEntry()
	entryAct, err := testQueries.GetEntry(context.Background(), entryExp.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entryAct)
	require.Equal(t, entryExp.ID, entryAct.ID)
	require.Equal(t, entryExp.AccountID, entryAct.AccountID)
	require.Equal(t, entryExp.Amount, entryAct.Amount)
	require.WithinDuration(t, entryExp.CreatedAt, entryAct.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		_, _, _ = CreateRandomEntry()
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
