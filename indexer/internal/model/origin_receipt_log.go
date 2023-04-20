package model

import (
	"encoding/json"

	"github.com/721tools/backend-go/index/pkg/blockchain/abi_parse"
	"github.com/721tools/backend-go/index/pkg/utils/hex"
)

type OriginReceiptLog struct {
	Height       uint64  `xorm:"height"`
	TxHash       hex.Hex `xorm:"tx_hash"`
	LogIndex     uint64  `xorm:"log_index"`
	Address      hex.Hex `xorm:"address"`
	EventSig     hex.Hex `xorm:"event_sig"`
	EventName    hex.Hex `xorm:"event_name"`
	EventPayload hex.Hex `xorm:"event_payload"`
}

func (o OriginReceiptLog) TableName() string {
	return "origin_receipt_log"
}

func (o *OriginReceiptLog) ToPayload() abi_parse.Args {
	args := make(abi_parse.Args)
	_ = json.Unmarshal(o.EventPayload, &args)
	return args
}

func (o *OriginReceiptLog) ToTokenFlow() TokenFlow {
	args := abi_parse.HexToArgs(o.EventPayload)
	return TokenFlow{
		Height:   o.Height,
		TxHash:   o.TxHash,
		LogIndex: int64(o.LogIndex),
		Address:  o.Address,
		From:     hex.HexstrToHex(args["from"].(string)),
		To:       hex.HexstrToHex(args["to"].(string)),
		Value:    hex.HexstrToHex(args["value"].(string)),
		Amount:   hex.HexstrToHex(args["amount"].(string)),
	}
}
