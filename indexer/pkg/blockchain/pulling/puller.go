package pulling

import (
	"context"
	"fmt"
	"time"

	"github.com/721tools/backend-go/indexer/internal/bus"
	"github.com/721tools/backend-go/indexer/pkg/utils/quit"

	"github.com/721tools/backend-go/indexer/pkg/blockchain/client"
	"github.com/721tools/backend-go/indexer/pkg/blockchain/model"
	"github.com/721tools/backend-go/indexer/pkg/utils/log16"
	"github.com/721tools/backend-go/indexer/pkg/utils/retry"
	"github.com/ethereum/go-ethereum/common/math"
)

var log = log16.NewLogger("module", "pulling")

type Puller struct {
	Cache     chan *BlockData
	cacheSize uint8
	pullStep  uint8

	lastCheckBlock *model.RawBlock
	client         client.Client

	Cancel bool

	Quit *quit.Quit
}

func NewPuller(cacheSize uint8, pullStep uint8, client client.Client) *Puller {
	q := quit.NewQuit()
	q.WatchOsSignal()
	return &Puller{cacheSize: cacheSize, Cache: make(chan *BlockData, cacheSize), client: client, pullStep: pullStep, Quit: q, Cancel: false}
}

// Continuous 持续 start -> last block height pulling block data
func (p *Puller) Continuous(start uint64) {
	block, err := p.client.BlockByHeight(context.Background(), math.HexOrDecimal64(start), true)
	if err != nil {
		panic(fmt.Errorf("get block failed, err is %v", err))
	}
	p.lastCheckBlock = &block
	for {
		if p.Cancel {
			break
		}
		time.Sleep(time.Second)
		if len(p.Cache) > int(p.cacheSize) {
			continue
		}
		var lastBlockHeight math.HexOrDecimal64
		exec, err := retry.Exec(func() (bool, error) {
			lastBlockHeight, err = p.client.LatestBlockHeight(context.Background())
			if err != nil {
				log.Warn("get latest block height failed")
				return false, err
			}
			return true, nil
		}, 3)
		if !exec {
			panic(fmt.Errorf("get latest block height failed, and err is %v", err))
		}

		if uint64(lastBlockHeight) < uint64(p.lastCheckBlock.Number) {
			log.Warn("last block height < start", "start", start, "end", uint64(lastBlockHeight))
			continue
		}
		for i := uint64(p.lastCheckBlock.Number) + 1; i <= uint64(lastBlockHeight); i += uint64(p.pullStep) {
			if p.Cancel {
				break
			}
			var s, e = i, i + uint64(p.pullStep)
			if e > uint64(lastBlockHeight) {
				e = uint64(lastBlockHeight)
			}
			if s >= e {
				break
			}
			log.Info(fmt.Sprintf("pulling start %d To end %d", s, e))
			blockData := NewBlockData(context.Background(), s, e, true, p.client)
			if err = blockData.GetBlocksFromNode(); err != nil {
				log.Critical(context.Background(), "get blocks from node occur a error", "err", err)
			}
			log.Info("pull 2", "data len", len(blockData.Data))
			for idx := range blockData.Data {
				if !p.lastCheckBlock.Hash.EqualTo(blockData.Data[idx].ParentHash) {
					err = fmt.Errorf("can not match last check height %d, and pull height %d", p.lastCheckBlock.Number, blockData.Data[idx].Number)
					panic(err)
				}
				p.lastCheckBlock = &blockData.Data[idx]
			}
			p.Cache <- blockData
			log.Info("pull block data", "len", len(p.Cache))
		}
		if p.Cancel {
			break
		}
	}
}

// Temporarily 从start -> end pulling block data
func (p *Puller) Temporarily(start, end uint64) {
	if start > end {
		return
	}
	for i := start; i <= end; i += uint64(p.pullStep) {
		var s, e = i, i + uint64(p.pullStep)
		if e > end {
			e = end
		}
		if s >= e {
			break
		}
		blockData := NewBlockData(context.Background(), s, e, true, p.client)

		if err := blockData.GetBlocksFromNode(); err != nil {
			log.Critical(context.Background(), "get blocks from node occur a error", "err", err)
		}
		p.Cache <- blockData
		log.Info("pull block data", "len", len(p.Cache))
	}
	p.Cache <- NewCancelBlockData()
}

func (p *Puller) DealWithHeights(heights []uint64) {
	if len(heights) == 0 {
		return
	}
	for _, height := range heights {
		blockData := NewBlockData(context.Background(), height, height+1, true, p.client)
		if err := blockData.GetBlocksFromNode(); err != nil {
			log.Critical(context.Background(), "get blocks from node occur a error", "err", err)
		}
		p.Cache <- blockData
	}
}

// Pop data from Cache
func (p *Puller) Pop() *BlockData {
	return <-p.Cache
}

func (p *Puller) Fix(height uint64) error {
	blockData := NewBlockData(context.Background(), height, height+1, true, p.client)
	if err := blockData.GetBlocksFromNode(); err != nil {
		log.Critical(context.Background(), "get blocks from node occur a error", "err", err)
	}
	convert := bus.NewConvert()
	convert.SetRawBlocks(blockData.Data)
	convert.DispatchEvents()
	return nil
}
