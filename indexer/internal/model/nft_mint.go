package model

import (
	"github.com/721tools/backend-go/indexer/pkg/utils/hex"
)

type NFTMint struct {
	Height   uint64  `xorm:"height"`
	TxHash   hex.Hex `xorm:"tx_hash"`
	Address  hex.Hex `xorm:"address"`
	TokenID  hex.Hex `xorm:"token_id"`
	Minter   hex.Hex `xorm:"minter"`
	Amount   hex.Hex `xorm:"amount"`
	TokenURL string  `xorm:"token_url"`
	Metadata string  `xorm:"metadata"`
}

func (n NFTMint) TableName() string {
	return "nft_mint"
}
