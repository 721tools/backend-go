package model

import "github.com/721tools/backend-go/indexer/pkg/utils/hex"

// OriginInTxLog deprecate
type OriginInTxLog struct {
	Height          uint    `xorm:"height"`
	TxHash          hex.Hex `xorm:"tx_hash"`
	From            hex.Hex `xorm:"from"`
	To              hex.Hex `xorm:"to"`
	Value           hex.Hex `xorm:"value"`
	ContractAddress hex.Hex `xorm:"contract_address"`
}

func (o OriginInTxLog) TableName() string {
	return "origin_in_tx_log"
}
