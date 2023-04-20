package implement

import (
	"context"

	"github.com/721tools/backend-go/indexer/implement/boot"
	"github.com/721tools/backend-go/indexer/implement/bus"
	"github.com/721tools/backend-go/indexer/pkg/blockchain/client"
	"github.com/721tools/backend-go/indexer/pkg/db"
	"github.com/721tools/backend-go/indexer/pkg/mq"
	"github.com/721tools/backend-go/indexer/pkg/utils/flags"
	"github.com/spf13/cobra"
)

func newFixCmd() *cobra.Command {
	type FixFlag struct {
		CacheSize uint8
		PullSize  uint8
		DSN       string
		RPC       string
		RedisURL  string
		Start     uint64
		End       uint64
	}
	var flag FixFlag
	fixCmd := &cobra.Command{
		Use:   "fix",
		Short: "fix blockchain",
		Long:  `fix data from blockchain`,
		Run: func(cmd *cobra.Command, args []string) {
			flags.PrintFlags(cmd.Flags())
			db.NewDBEngine(flag.DSN)
			mq.NewMQ(context.Background(), flag.RedisURL)
			client.NewClient(flag.RPC, client.EvmClient)
			bus.Init()
			boot.Temporarily(flag.Start, flag.End, flag.CacheSize, flag.PullSize)
		},
	}
	//
	fixCmd.Flags().Uint8Var(&flag.CacheSize, "cache_size", 3, "ex: 1 2 3")
	fixCmd.Flags().Uint8Var(&flag.PullSize, "pull_step", 30, "ex: 30 50 80")
	fixCmd.Flags().StringVar(&flag.RedisURL, "redis_url", "redis://<user>:<password>@<host>:<port>/<db_number>", "redis url")
	fixCmd.Flags().StringVar(&flag.DSN, "dsn", "${USER}:${PASSWORD}@tcp(${HOST}:${HOST_PORT})/${DB_NAME}?parseTime=true", "db dsn")
	fixCmd.Flags().StringVar(&flag.RPC, "rpc", "", "rpc node")
	fixCmd.Flags().Uint64Var(&flag.Start, "start", 1, "start point")
	fixCmd.Flags().Uint64Var(&flag.End, "end", 2, "end point")
	return fixCmd
}

func init() {
	rootCmd.AddCommand(newFixCmd())
}
