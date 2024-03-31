package web

import (
	"fmt"
	"github.com/redis/rueidis"
	"github.com/tonindexer/anton/internal/app/parser"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/allisson/go-env"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"

	"github.com/tonindexer/anton/abi"
	"github.com/tonindexer/anton/internal/api/http"
	"github.com/tonindexer/anton/internal/app"
	"github.com/tonindexer/anton/internal/app/query"
	"github.com/tonindexer/anton/internal/core/repository"
	"github.com/tonindexer/anton/internal/core/repository/contract"
)

var Command = &cli.Command{
	Name:  "web",
	Usage: "HTTP JSON API",

	Action: func(ctx *cli.Context) error {
		chURL := env.GetString("DB_CH_URL", "")
		pgURL := env.GetString("DB_PG_URL", "")
		rsURL := env.GetString("REDIS_URL", "")

		rsClient, err := rueidis.NewClient(rueidis.ClientOption{InitAddress: []string{rsURL}})
		if err != nil {
			return errors.Wrap(err, "cannot connect redis client")
		}
		defer rsClient.Close()

		conn, err := repository.ConnectDB(
			ctx.Context, chURL, pgURL)
		if err != nil {
			return errors.Wrap(err, "cannot connect to a database")
		}

		repository.WithCache(conn, rsClient)
		contractRepo := contract.NewRepository(conn.PG)

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
		qs, err := query.NewService(ctx.Context, &app.QueryConfig{
			DB:  conn,
			API: api,
		})
		if err != nil {
			return err
		}

		srv := http.NewServer(
			env.GetString("LISTEN", "0.0.0.0:80"),
		)
		srv.RegisterRoutes(http.NewController(qs, p, api))

		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-c
			conn.Close()
			os.Exit(0)
		}()

		if err = srv.Run(); err != nil {
			return err
		}

		return nil
	},
}
