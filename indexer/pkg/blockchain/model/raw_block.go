package model

import (
	"github.com/721tools/backend-go/index/pkg/utils/hex"
	"github.com/ethereum/go-ethereum/common/math"
)

// RawBlock 链上区块信息
type RawBlock struct {
	Number           math.HexOrDecimal64  `json:"number"`
	Hash             hex.Hex              `json:"hash"`
	ParentHash       hex.Hex              `json:"parentHash"`
	MixHash          hex.Hex              `json:"mixHash"`
	Nonce            hex.Hex              `json:"nonce"`
	Sha3Uncles       hex.Hex              `json:"sha3Uncles"`
	LogsBloom        hex.Hex              `json:"logsBloom"`
	TransactionsRoot hex.Hex              `json:"transactionsRoot"`
	StateRoot        hex.Hex              `json:"stateRoot"`
	ReceiptsRoot     hex.Hex              `json:"receiptsRoot"`
	Miner            hex.Hex              `json:"miner"`
	Difficulty       math.HexOrDecimal256 `json:"difficulty"`
	TotalDifficulty  math.HexOrDecimal256 `json:"totalDifficulty"`
	ExtraData        hex.Hex              `json:"extraData"`
	Size             math.HexOrDecimal64  `json:"size"`
	GasLimit         math.HexOrDecimal64  `json:"gasLimit"`
	GasUsed          math.HexOrDecimal64  `json:"gasUsed"`
	Timestamp        math.HexOrDecimal64  `json:"timestamp"`
	Transactions     []RawTx              `json:"transactions"`
}
