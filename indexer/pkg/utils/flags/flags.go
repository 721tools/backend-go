package flags

import (
	"github.com/721tools/backend-go/indexer/pkg/utils/log16"
	"github.com/spf13/pflag"
)

var log = log16.NewLogger("module", "flags")

func PrintFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		log.Debug("FLAG: ", flag.Name, flag.Value)
	})
}
