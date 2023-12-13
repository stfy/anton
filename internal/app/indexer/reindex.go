package indexer

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/tonindexer/anton/internal/app"
	"github.com/tonindexer/anton/internal/core"
	"github.com/tonindexer/anton/internal/core/repository/account"
	"github.com/tonindexer/anton/internal/core/repository/block"
	"github.com/tonindexer/anton/internal/core/repository/msg"
	"github.com/tonindexer/anton/internal/core/repository/tx"
	"sort"
	"time"
)

type Reindex struct {
	*Service
}

func NewReindexService(cfg *app.IndexerConfig) *Reindex {
	var s = new(Reindex)

	s.Service = NewService(cfg)
	s.IndexerConfig = cfg

	// validate config
	if s.Workers < 1 {
		s.Workers = 1
	}
	if s.FromBlock < 2 {
		s.FromBlock = 2
	}

	ch, pg := s.DB.CH, s.DB.PG
	s.txRepo = tx.NewRepository(ch, pg)
	s.msgRepo = msg.NewRepository(ch, pg)
	s.blockRepo = block.NewRepository(ch, pg)
	s.accountRepo = account.NewRepository(ch, pg)

	return s
}

func (s *Reindex) Start() error {
	ctx := context.Background()
	shard, err := s.API.LookupBlock(ctx, 0, 4000000000000000, 40600881)
	if err != nil {
		panic(err)
	}

	btx, err := s.Fetcher.BlockTransactions(ctx, shard)
	if err != nil {
		panic(err)
	}

	shardBlock := &core.Block{
		Workchain: shard.Workchain,
		Shard:     shard.Shard,
		SeqNo:     shard.SeqNo,
		RootHash:  shard.RootHash,
		FileHash:  shard.FileHash,
		MasterID: &core.BlockID{
			Workchain: 0,
			Shard:     0,
			SeqNo:     0,
		},
		Transactions: btx,
		ScannedAt:    time.Now(),
	}

	blk := core.Block{
		Workchain:         0,
		Shard:             0,
		SeqNo:             0,
		FileHash:          nil,
		RootHash:          nil,
		MasterID:          nil,
		Shards:            []*core.Block{shardBlock},
		TransactionsCount: 0,
		Transactions:      []*core.Transaction{},
		ScannedAt:         time.Time{},
	}

	msg, err := s.getBlockMessages(&blk)

	fmt.Println(msg)

	return err
}

func (s *Reindex) getParsedMessages(msg []*core.Message) []*core.Message {
	ctx := context.Background()

	parsedMessages := make([]*core.Message, 0)

	for _, message := range msg {
		err := s.Parser.ParseMessagePayload(ctx, message)
		if errors.Is(err, app.ErrImpossibleParsing) {
			continue
		}
		if err != nil {
			log.Error().Err(err).
				Hex("msg_hash", message.Hash).
				Hex("src_tx_hash", message.SrcTxHash).
				Str("src_addr", message.SrcAddress.String()).
				Hex("dst_tx_hash", message.DstTxHash).
				Str("dst_addr", message.DstAddress.String()).
				Uint32("op_id", message.OperationID).
				Msg("parse message payload")
		}

		parsedMessages = append(parsedMessages, message)
	}

	sort.Slice(parsedMessages, func(i, j int) bool {
		return parsedMessages[i].CreatedLT < parsedMessages[j].CreatedLT
	})

	return parsedMessages
}

func (s *Reindex) getBlockMessages(master *core.Block) ([]*core.Message, error) {
	newBlocks := append([]*core.Block{master}, master.Shards...)

	var newTransactions []*core.Transaction
	for i := range newBlocks {
		newTransactions = append(newTransactions, newBlocks[i].Transactions...)
	}

	parsed := s.getParsedMessages(s.uniqMessages(newTransactions))
	acc := make([]*core.AccountState, 0)
	for _, P := range parsed {
		acc = append(acc, P.DstState)
	}

	return parsed, nil
}
