package model

import (
	"github.com/721tools/backend-go/indexer/pkg/utils/hex"
	"github.com/ethereum/go-ethereum/common/math"
)

// RawTx 链上Transaction 交易记录
type RawTx struct {
	BlockHash        hex.Hex              `json:"blockHash"`
	BlockNumber      math.HexOrDecimal64  `json:"blockNumber"`
	Hash             hex.Hex              `json:"hash"`
	From             hex.Hex              `json:"from"`
	To               hex.Hex              `json:"to"`
	Gas              math.HexOrDecimal256 `json:"gas"`
	GasPrice         math.HexOrDecimal256 `json:"gasPrice"`
	Input            hex.Hex              `json:"input"`
	Nonce            math.HexOrDecimal64  `json:"nonce"`
	TransactionIndex math.HexOrDecimal64  `json:"transactionIndex"`
	Value            math.HexOrDecimal256 `json:"value"`
	Receipt          *RawTxReceipt        `json:"receipt"`
	Timestamp        math.HexOrDecimal64  `json:"timestamp"`
}

func (r *RawTx) GetMethodID() hex.Hex {
	if len(r.Input) > 4 {
		return r.Input[:4]
	}
	return nil
}

func (r *RawTx) GetInputBytes() []byte {
	return r.Input
}
