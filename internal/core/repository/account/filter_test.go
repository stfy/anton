package account

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/tonindexer/anton/abi"
	"github.com/tonindexer/anton/internal/addr"
	"github.com/tonindexer/anton/internal/core"
	"github.com/tonindexer/anton/internal/core/filter"
)

func TestFilterRepository(t *testing.T) {
	var (
		// filter by address
		address       *addr.Address
		addressStates []*core.AccountState

		// filter by latest
		latestState *core.AccountState

		// filter by contract type
		specialState *core.AccountState
	)

	initdb(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Run("drop tables", func(t *testing.T) {
		dropTables(t)
	})

	t.Run("create tables", func(t *testing.T) {
		createTables(t)
	})

	t.Run("insert test data", func(t *testing.T) {
		tx, err := pg.Begin()
		assert.Nil(t, err)

		for i := 0; i < 10; i++ { // 10 account states on 100 addresses
			var states []*core.AccountState

			for j := 0; j < 10; j++ {
				states = append(states, randAccountStates(10)...)
			}
			data := randAccountData(states)

			// filter by address
			address = &states[len(states)-10].Address
			addressStates = states[len(states)-10:]

			err = repo.AddAccountData(ctx, tx, data)
			assert.Nil(t, err)
			err = repo.AddAccountStates(ctx, tx, states)
			assert.Nil(t, err)
		}

		// filter by contract interfaces
		for i := 0; i < 15; i++ { // add 15 addresses with 10 states and "special" contract type
			states := randAccountStates(10)
			data := randContractData(states, "special")

			specialState = states[len(states)-1]
			specialState.StateData = data[len(data)-1]

			err = repo.AddAccountData(ctx, tx, data)
			assert.Nil(t, err)
			err = repo.AddAccountStates(ctx, tx, states)
			assert.Nil(t, err)
		}

		err = tx.Commit()
		assert.Nil(t, err)
	})

	t.Run("insert many states on some address", func(t *testing.T) {
		tx, err := pg.Begin()
		assert.Nil(t, err)

		for i := 0; i < 5; i++ { // 50 states on some address
			states := randAddressStates(address, 10)
			data := randAccountData(states)

			// filter by latest state
			latestState = states[len(states)-1]
			latestState.StateData = data[len(data)-1]

			err := repo.AddAccountData(ctx, tx, data)
			assert.Nil(t, err)
			err = repo.AddAccountStates(ctx, tx, states)
			assert.Nil(t, err)
		}

		err = tx.Commit()
		assert.Nil(t, err)
	})

	t.Run("filter states by address", func(t *testing.T) {
		results, err := repo.FilterAccounts(ctx, &filter.AccountsReq{
			Addresses: []*addr.Address{address},
			Order:     "ASC", Limit: len(addressStates),
		})
		assert.Nil(t, err)
		assert.Equal(t, 60, results.Total)
		assert.Equal(t, addressStates, results.Rows)
	})

	t.Run("filter latest state by address and exclude columns", func(t *testing.T) {
		latest := *latestState
		latest.StateData = nil
		latest.Code = nil

		results, err := repo.FilterAccounts(ctx, &filter.AccountsReq{
			Addresses:     []*addr.Address{&latest.Address},
			LatestState:   true,
			ExceptColumns: []string{"code"},
		})
		assert.Nil(t, err)
		assert.Equal(t, 1, results.Total)
		assert.Equal(t, []*core.AccountState{&latest}, results.Rows)
	})

	t.Run("filter latest state with data by address and exclude columns", func(t *testing.T) {
		latest := *latestState
		latest.Code = nil

		results, err := repo.FilterAccounts(ctx, &filter.AccountsReq{
			Addresses:   []*addr.Address{&latest.Address},
			LatestState: true, WithData: true,
			ExceptColumns: []string{"code"},
		})
		assert.Nil(t, err)
		assert.Equal(t, 1, results.Total)
		assert.Equal(t, []*core.AccountState{&latest}, results.Rows)
	})

	t.Run("filter latest state with data by contract types", func(t *testing.T) {
		results, err := repo.FilterAccounts(ctx, &filter.AccountsReq{
			ContractTypes: []abi.ContractName{"special"},
			LatestState:   true, WithData: true,
			Order: "DESC", Limit: 1,
		})
		assert.Nil(t, err)
		assert.Equal(t, 15, results.Total)
		assert.Equal(t, []*core.AccountState{specialState}, results.Rows)
	})

	t.Run("filter states by minter", func(t *testing.T) {
		results, err := repo.FilterAccounts(ctx, &filter.AccountsReq{
			WithData: true, MinterAddress: latestState.StateData.MinterAddress,
			Order: "DESC", Limit: 1,
		})
		assert.Nil(t, err)
		assert.Equal(t, 60, results.Total)
		assert.Equal(t, []*core.AccountState{latestState}, results.Rows)
	})

	t.Run("drop tables again", func(t *testing.T) {
		dropTables(t)
	})
}
