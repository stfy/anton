package notifier

import (
	"context"
	"encoding/json"
	"github.com/tonindexer/anton/internal/app"
	"github.com/tonindexer/anton/internal/core"
	"github.com/twmb/franz-go/pkg/kgo"
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

func (n *Kafka) NotifyMessages(ctx context.Context, msgs []*core.Message) error {
	for _, msg := range msgs {
		msgValue, err := json.Marshal(msg)

		if err != nil {
			return err
		}

		p := n.Client.ProduceSync(
			ctx,
			&kgo.Record{
				Value: msgValue,
				Topic: "MESSAGE",
			},
		)

		if err = p.FirstErr(); err != nil {
			return err
		}
	}

	return nil
}

func NewKafkaNotifier(cfg *KafkaConfig) *Kafka {
	return &Kafka{KafkaConfig: cfg}
}
