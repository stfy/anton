package latest

import (
	"context"
	"encoding/json"
	"github.com/redis/rueidis"
	"github.com/tonindexer/anton/abi"
	"github.com/tonindexer/anton/internal/app"
	"github.com/tonindexer/anton/internal/core"
	"time"
)

type Cache struct {
	client rueidis.Client
}

func AddAccount(ctx context.Context, client rueidis.Client, acc *core.AccountState) error {
	for _, t := range acc.Types {
		v, err := json.Marshal(&acc)
		if err != nil {
			return err
		}

		cmd :=
			client.
				B().
				Hset().
				Key(string(t)).
				FieldValue().
				FieldValue(acc.Address.String(), string(v)).
				Build()

		res := client.Do(ctx, cmd)
		if err := res.Error(); err != nil {
			return err
		}
	}

	return nil
}

func (srv *Cache) AddAccount(ctx context.Context, acc core.AccountState) error {
	for _, t := range acc.Types {
		v, err := json.Marshal(&acc)
		if err != nil {
			return err
		}

		cmd := srv.
			client.
			B().
			Hset().
			Key(string(t)).
			FieldValue().
			FieldValue(acc.Address.String(), string(v)).
			Build()

		res := srv.client.Do(ctx, cmd)
		if err := res.Error(); err != nil {
			return err
		}
	}

	return nil
}

func GetLatestAccounts(ctx context.Context, client rueidis.Client, value abi.ContractName) ([]*core.AccountState, error) {
	defer app.TimeTrack(time.Now(), "cache.GetLatestAccounts")

	result := make([]*core.AccountState, 0)

	cmd := client.B().Hgetall().Key(string(value)).Build()

	res := client.Do(ctx, cmd)
	if err := res.Error(); err != nil {
		return nil, err
	}

	accounts, err := res.AsMap()
	if err != nil {
		return nil, err
	}
	for key := range accounts {
		var acc core.AccountState

		message := accounts[key]
		err := message.DecodeJSON(&acc)
		if err != nil {
			return nil, err
		}

		result = append(result, &acc)
	}

	return result, nil
}
