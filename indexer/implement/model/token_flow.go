package model

import (
	"github.com/721tools/backend-go/indexer/pkg/utils/hex"
)

type TokenFlow struct {
	Height   uint64  `xorm:"height"`
	TxHash   hex.Hex `xorm:"tx_hash"`
	LogIndex int64   `xorm:"log_index"`
	Address  hex.Hex `xorm:"address"`
	From     hex.Hex `xorm:"from"`
	To       hex.Hex `xorm:"to"`
	Value    hex.Hex `xorm:"value"`
	Amount   hex.Hex `xorm:"amount"`
}

func (t TokenFlow) TableName() string {
	return "token_flow"
}

func (t TokenFlow) IsMint() bool {
	return t.From.EqualTo(hex.EmptyHexStr)
}

func (t TokenFlow) ToMintNFTToken(minter hex.Hex) NFTMint {
	return NFTMint{
		Height:  t.Height,
		TxHash:  t.TxHash,
		Address: t.Address,
		TokenID: t.Value,
		Minter:  minter,
		Amount:  t.Amount,
	}
}
