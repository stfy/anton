package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/tonindexer/anton/abi"
	"github.com/tonindexer/anton/internal/app"
	"github.com/tonindexer/anton/internal/core/filter"
	"github.com/xssnick/tonutils-go/ton"
	"net/http"
	"time"
)

var v1Path = "/api/v1"

type StormController struct {
	svc    app.QueryService
	parser app.ParserService
	API    ton.APIClientWrapped
}

func NewStormController(svc app.QueryService, parser app.ParserService, api ton.APIClientWrapped) *StormController {
	return &StormController{svc: svc, parser: parser, API: api}
}

func (c *StormController) GetPositionManagers(ctx *gin.Context) {
	ts := time.Now()

	var req filter.AccountsReq

	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		paramErr(ctx, "account_filter", err)
		return
	}

	ret, err := c.svc.FilterAccounts(ctx, &filter.AccountsReq{
		LatestState:   true,
		ContractTypes: []abi.ContractName{"position_manager_v2r1"},
		Columns: []string{
			"address",
			"data_hash",
			"updated_at",
			"last_tx_lt",
			"last_tx_hash",
			"shard",
			"workchain",
			"data",
			"block_seq_no",
		},
		Order:     "ASC",
		AfterTxLT: req.AfterTxLT,
	})
	if err != nil {
		internalErr(ctx, err)
		return
	}

	app.TimeTrack(ts, "load pm")

	ctx.IndentedJSON(http.StatusOK, ret)
}

func (c *StormController) GetNftItems(ctx *gin.Context) {
	var req filter.AccountsReq
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		paramErr(ctx, "account_filter", err)
		return
	}

	req.OwnerAddress, err = unmarshalAddress(ctx.Query("owner_address"))
	if err != nil {
		paramErr(ctx, "owner_address", err)
		return
	}
	req.MinterAddress, err = unmarshalAddress(ctx.Query("minter_address"))
	if err != nil {
		paramErr(ctx, "minter_address", err)
		return
	}

	if req.MinterAddress == nil {
		paramErr(ctx, "minter_address", errors.New("minter address should be provided"))
		return
	}

	if req.OwnerAddress == nil {
		paramErr(ctx, "owner_address", errors.New("owner address should be provided"))
		return
	}

	ret, err := c.svc.FilterNftAccounts(ctx, &filter.AccountsReq{
		LatestState:   true,
		ContractTypes: []abi.ContractName{"nft_item"},
		OwnerAddress:  req.OwnerAddress,
		MinterAddress: req.MinterAddress,
		Columns: []string{
			"address",
			"data_hash",
			"updated_at",
			"last_tx_lt",
			"last_tx_hash",
			"shard",
			"workchain",
			"minter_address",
			"owner_address",
			"data",
			"block_seq_no",
			"content_uri",
			"types",
			"executed_get_methods",
		},
	})
	if err != nil {
		internalErr(ctx, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, ret)
}
