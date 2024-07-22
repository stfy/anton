package notifier

import (
	"context"
	"encoding/json"
	"github.com/tmthrgd/go-hex"
	"github.com/tonindexer/anton/internal/app"
	"github.com/tonindexer/anton/internal/core"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/xssnick/tonutils-go/tlb"
)

var _ app.NotifierService = (*Kafka)(nil)

type KafkaConfig struct {
	Client *kgo.Client
}

type Kafka struct {
	*KafkaConfig
}

func (n *Kafka) Notify(ctx context.Context, entity any) error {
	//TODO implement me
	panic("implement me")
}

func (n *Kafka) NotifyAccounts(ctx context.Context, accs []*core.AccountState) error {
	for _, msg := range accs {
		msgValue, err := json.Marshal(msg)

		if err != nil {
			return err
		}

		p := n.Client.ProduceSync(
			ctx,
			&kgo.Record{
				Value: msgValue,
				Topic: "ACCOUNT",
			},
		)

		if err = p.FirstErr(); err != nil {
			return err
		}
	}

	return nil
}

func (n *Kafka) NotifyMessages(ctx context.Context, msgs []*core.Message, ext []*core.Message) error {
	records := make([]*kgo.Record, 0)

	for _, msg := range msgs {
		msgValue, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		records = append(records, &kgo.Record{
			Value: msgValue,
			Topic: "MESSAGE",
		})
	}

	for _, msg := range ext {
		data, hash, err := MessageToCell(msg.RawMessage)
		if err != nil {
			continue
		}

		records = append(records, &kgo.Record{
			Value: data,
			Topic: "EXT_MESSAGE",
			Key:   hash,
		})
	}

	p := n.Client.ProduceSync(ctx, records...)

	if err := p.FirstErr(); err != nil {
		return err
	}

	return nil
}

func NewKafkaNotifier(cfg *KafkaConfig) *Kafka {
	return &Kafka{KafkaConfig: cfg}
}

func MessageToCell(message tlb.AnyMessage) ([]byte, []byte, error) {
	c, err := tlb.ToCell(message)
	if err != nil {
		return nil, nil, err
	}

	msg := struct {
		RawData string `json:"rawData,omitempty"`
	}{}

	msgValue, err := json.Marshal(msg)
	if err != nil {
		return nil, nil, err
	}

	return msgValue, []byte(hex.EncodeToString(c.Hash())), nil
}
