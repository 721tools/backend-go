package model

import (
	"time"

	"github.com/721tools/backend-go/indexer/pkg/blockchain/model"
	"github.com/721tools/backend-go/indexer/pkg/utils/hex"
)

// OriginBlock 区块链原始block
type OriginBlock struct {
	Height     uint64          `xorm:"pk"`
	BlockHash  hex.Hex         `xorm:"block_hash"`
	ParentHash hex.Hex         `xorm:"parent_hash"`
	GasLimit   uint64          `xorm:"gas_limit"`
	GasUsed    uint64          `xorm:"gas_used"`
	Timestamp  time.Time       `xorm:"timestamp"`
	CreatedAt  time.Time       `xorm:"created_at"`
	Done       bool            `xorm:"done"`
	Txs        []OriginTx      `xorm:"-"`
	RawBlock   *model.RawBlock `xorm:"-"`
}

func (o *OriginBlock) TableName() string {
	return "origin_block"
}
