package app

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/xssnick/tonutils-go/tl"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/tvm/cell"

	"github.com/tonindexer/anton/addr"
	"github.com/tonindexer/anton/internal/core"
)

var ErrImpossibleParsing = errors.New("parsing is impossible")

type ParserConfig struct {
	BlockchainConfig *cell.Cell
	ContractRepo     core.ContractRepository
}

func GetBlockchainConfig(ctx context.Context, api ton.APIClientWrapped) (*cell.Cell, error) {
	var res tl.Serializable

	b, err := api.GetMasterchainInfo(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get masterchain info")
	}

	err = api.Client().QueryLiteserver(ctx, ton.GetConfigAll{Mode: 0, BlockID: b}, &res)
	if err != nil {
		return nil, err
	}

	switch t := res.(type) {
	case ton.ConfigAll:
		var state tlb.ShardStateUnsplit

		ref, err := t.ConfigProof.BeginParse().LoadRef()
		if err != nil {
			return nil, err
		}

		err = tlb.LoadFromCell(&state, ref)
		if err != nil {
			return nil, err
		}

		var mcStateExtra tlb.McStateExtra
		err = tlb.LoadFromCell(&mcStateExtra, state.McStateExtra.BeginParse())
		if err != nil {
			return nil, err
		}

		return tlb.ToCell(mcStateExtra.ConfigParams.Config.Params)
		//return mcStateExtra.ConfigParams.Config.Params.ToCell()
	case ton.LSError:
		return nil, t

	default:
		return nil, fmt.Errorf("got unknown response")
	}
}

type ParserService interface {
	ParseAccountData(
		ctx context.Context,
		acc *core.AccountState,
		others func(context.Context, addr.Address) (*core.AccountState, error),
	) error

	ParseMessagePayload(
		ctx context.Context,
		message *core.Message, // source and destination account states must be known
	) error
}
