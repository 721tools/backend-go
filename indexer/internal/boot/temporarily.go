package boot

import (
	"context"
	"time"

	"github.com/721tools/backend-go/index/internal/bus"
	"github.com/721tools/backend-go/index/internal/service"
	"github.com/721tools/backend-go/index/pkg/blockchain/client"
	"github.com/721tools/backend-go/index/pkg/blockchain/pulling"
)

func Temporarily(start, end uint64, cacheSize, pullStep uint8) {
	svc := service.NewBlockService()
	var blockNumbers []uint64
	for i := start; i < end; i++ {
		blockNumbers = append(blockNumbers, i)
	}
	if err := svc.ClearBlocks(context.Background(), blockNumbers); err != nil {
		return
	}
	puller := pulling.NewPuller(cacheSize, pullStep, client.GetClient())
	go puller.Temporarily(start, end)
	for {
		select {
		case pipe := <-puller.Cache:
			if pipe.Cancel {
				time.Sleep(5 * time.Second)
				return
			}
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
