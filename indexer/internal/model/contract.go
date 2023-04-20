package model

import "github.com/721tools/backend-go/indexer/pkg/utils/hex"

type Contract struct {
	Height      uint64     `xorm:"height"`
	TxHash      hex.Hex    `xorm:"tx_hash"`
	Address     hex.Hex    `xorm:"address"`
	Type        uint       `xorm:"type"`
	IsStandard  bool       `xorm:"is_standard"`
	Name        string     `xorm:"name"`
	Symbol      string     `xorm:"symbol"`
	Decimal     uint64     `xorm:"decimal"`
	TotalSupply hex.BigInt `xorm:"total_supply"`
}

func (c Contract) TableName() string {
	return "contract"
}
