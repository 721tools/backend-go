package block_listener

import (
	"context"
	"encoding/json"
	"math"

	"github.com/721tools/backend-go/index/pkg/consts"
	"github.com/721tools/backend-go/index/pkg/mq"
)

type GasUsedListener struct {
	mq *mq.RedisMQ
}

func NewGasUsedListener() *GasUsedListener {
	return &GasUsedListener{mq: mq.GetMQ()}
}

func (t *GasUsedListener) Handle(event *Event) (err error) {
	gas := &GasPerBlock{}
	gas.Height = event.Block.Height
	gas.GasLimit = event.Block.GasLimit
	gas.GasUsed = event.Block.GasUsed
	gas.GasUsedRatio = float64(gas.GasUsed) / float64(gas.GasLimit)
	gas.Timestamp = event.Block.Timestamp
	sum := uint64(0)
	max := uint64(0)
	min := ^uint64(0)
	for _, rawTx := range event.Block.Txs {
		gas_price := rawTx.GasPrice.Uint64()
		if gas_price > max {
			max = gas_price
		}
		if gas_price < min {
			min = gas_price
		}
		sum += gas_price
	}

	ulen := uint64(len(event.Block.Txs))
	if ulen == 0 {
		return nil
	}

	gas.MaxGasPriceGwei = max / uint64(math.Pow10(9))
	gas.MinGasPriceGwei = min / uint64(math.Pow10(9))
	gas.AvgGasPriceGwei = sum / ulen / uint64(math.Pow10(9))

	mq_ := mq.GetMQ()
	gas_str, _ := json.Marshal(gas)
	mq_.Publish(context.TODO(), consts.GAS_PRICE_NOW_MQ, string(gas_str))
	mq_.Set(context.TODO(), consts.GAS_PRICE_NOW_KEY, string(gas_str), 0)

	log.Info("debug gas used in current block", "height", event.Block.Height, "gas is", gas)
	return nil
}
