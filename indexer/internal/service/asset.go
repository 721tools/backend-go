package service

import (
	"context"

	model2 "github.com/721tools/backend-go/index/internal/model"
	"github.com/721tools/backend-go/index/pkg/blockchain/client"
	"github.com/721tools/backend-go/index/pkg/blockchain/model"
	"github.com/721tools/backend-go/index/pkg/db"
	"github.com/721tools/backend-go/index/pkg/utils/hex"
	"xorm.io/xorm"
)

type AssetSvc struct {
	client client.Client
	orm    *xorm.Engine
}

func NewAssetSvc() AssetIface {
	return &AssetSvc{orm: db.GetDBEngine(), client: client.GetClient()}
}

func (a *AssetSvc) GetERC20BalanceByAccount(ctx context.Context, wallet hex.Hex) (balances []ERC20Balance, err error) {
	var tokens []hex.Hex
	var tokenFlows []model2.TokenFlow
	session := a.orm.NewSession()
	err = session.Table(&model2.TokenFlow{}).Distinct("token_flow.address").
		Join("inner", "contract", "contract.address = token_flow.address").
		Where("contract.type = ?", model.ERC20.ToUInt()).
		Where("token_flow.to = ?", wallet).Find(&tokenFlows)
	if err != nil {
		return
	}
	for idx := range tokenFlows {
		tokens = append(tokens, tokenFlows[idx].Address)
	}
	rawBalances := make([]*client.QueryTokenBalance, 0)
	for _, token := range tokens {
		rawBalances = append(rawBalances, &client.QueryTokenBalance{
			WalletAddress: wallet,
			TokenAddress:  token,
		})
	}
	err = a.client.ERC20Balances(ctx, rawBalances)
	if err != nil {
		return
	}
	balances = make([]ERC20Balance, 0)
	tokenAddress := make([]interface{}, 0)
	for _, rawBalance := range rawBalances {
		if rawBalance.Balance.GreaterZero() {
			tokenAddress = append(tokenAddress, rawBalance.TokenAddress)
		}
	}
	if len(tokenAddress) == 0 {
		return balances, nil
	}
	var tokenInfos []model2.Contract
	err = session.Table(&model2.Contract{}).In("address", tokenAddress).Find(&tokenInfos)
	if err != nil {
		return balances, err
	}
	m := make(map[string]model2.Contract)
	for idx := range tokenInfos {
		m[tokenInfos[idx].Address.String()] = tokenInfos[idx]
	}
	for _, rawBalance := range rawBalances {
		if !rawBalance.Balance.GreaterZero() {
			continue
		}
		tokenInfo, ok := m[rawBalance.TokenAddress.String()]
		if !ok {
			continue
		}
		balances = append(balances, ERC20Balance{
			TokenInfo: TokenInfo{
				Height:      tokenInfo.Height,
				Address:     tokenInfo.Address,
				IsStandard:  tokenInfo.IsStandard,
				Name:        tokenInfo.Name,
				Symbol:      tokenInfo.Symbol,
				Decimal:     tokenInfo.Decimal,
				TotalSupply: tokenInfo.TotalSupply,
			},
			Token:  rawBalance.TokenAddress,
			Amount: rawBalance.Balance,
		})
	}
	return balances, nil
}

func (a *AssetSvc) GetERC721BalanceByAccount(ctx context.Context, wallet hex.Hex) (balances []ERC721Balance, err error) {
	var nfts []struct {
		Token   hex.Hex
		TokenID hex.Hex
	}
	var tokenFlows []model2.TokenFlow
	session := a.orm.NewSession()
	err = session.Table(&model2.TokenFlow{}).Distinct("token_flow.address, token_flow.value").
		Join("inner", "contract", "contract.address = token_flow.address").
		Where("contract.type = ?", model.ERC721.ToUInt()).
		Where("token_flow.to = ?", wallet).Find(&tokenFlows)
	if err != nil {
		return
	}
	for idx := range tokenFlows {
		nfts = append(nfts, struct {
			Token   hex.Hex
			TokenID hex.Hex
		}{Token: tokenFlows[idx].Address, TokenID: tokenFlows[idx].Value})
	}
	if len(nfts) == 0 {
		return
	}

	rawBalances := make([]*client.QueryTokenBalance, 0)
	for _, nft := range nfts {
		rawBalances = append(rawBalances, &client.QueryTokenBalance{
			WalletAddress: wallet,
			TokenAddress:  nft.Token,
			TokenID:       hex.HexstrToBigInt(nft.TokenID.HexStr()),
		})
	}

	err = a.client.ERC721Balances(ctx, rawBalances)
	if err != nil {
		return
	}

	tokenAddress := make([]interface{}, 0)
	for _, rawBalance := range rawBalances {
		if rawBalance.Balance.GreaterZero() {
			tokenAddress = append(tokenAddress, rawBalance.TokenAddress)
		}
	}
	if len(tokenAddress) == 0 {
		return
	}

	var tokenInfos []model2.Contract
	err = session.Table(&model2.Contract{}).In("address", tokenAddress).Find(&tokenInfos)
	if err != nil {
		return balances, err
	}
	m := make(map[string]model2.Contract)
	for idx := range tokenInfos {
		m[tokenInfos[idx].Address.String()] = tokenInfos[idx]
	}

	for _, rawBalance := range rawBalances {
		if !rawBalance.Balance.GreaterZero() {
			continue
		}
		tokenInfo, ok := m[rawBalance.TokenAddress.String()]
		if !ok {
			continue
		}
		balances = append(balances, ERC721Balance{
			TokenInfo: TokenInfo{
				Height:      tokenInfo.Height,
				Address:     tokenInfo.Address,
				IsStandard:  tokenInfo.IsStandard,
				Name:        tokenInfo.Name,
				Symbol:      tokenInfo.Symbol,
				Decimal:     tokenInfo.Decimal,
				TotalSupply: tokenInfo.TotalSupply,
			},
			Token:   rawBalance.TokenAddress,
			TokenID: rawBalance.TokenID,
		})
	}
	return
}
