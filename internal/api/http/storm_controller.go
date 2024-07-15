package http

import (
	"github.com/gin-gonic/gin"
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

// GetPositionManagers godoc
//
//	@Summary		statistics on all tables
//	@Description	Returns statistics on blocks, transactions, messages and accounts
//	@Tags			statistics
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}		aggregate.Statistics
//	@Router			/statistics [get]
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
