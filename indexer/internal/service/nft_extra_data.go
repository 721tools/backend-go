package service

import (
	model2 "github.com/721tools/backend-go/index/internal/model"
	"github.com/721tools/backend-go/index/pkg/db"
	"github.com/721tools/backend-go/index/pkg/utils/hex"
	"xorm.io/xorm"
)

type NFTItemService struct {
	xorm *xorm.Engine
}

type NFTExtraData struct {
	Item model2.NFTItem
	Col  model2.NFTCollection
}

func NewNFTItemService() *NFTItemService {
	return &NFTItemService{xorm: db.GetDBEngine()}
}

func (n *NFTItemService) GetNFTExtraDataByContractAndTokenId(contract_address, token_address hex.Hex) *NFTExtraData {
	var item model2.NFTItem
	var item2 model2.NFTCollection
	get0, _ := n.xorm.Table(&model2.NFTItem{}).Where("contract_address = ? and token_id = ?", contract_address, token_address).Desc("update_time").Limit(1).Get(&item)
	get1, _ := n.xorm.Table(&model2.NFTCollection{}).Where("contract_address = ?", contract_address).Desc("update_time").Limit(1).Get(&item2)

	log.Info("debug orm", "item", item, "item2", item2)
	if get0 && get1 {
		return &NFTExtraData{
			Item: item,
			Col:  item2,
		}
	}
	return nil
}
