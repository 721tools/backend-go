package log16

import (
	log "github.com/inconshreveable/log15"
)

type Option func(logger *Logger)

func WithFileOption(logger *Logger) {
	logger.SetHandler(log.MultiHandler(
		//log.StreamHandler(os.Stderr, log.LogfmtFormat()),
		log.LvlFilterHandler(log.LvlError, log.Must.FileHandler("bamboo.json", log.JsonFormat())),
	))
}
