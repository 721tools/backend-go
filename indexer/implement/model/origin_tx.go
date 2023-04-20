package model

import (
	"time"

	"github.com/721tools/backend-go/indexer/pkg/blockchain/model"

	"github.com/721tools/backend-go/indexer/pkg/utils/hex"
)

type OriginTx struct {
	TxHash          hex.Hex             `xorm:"tx_hash pk"`
	TxIdx           uint64              `xorm:"tx_idx pk"`
	Height          uint64              `xorm:"height"`
	From            hex.Hex             `xorm:"from"`
	To              hex.Hex             `xorm:"to"`
	Value           hex.BigInt          `xorm:"value"`
	Status          uint64              `xorm:"status"`
	ErrCode         uint                `xorm:"err_code"`
	ErrMsg          hex.Hex             `xorm:"err_msg"`
	MethodID        hex.Hex             `xorm:"method_id"`
	MethodName      hex.Hex             `xorm:"method_name"`
	MethodArgs      hex.Hex             `xorm:"method_args"`
	ContractAddress hex.Hex             `xorm:"contract_address"`
	GasLimit        hex.BigInt          `xorm:"gas_limit"`
	GasPrice        hex.BigInt          `xorm:"gas_price"`
	GasUsed         hex.BigInt          `xorm:"gas_used"`
	Nonce           uint64              `xorm:"nonce"`
	Timestamp       time.Time           `xorm:"timestamp"`
	Input           hex.Hex             `xorm:"-"`
	ReceiptLogs     []OriginReceiptLog  `xorm:"-"`
	RawTxReceipt    *model.RawTxReceipt `xorm:"-"`
}

func (o OriginTx) TableName() string {
	return "origin_tx"
}
