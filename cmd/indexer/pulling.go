package main

import (
	"context"

	"github.com/721tools/backend-go/indexer/internal/boot"
	"github.com/721tools/backend-go/indexer/internal/bus"
	"github.com/721tools/backend-go/indexer/pkg/blockchain/client"
	"github.com/721tools/backend-go/indexer/pkg/db"
	"github.com/721tools/backend-go/indexer/pkg/mq"
	"github.com/721tools/backend-go/indexer/pkg/utils/flags"
	"github.com/spf13/cobra"
)

func newPullingCmd() *cobra.Command {
	type PullingFlag struct {
		CacheSize uint8
		PullStep  uint8
		DSN       string
		RPC       string
		RedisURL  string
	}
	var flag PullingFlag
	pullingCmd := &cobra.Command{
		Use:   "pulling",
		Short: "continue pulling blockchain",
		Long:  `continue pulling data from blockchain`,
		Run: func(cmd *cobra.Command, args []string) {
			flags.PrintFlags(cmd.Flags())
			db.NewDBEngine(flag.DSN)
			mq.NewMQ(context.Background(), flag.RedisURL)
			client.NewClient(flag.RPC, client.EvmClient)
			bus.Init()
			boot.Continuous(flag.CacheSize, flag.PullStep)
		},
	}
	pullingCmd.Flags().Uint8Var(&flag.CacheSize, "cache_size", 3, "ex: 1 2 3")
	pullingCmd.Flags().Uint8Var(&flag.PullStep, "pull_step", 3, "ex: 30 50 80")
	pullingCmd.Flags().StringVar(&flag.RedisURL, "redis_url", "redis://<user>:<password>@<host>:<port>/<db_number>", "redis url")
	pullingCmd.Flags().StringVar(&flag.DSN, "dsn", "${USER}:${PASSWORD}@tcp(${HOST}:${HOST_PORT})/${DB_NAME}?parseTime=true", "db dsn")
	pullingCmd.Flags().StringVar(&flag.RPC, "rpc", "blockchain rpc", "block chain rpc")
	return pullingCmd
}

func init() {
	rootCmd.AddCommand(newPullingCmd())
}
