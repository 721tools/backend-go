package inter

import (
	"context"
	"os"

	"github.com/721tools/backend-go/indexer/pkg/utils/log16"
	"github.com/spf13/cobra"
)

var (
	debug   bool
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "bamboo",
	Short: "daemon for bamboo",
}

var log = log16.NewLogger("module", "rootCmd")

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logger")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bamboo.yaml)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Critical(context.Background(), "root cmd exec failed", "err", err)
		os.Exit(1)
	}
}

func initConfig() {}
