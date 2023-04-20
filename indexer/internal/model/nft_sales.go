package model

import (
	"time"
)

type NFTSales struct {
	Height    uint64    `xorm:"height"`
	Timestamp time.Time `xorm:"timestamp"`
	TxHash    string    `xorm:"tx_hash"`
	LogIndex  uint      `xorm:"logIndex"`
	Plateform uint      `xorm:"plateform"`
	IsBundle  uint      `xorm:"isBundle"`
	Address   string    `xorm:"address"`
	TokenId   string    `xorm:"tokenId"`
	Amount    string    `xorm:"amount"`
	Buyer     string    `xorm:"buyer"`
	Seller    string    `xorm:"seller"`
	PriceETH  string    `xorm:"priceETH"`
	Direction uint      `xorm:"direction"`
}

func (n NFTSales) TableName() string {
	return "nft_trades"
}
