package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func createRandomTransfer(t *testing.T, account1, account2 Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}
	Transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, Transfer)

	require.Equal(t, arg.FromAccountID, Transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, Transfer.ToAccountID)
	require.Equal(t, arg.Amount, Transfer.Amount)

	require.NotZero(t, Transfer.ID)
	require.NotZero(t, Transfer.CreatedAt)

	return Transfer
}
func TestCreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	createRandomTransfer(t, account1, account2)
}

func TestGetTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	Transfer1 := createRandomTransfer(t, account1, account2)

	Transfer2, err := testQueries.GetTransfer(context.Background(), Transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, Transfer2)
	require.Equal(t, Transfer1.ID, Transfer2.ID)

	require.Equal(t, Transfer1.ToAccountID, Transfer2.ToAccountID)
	require.Equal(t, Transfer1.FromAccountID, Transfer2.FromAccountID)
	require.Equal(t, Transfer1.Amount, Transfer2.Amount)
	require.WithinDuration(t, Transfer1.CreatedAt, Transfer2.CreatedAt, time.Second)
}

func TestUpdateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	Transfer1 := createRandomTransfer(t, account1, account2)

	arg := UpdateTransfersParams{
		ID:     Transfer1.ID,
		Amount: util.RandomMoney(),
	}

	Transfer2, err := testQueries.UpdateTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, Transfer2)

	require.Equal(t, Transfer1.ID, Transfer2.ID)
	require.Equal(t, Transfer1.ToAccountID, Transfer2.ToAccountID)
	require.Equal(t, Transfer1.FromAccountID, Transfer2.FromAccountID)
	require.Equal(t, arg.Amount, Transfer2.Amount)
	require.WithinDuration(t, Transfer1.CreatedAt, Transfer1.CreatedAt, time.Second)
}

func TestDeleteTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	Transfer1 := createRandomTransfer(t, account1, account2)
	err := testQueries.DeleteTransfer(context.Background(), Transfer1.ID)
	require.NoError(t, err)

	Transfer2, err := testQueries.GetTransfer(context.Background(), Transfer1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, Transfer2)
}

func TestListTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, account1, account2)
		createRandomTransfer(t, account2, account1)
	}

	arg := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account1.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == account1.ID || transfer.ToAccountID == account1.ID)
	}
}
