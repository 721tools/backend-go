package service

import (
	"testing"

	"github.com/721tools/backend-go/index/pkg/db"
	"github.com/721tools/backend-go/index/pkg/utils/hex"
)

func TestGetNFTExtraDataByContractAndTokenId(t *testing.T) {
	dsn := "root:@tcp(localhost:3306)/721tools_ethereum_nft_assets?parseTime=true"
	db.NewDBEngine(dsn)
	service := NewNFTItemService()
	addrs := hex.HexstrToHex("0x91673149ffae3274b32997288395d07a8213e41f")
	id := hex.IntstrToBigInt("230").Bytes()
	item := service.GetNFTExtraDataByContractAndTokenId(addrs, id)

	// for i, _ := range item {
	// 	t.Log("debug tx", " height", tx[i].Height, " hash", tx[i].TxHash, " from ", tx[i].From, " to ", tx[i].To, " token", tx[i].Symbol)
	// }
	t.Log("debug tx", " height", item)
}
