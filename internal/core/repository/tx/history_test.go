package tx_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/tonindexer/anton/addr"
	"github.com/tonindexer/anton/internal/core/aggregate/history"
	"github.com/tonindexer/anton/internal/core/rndm"
)

func TestRepository_AggregateTransactionsHistory(t *testing.T) {
	initdb(t)

	transactions := rndm.BlockTransactions(rndm.BlockID(-1), 10)

	a := rndm.Address()
	addrTransactions := rndm.AddressTransactions(a, 20)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	t.Run("drop tables", func(t *testing.T) {
		dropTables(t)
	})

	t.Run("create tables", func(t *testing.T) {
		createTables(t)
	})

	t.Run("add transactions", func(t *testing.T) {
		dbtx, err := pg.Begin()
		require.Nil(t, err)

		err = repo.AddTransactions(ctx, dbtx, transactions)
		require.Nil(t, err)

		err = repo.AddTransactions(ctx, dbtx, addrTransactions)
		require.Nil(t, err)

		err = dbtx.Commit()
		require.Nil(t, err)
	})

	t.Run("transaction count", func(t *testing.T) {
		res, err := repo.AggregateTransactionsHistory(ctx, &history.TransactionsReq{
			Metric:    history.TransactionCount,
			ReqParams: history.ReqParams{Interval: 4 * time.Hour},
		})
		require.Nil(t, err)
		require.Equal(t, 1, len(res.CountRes))
		require.Equal(t, 30, res.CountRes[0].Value)
	})

	t.Run("transaction count by workchain", func(t *testing.T) {
		res, err := repo.AggregateTransactionsHistory(ctx, &history.TransactionsReq{
			Metric:    history.TransactionCount,
			Workchain: new(int32),
			ReqParams: history.ReqParams{Interval: 4 * time.Hour},
		})
		require.Nil(t, err)
		require.Equal(t, 1, len(res.CountRes))
		require.Equal(t, 20, res.CountRes[0].Value)
	})

	t.Run("transaction count by workchain", func(t *testing.T) {
		res, err := repo.AggregateTransactionsHistory(ctx, &history.TransactionsReq{
			Metric:    history.TransactionCount,
			Addresses: []*addr.Address{a},
			ReqParams: history.ReqParams{Interval: 4 * time.Hour},
		})
		require.Nil(t, err)
		require.Equal(t, 1, len(res.CountRes))
		require.Equal(t, 20, res.CountRes[0].Value)
	})

	t.Run("drop tables again", func(t *testing.T) {
		dropTables(t)
	})
}
