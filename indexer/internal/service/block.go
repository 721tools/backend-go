package service

import (
	"context"

	"github.com/721tools/backend-go/indexer/pkg/utils/hex"

	model2 "github.com/721tools/backend-go/indexer/internal/model"
	"github.com/721tools/backend-go/indexer/pkg/db"
	"xorm.io/xorm"
)

type BlockService struct {
	xorm *xorm.Engine
}

func NewBlockService() BlockIface {
	return &BlockService{xorm: db.GetDBEngine()}
}

func (b *BlockService) GetTimeStampByHeights(ctx context.Context, blockNumbers []uint64) []model2.OriginBlock {
	var blocks []model2.OriginBlock
	b.xorm.Table(&model2.OriginBlock{}).In("height", blockNumbers).Find(&blocks)
	return blocks
}

func (b *BlockService) ClearBlocks(ctx context.Context, blockNumbers []uint64) error {
	for _, blockNumber := range blockNumbers {
		if err := b.clearBlock(ctx, blockNumber); err != nil {
			log.Err(ctx, "clear block occur a error", "blockNumber", blockNumber, "err", err)
			return err
		}
	}
	return nil
}

func (b *BlockService) LatestDBBlockHeight(ctx context.Context) uint64 {
	var block model2.OriginBlock
	get, _ := b.xorm.Table(&model2.OriginBlock{}).Where("done = 1").Desc("height").Limit(1).Get(&block)
	if get {
		return block.Height
	}
	return 0
}

func (b *BlockService) DBNotDoneHeights(ctx context.Context) []uint64 {
	var heights []uint64
	_ = b.xorm.Table(&model2.OriginBlock{}).Where("done = ?", false).Cols("height").Find(&heights)
	return heights
}

func (b *BlockService) FindContractInBlock(ctx context.Context, contractAddress hex.Hex, height uint64) bool {
	exist, _ := b.xorm.Table(&model2.OriginTx{}).Where("from = ? or to = ?", contractAddress, contractAddress).Exist()
	if exist {
		return true
	}
	exist, _ = b.xorm.Table(&model2.OriginReceiptLog{}).Where("address = ?", contractAddress).Exist()
	return exist
}

func (b *BlockService) clearBlock(ctx context.Context, blockNumber uint64) error {
	session := b.xorm.NewSession()
	if err := session.Begin(); err != nil {
		return err
	}
	defer session.Close()
	if _, err := session.Table(&model2.OriginBlock{}).Where("height = ?", blockNumber).Delete(); err != nil {
		return err
	}
	if _, err := session.Table(&model2.OriginTx{}).Where("height = ?", blockNumber).Delete(); err != nil {
		return err
	}
	if _, err := session.Table(&model2.OriginReceiptLog{}).Where("height = ?", blockNumber).Delete(); err != nil {
		return err
	}
	if _, err := session.Table(&model2.TokenFlow{}).Where("height = ?", blockNumber).Delete(); err != nil {
		return err
	}
	if _, err := session.Table(&model2.NFTMint{}).Where("height = ?", blockNumber).Delete(); err != nil {
		return err
	}
	if _, err := session.Table(&model2.NFTSales{}).Where("height = ?", blockNumber).Delete(); err != nil {
		return err
	}
	if _, err := session.Table(&model2.Contract{}).Where("height = ?", blockNumber).Delete(); err != nil {
		return err
	}
	return session.Commit()
}
