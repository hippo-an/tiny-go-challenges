package db

import (
	"context"
	"testing"

	"github.com/hippo-an/tiny-go-challenges/mombank_11/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, fromA, toA Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: fromA.ID,
		ToAccountID:   toA.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	return transfer

}

func TestQueries_CreateTransfer(t *testing.T) {
	fromA := createRandomAccount(t)
	toA := createRandomAccount(t)
	createRandomTransfer(t, fromA, toA)
}

func TestQueries_GetTransfer(t *testing.T) {
	fromA := createRandomAccount(t)
	toA := createRandomAccount(t)
	transfer := createRandomTransfer(t, fromA, toA)

	trans, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.Equal(t, trans, transfer)
}

func TestQueries_ListTransfers(t *testing.T) {
	fromA := createRandomAccount(t)
	toA := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomTransfer(t, fromA, toA)
	}

	arg := ListTransfersParams{
		FromAccountID: fromA.ID,
		ToAccountID:   toA.ID,
		Limit:         5,
		Offset:        0,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, tr := range transfers {
		require.Equal(t, tr.FromAccountID, fromA.ID)
		require.Equal(t, tr.ToAccountID, toA.ID)
	}
}
