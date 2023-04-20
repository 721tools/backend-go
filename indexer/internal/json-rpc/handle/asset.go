package handle

import (
	"context"
	"encoding/json"

	"github.com/721tools/backend-go/indexer/internal/service"
	"github.com/721tools/backend-go/indexer/pkg/blockchain/client"
	"github.com/721tools/backend-go/indexer/pkg/blockchain/model"
	"github.com/721tools/backend-go/indexer/pkg/consts"
	"github.com/721tools/backend-go/indexer/pkg/utils/hex"
	"github.com/721tools/backend-go/indexer/pkg/utils/log16"
	"github.com/osamingo/jsonrpc/v2"
)

var log = log16.NewLogger("module", "rpc.handle")

type EmptyReq struct{}

type WalletReq struct {
	Wallet hex.Hex `json:"wallet"`
}

type AddressAndNFTItem struct {
	Wallet       hex.Hex    `json:"wallet"`
	TokenAddress hex.Hex    `json:"token_address"`
	TokenID      hex.BigInt `json:"token_id"`
	Owning       bool       `json:"owning"`
}

type AddressAndNFTItemReq struct {
	Items []*AddressAndNFTItem `json:"items"`
}

type AddressAndNFTItemRes struct {
	Items []*AddressAndNFTItem `json:"items"`
}

type WalletPagenationReq struct {
	Wallet hex.Hex `json:"wallet"`
	Page   uint64  `json:"page"`
	Limit  uint64  `json:"limit"`
}

type ERC721AssetResp struct {
	Balances []service.ERC721Balance `json:"balances"`
}

type ERC20AssetResp struct {
	Balances []service.ERC20Balance `json:"balances"`
}

type UserBalanceResp struct {
	Wallet hex.Hex                     `json:"wallet"`
	ERC20  []*client.QueryTokenBalance `json:"erc20"`
	ERC721 []*client.QueryTokenBalance `json:"erc721"`
	ETH    hex.BigInt                  `json:"eth_balance"`
}

type TransactionResp struct {
	// token info: address, erc20 or erc721, tx
	// from to
	// blockid
	Wallet hex.Hex      `json:"wallet"`
	Tx     []service.Tx `json:"tx"`
	Total  uint64       `json:"total"`
	Limit  uint64       `json:"limit"`
	Page   uint64       `json:"page"`
}

func GetUserTx(c context.Context, rawReq *json.RawMessage) (interface{}, *jsonrpc.Error) {
	var req *WalletPagenationReq
	if err := json.Unmarshal(*rawReq, &req); err != nil {
		return nil, &jsonrpc.Error{Code: jsonrpc.ErrorCodeInvalidRequest, Message: "req invalid", Data: nil}
	}

	resp := TransactionResp{}
	resp.Tx = make([]service.Tx, 0)

	wallet := req.Wallet
	page := req.Page
	limit := req.Limit

	resp.Wallet = wallet
	db := service.NewAcountService()
	db.GetAcountTx(context.Background(), wallet, page, limit, &resp.Tx)

	for i, tx := range resp.Tx {
		if tx.Address.EqualTo(consts.NativeToken) {
			resp.Tx[i].Type = uint64(consts.NativeTokenType)
			resp.Tx[i].Decimal = uint64(consts.NativeTokenDecimal)
			resp.Tx[i].Name = consts.NativeTokenName
			resp.Tx[i].Symbol = consts.NativeTokenSymbol
			resp.Tx[i].TotalSupply = consts.NativeTokenTotalSupply
		}
	}
	resp.Total, _ = db.GetAcountTxCnt(context.Background(), wallet)
	resp.Page = page
	resp.Limit = limit
	return &resp, nil
}

func GetUserBalance(c context.Context, rawReq *json.RawMessage) (interface{}, *jsonrpc.Error) {
	var erc20Req []*WalletReq
	if err := json.Unmarshal(*rawReq, &erc20Req); err != nil {
		return nil, &jsonrpc.Error{Code: jsonrpc.ErrorCodeInvalidRequest, Message: "req invalid", Data: nil}
	}
	req := erc20Req[0]
	wallet := req.Wallet
	db := service.NewAcountService()
	balances := make([]service.TokenBalance, 0)
	db.GetAcountBalance(context.TODO(), wallet, &balances)
	erc20_balances := make([]*client.QueryTokenBalance, 0)
	erc721_balances := make([]*client.QueryTokenBalance, 0)

	eth_client := client.GetClient()
	for _, balance := range balances {
		if balance.Type == model.ERC20.ToUInt() {
			erc20_balances = append(erc20_balances, &client.QueryTokenBalance{
				WalletAddress: wallet,
				TokenAddress:  balance.Address,
				Name:          balance.Name,
				Symbol:        balance.Symbol,
				Decimal:       balance.Decimal,
			})
		} else if balance.Type == model.ERC721.ToUInt() {
			erc721_balances = append(erc721_balances, &client.QueryTokenBalance{
				WalletAddress: wallet,
				TokenAddress:  balance.Address,
				TokenID:       hex.HexstrToBigInt(balance.Value.String()),
				Name:          balance.Name,
				Symbol:        balance.Symbol,
			})
		}
	}

	if len(erc20_balances) > 0 {
		eth_client.ERC20Balances(context.TODO(), erc20_balances)
	}

	if len(erc721_balances) > 0 {
		eth_client.ERC721Balances(context.TODO(), erc721_balances)
	}

	eth_balance, _ := eth_client.Balance(context.TODO(), wallet)

	resp := UserBalanceResp{
		Wallet: wallet,
		ERC20:  erc20_balances,
		ERC721: erc721_balances,
		ETH:    hex.BigInt(eth_balance),
	}
	return &resp, nil
}

// GetUserErc20TokenBalance 返回erc20 token 余额信息
func GetUserErc20TokenBalance(ctx context.Context, rawReq *json.RawMessage) (interface{}, *jsonrpc.Error) {
	var erc20Req []*WalletReq
	if err := json.Unmarshal(*rawReq, &erc20Req); err != nil {
		return nil, &jsonrpc.Error{Code: jsonrpc.ErrorCodeInvalidRequest, Message: "req invalid", Data: nil}
	}
	walletReq := erc20Req[0]
	svc := service.NewAssetSvc()
	erc20Balance, err := svc.GetERC20BalanceByAccount(ctx, walletReq.Wallet)
	if err != nil {
		log.Err(ctx, "query erc20 balances failed", "err", err)
		return nil, &jsonrpc.Error{Code: jsonrpc.ErrorCodeInternal, Message: "erc20 balance occur a error", Data: nil}
	}
	resp := ERC20AssetResp{
		Balances: erc20Balance,
	}
	return &resp, nil
}

// GetUserErc721TokenBalance 返回erc721 余额信息
func GetUserErc721TokenBalance(ctx context.Context, rawReq *json.RawMessage) (interface{}, *jsonrpc.Error) {
	var erc721Req []*WalletReq
	if err := json.Unmarshal(*rawReq, &erc721Req); err != nil {
		return nil, &jsonrpc.Error{Code: jsonrpc.ErrorCodeInvalidRequest, Message: "req invalid", Data: nil}
	}
	walletReq := erc721Req[0]
	svc := service.NewAssetSvc()
	erc721Balance, err := svc.GetERC721BalanceByAccount(ctx, walletReq.Wallet)
	if err != nil {
		log.Err(ctx, "query erc721 balances failed", "err", err)
		return nil, &jsonrpc.Error{Code: jsonrpc.ErrorCodeInternal, Message: "erc20 balance occur a error", Data: nil}
	}
	resp := ERC721AssetResp{
		Balances: erc721Balance,
	}
	return &resp, nil
}

// 批量查询地址是否拥有某个 NFT
// req: {
// address - contract address - tokenid
// }
// res:
// address - contract address - tokenid - true | false
func BatchCheckErc721TokenOwner(ctx context.Context, rawReq *json.RawMessage) (interface{}, *jsonrpc.Error) {
	var items []*AddressAndNFTItem
	var err error

	log.Info("output input", "raw req=", string(*rawReq))
	if err = json.Unmarshal(*rawReq, &items); err != nil {
		log.Info("output input", "err", err.Error())
		return nil, &jsonrpc.Error{Code: jsonrpc.ErrorCodeInvalidParams, Message: "req invalid", Data: nil}
	}

	eth_client := client.GetClient()
	rawBalances := make([]*client.QueryTokenBalance, 0)
	for _, i := range items {
		rawBalances = append(rawBalances, &client.QueryTokenBalance{
			WalletAddress: i.Wallet,
			TokenAddress:  i.TokenAddress,
			TokenID:       i.TokenID,
		})
	}

	// use `ownerOf(tokenId)` call to check owner address
	err = eth_client.ERC721Balances(ctx, rawBalances)
	if err != nil {
		return nil, &jsonrpc.Error{Code: jsonrpc.ErrorCodeInvalidRequest, Message: "failed to get token balance", Data: nil}
	}

	for _, i := range items {
		i.Owning = false
		for _, b := range rawBalances {
			log.Info("output rawBalances", "rawBalances TokenAddress=", b.TokenAddress, "Balance", b.Balance)
			if i.TokenAddress.EqualTo(b.TokenAddress) && i.TokenID.EqualTo(&b.TokenID) && b.Balance.GreaterZero() {
				i.Owning = true
				break
			}
		}
		log.Info("output input", "items TokenAddress=", i.TokenAddress, "owning", i.Owning)
	}

	resp := AddressAndNFTItemRes{
		Items: items,
	}
	return &resp, nil
}
