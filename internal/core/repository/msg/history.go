package msg

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/uptrace/go-clickhouse/ch"

	"github.com/tonindexer/anton/internal/core"
	"github.com/tonindexer/anton/internal/core/aggregate/history"
)

func addMessagesHistoryFilters(q *ch.SelectQuery, req *history.MessagesReq) *ch.SelectQuery {
	if len(req.SrcAddresses) > 0 {
		q = q.Where("src_address in (?)", ch.In(req.SrcAddresses))
	}
	if len(req.DstAddresses) > 0 {
		q = q.Where("dst_address in (?)", ch.In(req.DstAddresses))
	}
	if req.SrcWorkchain != nil {
		q = q.Where("src_workchain = ?", *req.SrcWorkchain)
	}
	if req.DstWorkchain != nil {
		q = q.Where("dst_workchain = ?", *req.DstWorkchain)
	}
	if len(req.SrcContracts) > 0 {
		q = q.Where("src_contract IN (?)", ch.In(req.SrcContracts))
	}
	if len(req.DstContracts) > 0 {
		q = q.Where("dst_contract IN (?)", ch.In(req.DstContracts))
	}
	if len(req.OperationNames) > 0 {
		q = q.Where("operation_name IN (?)", ch.In(req.OperationNames))
	}
	if req.MinterAddress != nil {
		q = q.Where("minter_address = ?", req.MinterAddress)
	}
	return q
}

func (r *Repository) AggregateMessagesHistory(ctx context.Context, req *history.MessagesReq) (*history.MessagesRes, error) {
	var res history.MessagesRes
	var bigIntRes bool // do we need to count account_data or account_states

	q := addMessagesHistoryFilters(r.ch.NewSelect().Model((*core.Message)(nil)), req)

	switch req.Metric {
	case history.MessageCount:
		q = q.ColumnExpr("count() as value")
	case history.MessageAmountSum:
		q, bigIntRes = q.ColumnExpr("sum(amount) as value"), true
	default:
		return nil, errors.Wrapf(core.ErrInvalidArg, "invalid message metric %s", req.Metric)
	}

	rounding, err := history.GetRoundingFunction(req.Interval)
	if err != nil {
		return nil, err
	}
	q = q.ColumnExpr(fmt.Sprintf(rounding, "created_at") + " as timestamp")
	q = q.Group("timestamp")

	if !req.From.IsZero() {
		q = q.Where("created_at > ?", req.From)
	}
	if !req.To.IsZero() {
		q = q.Where("created_at < ?", req.To)
	}

	q = q.Order("timestamp ASC")

	if bigIntRes {
		if err := q.Scan(ctx, &res.BigIntRes); err != nil {
			return nil, err
		}
	} else {
		if err := q.Scan(ctx, &res.CountRes); err != nil {
			return nil, err
		}
	}

	return &res, nil
}
