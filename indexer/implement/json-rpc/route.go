package json_rpc

import (
	"context"

	"github.com/721tools/backend-go/indexer/implement/json-rpc/handle"
	"github.com/721tools/backend-go/indexer/pkg/utils/log16"
	"github.com/gin-gonic/gin"
	"github.com/osamingo/jsonrpc/v2"
)

var log = log16.NewLogger("module", "json-rpc")

func Route() gin.HandlerFunc {
	mr := jsonrpc.NewMethodRepository()

	// BatchCheckErc721TokenOwner
	err := mr.RegisterMethod("asset.get_user_balance", jsonrpc.HandlerFunc(handle.GetUserBalance), &[]*handle.WalletReq{}, handle.UserBalanceResp{})
	if err != nil {
		log.Critical(context.Background(), "register asset.get_user_balance failed", "err", err)
	}

	err = mr.RegisterMethod("asset.check_nft_owner", jsonrpc.HandlerFunc(handle.BatchCheckErc721TokenOwner), &[]*handle.AddressAndNFTItemReq{}, handle.AddressAndNFTItemRes{})
	if err != nil {
		log.Critical(context.Background(), "register asset.check_nft_owner failed", "err", err)
	}

	err = mr.RegisterMethod("asset.get_user_tx", jsonrpc.HandlerFunc(handle.GetUserTx), &[]*handle.WalletPagenationReq{}, handle.TransactionResp{})
	if err != nil {
		log.Critical(context.Background(), "register asset.get_user_balance failed", "err", err)
	}

	err = mr.RegisterMethod("asset.erc20_balance", jsonrpc.HandlerFunc(handle.GetUserErc20TokenBalance), &[]*handle.WalletReq{}, handle.ERC20AssetResp{})
	if err != nil {
		log.Critical(context.Background(), "register asset.erc20_balance failed", "err", err)
	}
	err = mr.RegisterMethod("asset.erc721_balance", jsonrpc.HandlerFunc(handle.GetUserErc721TokenBalance), &[]*handle.WalletReq{}, handle.ERC721AssetResp{})
	if err != nil {
		log.Critical(context.Background(), "register asset.erc20_balance failed", "err", err)
	}

	err = mr.RegisterMethod("gas.get_gas_now", jsonrpc.HandlerFunc(handle.GASNOW), &[]*handle.EmptyReq{}, handle.GasNowResp{})
	if err != nil {
		log.Critical(context.Background(), "register asset.get_gas_price_now failed", "err", err)
	}

	err = mr.RegisterMethod("block.ts", jsonrpc.HandlerFunc(handle.Ts), &[]*handle.HeightReq{}, []handle.TimeStampResp{})
	if err != nil {
		log.Critical(context.Background(), "register block.ts", "err", err)
	}

	err = mr.RegisterMethod("nft.get_most_popular", jsonrpc.HandlerFunc(handle.GetTheMostPopularCollections), &[]*handle.EmptyReq{}, []handle.CollectionsResp{})
	if err != nil {
		log.Critical(context.Background(), "register block.ts", "err", err)
	}
	return gin.WrapF(mr.ServeHTTP)
}
