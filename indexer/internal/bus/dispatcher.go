package bus

import (
	"context"
	"fmt"

	"github.com/721tools/backend-go/indexer/internal/bus/block_listener"
)

var blockRawDataRegisters []block_listener.Listener

func Init() {
	blockRawDataRegisters = append(
		blockRawDataRegisters,
		block_listener.NewOriginDBStartListener(),
		block_listener.NewContractCreateListener(), // 同一个区块内，需先执行合约创建逻辑
		block_listener.NewTokenFlowListener(),
		block_listener.NewNFTSalesListener(),
		block_listener.NewNFTMintListener(),
		block_listener.NewBlurSalesListener(),
		block_listener.NewGasUsedListener(),
		block_listener.NewOriginDBEndListener(),
	)
}

func DispatchRawBlockData(event *block_listener.Event) {
	for _, register := range blockRawDataRegisters {
		if err := register.Handle(event); err != nil {
			log.Err(context.Background(), "occur a error", "err", err)
			panic(fmt.Errorf("bus dispatch raw block data occur a err, err is %v", err))
		}
	}
}
