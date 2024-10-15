package latest

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/rueidis"
	"github.com/tonindexer/anton/abi"
	"github.com/tonindexer/anton/addr"
	"github.com/tonindexer/anton/internal/app"
	"github.com/tonindexer/anton/internal/core"
	"slices"
	"strconv"
	"time"
)

type Cache struct {
	client rueidis.Client
}

func nftCollectionKey(acc *addr.Address) string {
	return fmt.Sprintf("nft_collection/%s", acc.String())
}

func nftOwnerKey(acc *addr.Address) string {
	return fmt.Sprintf("nft_collection_owners/%s", acc.String())
}

func getNftItemId(state *core.AccountState) (int64, error) {
	idx := slices.IndexFunc(state.ExecutedGetMethods["nft_item"], func(execution abi.GetMethodExecution) bool {
		return execution.Name == "get_nft_data"
	})

	switch v := state.ExecutedGetMethods["nft_item"][idx].Returns[1].(type) {
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(Base64ToHex(v), 16, 64)
	default:
		return -1, errors.New("unreqcoginzed interafce type")
	}
}

func ClearCacheNftItems(ctx context.Context, client rueidis.Client, address *addr.Address, owner *addr.Address) error {
	cmd := client.
		B().
		Del().
		Key(fmt.Sprintf("%s/%s", nftCollectionKey(address), owner)).Build()

	clearCached := client.B().Del().Key(fmt.Sprintf("cached_nft_collection/%s/%s", address, owner)).Build()

	for _, resp := range client.DoMulti(ctx, cmd, clearCached) {
		if err := resp.Error(); err != nil {
			return err
		}
	}

	return nil
}

func AddAccount(ctx context.Context, client rueidis.Client, acc *core.AccountState) error {
	for _, t := range acc.Types {
		switch t {
		case "nft_item":
			v, err := json.Marshal(&acc)
			if err != nil {
				return err
			}

			if acc.OwnerAddress == nil {
				continue
			}

			cmd := client.
				B().
				Hset().
				Key(fmt.Sprintf("%s/%s", nftCollectionKey(acc.MinterAddress), acc.OwnerAddress)).
				FieldValue().
				FieldValue(acc.Address.String(), string(v)).
				Build()

			res := client.Do(ctx, cmd)
			if err := res.Error(); err != nil {
				return err
			}
		}
	}

	return nil
}

func SetNftCollectionAsCached(ctx context.Context, client rueidis.Client, address *addr.Address, owner *addr.Address) (bool, error) {
	defer app.TimeTrack(time.Now(), "cache.GetLatestAccounts")

	cmd := client.B().Setex().Key(fmt.Sprintf("cached_nft_collection/%s/%s", address, owner)).Seconds(60).Value(address.String()).Build()
	res := client.Do(ctx, cmd)

	if err := res.Error(); err != nil {
		return false, err
	}

	return true, nil
}

func GetNftCollectionCached(ctx context.Context, client rueidis.Client, address *addr.Address, owner *addr.Address) (bool, error) {
	defer app.TimeTrack(time.Now(), "cache.GetLatestAccounts")

	cmd := client.B().Ttl().Key(fmt.Sprintf("cached_nft_collection/%s/%s", address, owner)).Build()
	res := client.Do(ctx, cmd)

	ttl, err := res.AsInt64()
	if err != nil {
		return false, err
	}

	switch {
	case ttl <= 0:
		return false, nil
	default:
		return true, nil
	}
}

func GetNftCollectionItems(ctx context.Context, client rueidis.Client, address *addr.Address, owner *addr.Address) ([]*core.AccountState, error) {
	defer app.TimeTrack(time.Now(), "cache.GetLatestAccounts")

	result := make([]*core.AccountState, 0)
	cmd := client.B().Hgetall().Key(fmt.Sprintf("%s/%s", nftCollectionKey(address), owner)).Build()

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

func GetLatestAccount(ctx context.Context, client rueidis.Client, accountType, address string) (*core.AccountState, error) {
	defer app.TimeTrack(time.Now(), "cache.GetLatestAccounts")
	cmd := client.B().Hget().Key(accountType).Field(address).Build()

	res := client.Do(ctx, cmd)
	if err := res.Error(); err != nil {
		return nil, err
	}

	var acc core.AccountState
	err := res.DecodeJSON(&acc)
	if err != nil {
		return nil, err
	}

	return &acc, nil
}

func Base64ToHex(str string) string {
	v, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(v)
}
