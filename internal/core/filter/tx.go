package filter

import (
	"context"

	"github.com/tonindexer/anton/addr"
	"github.com/tonindexer/anton/internal/core"
)

type TransactionsReq struct {
	Hash      []byte // `form:"hash"`
	InMsgHash []byte // `form:"in_msg_hash"`

	Addresses []*addr.Address //

	Workchain *int32 `form:"workchain"`

	BlockID *core.BlockID

	WithAccountState bool
	WithMessages     bool

	ExcludeColumn []string // TODO: support relations

	Order string `form:"order"` // ASC, DESC

	CreatedLT *uint64 `form:"created_lt"`

	AfterTxLT *uint64 `form:"after"`
	Limit     int     `form:"limit"`
}

type TraceReq struct {
	Hash                []byte // `form:"hash"`
	ExternalMessageHash []byte // `form:"ext_msg_hash"`
}

type TransactionsRes struct {
	Total int                 `json:"total"`
	Rows  []*core.Transaction `json:"results"`
}

type TraceRes struct {
	Total int                 `json:"total"`
	Root  *core.Transaction   `json:"root"`
	Rows  []*core.Transaction `json:"results"`
}

type TransactionRepository interface {
	FilterTransactions(context.Context, *TransactionsReq) (*TransactionsRes, error)
	FilterTrace(context.Context, *TraceReq) (*TraceRes, error)
}
