package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}
	Entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, Entry)

	require.Equal(t, arg.AccountID, Entry.AccountID)
	require.Equal(t, arg.Amount, Entry.Amount)

	require.NotZero(t, Entry.AccountID)
	require.NotZero(t, Entry.CreatedAt)

	return Entry
}
func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	Entry1 := createRandomEntry(t, account)
	Entry2, err := testQueries.GetEntry(context.Background(), Entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, Entry2)
	require.Equal(t, Entry1.ID, Entry2.ID)
	require.Equal(t, Entry1.AccountID, Entry2.AccountID)
	require.Equal(t, Entry1.Amount, Entry2.Amount)
	require.WithinDuration(t, Entry1.CreatedAt, Entry2.CreatedAt, time.Second)
}

func TestUpdateEntry(t *testing.T) {
	account := createRandomAccount(t)
	Entry1 := createRandomEntry(t, account)
	arg := UpdateEntryParams{
		ID:     Entry1.ID,
		Amount: util.RandomMoney(),
	}

	Entry2, err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, Entry2)

	require.Equal(t, Entry1.ID, Entry2.ID)
	require.Equal(t, Entry1.AccountID, Entry2.AccountID)
	require.Equal(t, arg.Amount, Entry2.Amount)
	require.WithinDuration(t, Entry1.CreatedAt, Entry2.CreatedAt, time.Second)
}

func TestDeleteEntry(t *testing.T) {
	account := createRandomAccount(t)
	Entry1 := createRandomEntry(t, account)
	err := testQueries.DeleteEntry(context.Background(), Entry1.ID)
	require.NoError(t, err)

	Entry2, err := testQueries.GetEntry(context.Background(), Entry1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, Entry2)
}

func TestListEntry(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, account)
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entries := range entries {
		require.NotEmpty(t, entries)
	}
}
