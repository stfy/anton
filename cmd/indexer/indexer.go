package indexer

import (
	"fmt"
	"github.com/allisson/go-env"
	"github.com/pkg/errors"
	"github.com/tonindexer/anton/internal/app/notifier"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/scram"
	"github.com/urfave/cli/v2"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/tonindexer/anton/abi"
	"github.com/tonindexer/anton/internal/app"
	"github.com/tonindexer/anton/internal/app/fetcher"
	"github.com/tonindexer/anton/internal/app/indexer"
	"github.com/tonindexer/anton/internal/app/parser"
	"github.com/tonindexer/anton/internal/core/repository"
	"github.com/tonindexer/anton/internal/core/repository/contract"
)

var Command = &cli.Command{
	Name:    "indexer",
	Aliases: []string{"idx"},
	Usage:   "Scans new blocks",

	Action: func(ctx *cli.Context) error {
		chURL := env.GetString("DB_CH_URL", "")
		pgURL := env.GetString("DB_PG_URL", "")

		brokerSeeds := env.GetStringSlice("BROKER_URL", ",", []string{""})
		kafkaOpts := []kgo.Opt{
			kgo.SeedBrokers(brokerSeeds...),
			kgo.AllowAutoTopicCreation(),
			kgo.ProducerBatchMaxBytes(int32(104857600)),
		}

		if env.GetBool("KAFKA_SASL_ENABLED", true) {
			kafkaOpts = append(
				kafkaOpts,
				kgo.SASL(
					scram.Auth{
						User: env.GetString("KAFKA_SASL_USERNAME", ""),
						Pass: env.GetString("KAFKA_SASL_PASSWORD", ""),
					}.AsSha256Mechanism(),
				),
			)
		}

		cl, err := kgo.NewClient(kafkaOpts...)
		if err = cl.Ping(ctx.Context); err != nil {
			return errors.Wrap(err, "initialize kafka error")
		}

		conn, err := repository.ConnectDB(ctx.Context, chURL, pgURL)
		if err != nil {
			return errors.Wrap(err, "cannot connect to a database")
		}

		contractRepo := contract.NewRepository(conn.PG)

		interfaces, err := contractRepo.GetInterfaces(ctx.Context)
		if err != nil {
			return errors.Wrap(err, "get interfaces")
		}
		if len(interfaces) == 0 {
			return errors.New("no contract interfaces")
		}

		def, err := contractRepo.GetDefinitions(ctx.Context)
		if err != nil {
			return errors.Wrap(err, "get definitions")
		}
		err = abi.RegisterDefinitions(def)
		if err != nil {
			return errors.Wrap(err, "get definitions")
		}

		client := liteclient.NewConnectionPool()
		api := ton.NewAPIClient(client, ton.ProofCheckPolicyUnsafe).WithRetry()
		for _, addr := range strings.Split(env.GetString("LITESERVERS", ""), ",") {
			split := strings.Split(addr, "|")
			if len(split) != 2 {
				return fmt.Errorf("wrong server address format '%s'", addr)
			}
			host, key := split[0], split[1]
			if err := client.AddConnection(ctx.Context, host, key); err != nil {
				return errors.Wrapf(err, "cannot add connection with %s host and %s key", host, key)
			}
		}
		bcConfig, err := app.GetBlockchainConfig(ctx.Context, api)
		if err != nil {
			return errors.Wrap(err, "cannot get blockchain config")
		}

		p := parser.NewService(&app.ParserConfig{
			BlockchainConfig: bcConfig,
			ContractRepo:     contractRepo,
		})
		f := fetcher.NewService(&app.FetcherConfig{
			API:    api,
			Parser: p,
		})
		i := indexer.NewService(&app.IndexerConfig{
			DB:        conn,
			API:       api,
			Parser:    p,
			Fetcher:   f,
			FromBlock: uint32(env.GetInt32("FROM_BLOCK", 1)),
			Workers:   env.GetInt("WORKERS", 4),
			Notifier:  notifier.NewKafkaNotifier(&notifier.KafkaConfig{Client: cl}),
		})
		if err = i.Start(); err != nil {
			return err
		}

		c := make(chan os.Signal, 1)
		done := make(chan struct{}, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-c
			i.Stop()
			conn.Close()
			done <- struct{}{}
		}()

		<-done

		return nil
	},
}
