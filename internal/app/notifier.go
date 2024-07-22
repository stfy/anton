package app

import (
	"context"
	"github.com/tonindexer/anton/internal/core"
)

type NotifierService interface {
	Notify(ctx context.Context, entity any) error

	NotifyAccounts(ctx context.Context, msg []*core.AccountState) error
	NotifyMessages(ctx context.Context, msg []*core.Message, ext []*core.Message) error
}
