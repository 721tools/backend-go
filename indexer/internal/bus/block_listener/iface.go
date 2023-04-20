package block_listener

import (
	"time"

	model2 "github.com/721tools/backend-go/indexer/internal/model"
)

type chunkSize = int64

type GasPerBlock struct {
	Height          uint64    `json:"height"`
	GasLimit        uint64    `json:"gas_limit"`
	GasUsed         uint64    `json:"gas_used"`
	GasUsedRatio    float64   `json:"gas_used_ratio"`
	Timestamp       time.Time `json:"ts"`
	MaxGasPriceGwei uint64    `json:"max_gas_price_gwei"`
	MinGasPriceGwei uint64    `json:"min_gas_price_gwei"`
	AvgGasPriceGwei uint64    `json:"avg_gas_price_gwei"`
}

type Event struct {
	Block *model2.OriginBlock
}

func NewEvent(block *model2.OriginBlock) *Event {
	return &Event{Block: block}
}

type Listener interface {
	Handle(event *Event) error
}
