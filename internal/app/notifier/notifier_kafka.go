package notifier

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tonindexer/anton/internal/app"
	"github.com/tonindexer/anton/internal/core"
	"github.com/twmb/franz-go/pkg/kgo"
	"golang.org/x/exp/slices"
)

var skipNftAddresses = []string{
	"EQCcuodUc7NuhiifxAaEvVf8Yu2C3xRHhhFrsPqKfTfHI4lO",
}

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
	for id, msg := range msgs {
		if slices.Contains(skipNftAddresses, msg.DstState.Address.Base64()) {
			continue
		}

		msgValue, err := json.Marshal(msg)

		if err != nil {
			return err
		}

		fmt.Println(id, "bytes", len(msgValue))

		fmt.Println(id, "bytes", string(msgValue))

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
