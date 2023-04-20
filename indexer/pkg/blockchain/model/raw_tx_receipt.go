package model

import (
	"github.com/721tools/backend-go/indexer/pkg/utils/hex"
	"github.com/ethereum/go-ethereum/common/math"
)

// RawTxReceipt 链上收据数据
type RawTxReceipt struct {
	BlockHash         hex.Hex              `json:"blockHash"`
	BlockNumber       math.HexOrDecimal64  `json:"blockNumber"`
	TransactionHash   hex.Hex              `json:"transactionHash"`
	TransactionIndex  math.HexOrDecimal64  `json:"transactionIndex"`
	Type              math.HexOrDecimal64  `json:"type"`
	From              hex.Hex              `json:"from"`
	To                hex.Hex              `json:"to"`
	CumulativeGasUsed math.HexOrDecimal256 `json:"cumulativeGasUsed"`
	GasUsed           math.HexOrDecimal256 `json:"gasUsed"`
	EffectiveGasPrice math.HexOrDecimal256 `json:"effectiveGasPrice"`
	ContractAddress   hex.Hex              `json:"contractAddress"`
	Logs              []RawTxEventLog      `json:"logs"`
	LogsBloom         hex.Hex              `json:"logsBloom"`
	Status            math.HexOrDecimal64  `json:"status"`
}
