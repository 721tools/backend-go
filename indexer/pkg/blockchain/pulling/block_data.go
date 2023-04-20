package pulling

import (
	"context"
	"fmt"

	"github.com/721tools/backend-go/index/pkg/blockchain/client"
	"github.com/721tools/backend-go/index/pkg/blockchain/model"
	"github.com/721tools/backend-go/index/pkg/utils/retry"
	"github.com/ethereum/go-ethereum/common/math"
)

type BlockData struct {
	// pull blockchain data params
	Start          uint64
	End            uint64
	heights        []math.HexOrDecimal64
	SkipInternalTx bool

	// block data
	Data   []model.RawBlock
	Cancel bool

	// evm client
	client client.Client
	ctx    context.Context
}

func NewBlockData(ctx context.Context, start, end uint64, skipInternal bool, client client.Client) *BlockData {
	data := &BlockData{Start: start, End: end, SkipInternalTx: skipInternal, ctx: ctx, client: client, Cancel: false}
	for i := int(start); i < int(end); i++ {
		data.heights = append(data.heights, math.HexOrDecimal64(i))
	}
	return data
}

func NewCancelBlockData() *BlockData {
	return &BlockData{Cancel: true}
}

func (b *BlockData) GetBlocksFromNode() error {
	success, err := retry.Exec(func() (bool, error) {
		log.Info("get heights", "height", b.heights)
		blocks, err := b.client.BlocksByHeight(b.ctx, b.heights, b.SkipInternalTx)
		if err != nil {
			log.Error("get block from node occur a failed", "heights", b.heights, "err", err)
			return false, err
		}
		b.Data = append(b.Data, blocks...)
		return true, nil
	}, 3)

	if !success {
		return fmt.Errorf("get Data from node occur a error, err is %v", err)
	}
	return nil
}

func (b *BlockData) Done(ctx context.Context, f func(context context.Context, heights []math.HexOrDecimal64) error) error {
	return f(ctx, b.heights)
}
