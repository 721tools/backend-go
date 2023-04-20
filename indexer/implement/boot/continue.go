package boot

import (
	"context"
	"time"

	"github.com/721tools/backend-go/indexer/implement/bus"
	"github.com/721tools/backend-go/indexer/implement/service"
	"github.com/721tools/backend-go/indexer/pkg/blockchain/client"
	"github.com/721tools/backend-go/indexer/pkg/blockchain/pulling"
	"github.com/721tools/backend-go/indexer/pkg/utils/log16"
)

var log = log16.NewLogger("module", "boot")

func Continuous(cacheSize, PullStep uint8) {
	svc := service.NewBlockService()
	puller := pulling.NewPuller(cacheSize, PullStep, client.GetClient())
	heights := svc.DBNotDoneHeights(context.Background())
	fixDBBlocks(svc, puller, heights)
	latestDBHeight := svc.LatestDBBlockHeight(context.Background())
	if latestDBHeight == 0 {
		genesisBlock(svc, puller)
	}
	log.Info("start to run pulling block", "start", latestDBHeight)
	go puller.Continuous(latestDBHeight)
	for {
		select {
		case pipe := <-puller.Cache:
			if puller.Cancel {
				return
			}
			convert := bus.NewConvert()
			convert.SetRawBlocks(pipe.Data)
			convert.DispatchEvents()
			_ = pipe.Done(context.Background(), svc.Done)

		case <-puller.Quit.QuitChan:
			log.Warn("send cancel to puller to stop pulling data from block chain.")
			puller.Cancel = true
			log.Warn("wait 5s to shutdown program")
			time.Sleep(5 * time.Second)
			return
		}
	}
}

func fixDBBlocks(svc service.BlockIface, puller *pulling.Puller, heights []uint64) {
	for _, height := range heights {
		if err := svc.ClearBlocks(context.Background(), []uint64{height}); err != nil {
			log.Err(context.Background(), "clear block failed", "err", err)
			break
		}
	}
}

func genesisBlock(svc service.BlockIface, puller *pulling.Puller) {
	fixDBBlocks(svc, puller, []uint64{0})
}
