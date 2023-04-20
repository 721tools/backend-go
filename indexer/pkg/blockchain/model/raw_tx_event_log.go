package model

import (
	"github.com/721tools/backend-go/index/pkg/blockchain/alg"
	"github.com/721tools/backend-go/index/pkg/utils/hex"
	"github.com/ethereum/go-ethereum/common/math"
)

var (
	transferEventsSig = []alg.SigID{
		alg.EventSig("Transfer(address,address,uint256)"),
		alg.EventSig("TransferSingle(address,address,address,uint256,uint256)"),
		alg.EventSig("TransferBatch(address,address,address,uint256[],uint256[])")}
)

// RawTxEventLog 链上EventLog
type RawTxEventLog struct {
	BlockNumber      math.HexOrDecimal64 `json:"blockNumber"`
	BlockHash        hex.Hex             `json:"blockHash"`
	TransactionHash  hex.Hex             `json:"transactionHash"`
	TransactionIndex math.HexOrDecimal64 `json:"transactionIndex"`
	Address          hex.Hex             `json:"address"`
	Topics           []hex.Hex           `json:"topics"`
	Data             hex.Hex             `json:"data"`
	LogIndex         math.HexOrDecimal64 `json:"logIndex"`
	Removed          bool                `json:"removed"`
}

// https://etherscan.io/tx/0x428c4cb8dd41cf55904594a7558cdff3f5573cb8ea5a6afd39b33688fde6d189#eventlog
func (r *RawTxEventLog) GetLogBytes() []byte {
	var bytes []byte
	if len(r.Topics) == 0 {
		return bytes
	}

	for idx, topic := range r.Topics {
		if idx == 0 {
			continue
		}
		bytes = append(bytes, topic...)
	}
	if len(r.Data) == 0 {
		return bytes
	}
	bytes = append(bytes, r.Data...)
	return bytes
}

func (r *RawTxEventLog) GetTopic0() hex.Hex {
	if len(r.Topics) == 0 {
		return nil
	}
	return r.Topics[0]
}

func (r *RawTxEventLog) IsTokenFlow() bool {
	if len(r.Topics) == 0 {
		return false
	}
	topic0 := r.GetTopic0()
	for _, sig := range transferEventsSig {
		if topic0.EqualTo(sig.ToHex()) {
			return true
		}
	}
	return false
}
