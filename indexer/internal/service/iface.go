package service

import (
	"context"

	"github.com/ethereum/go-ethereum/common/math"

	model2 "github.com/721tools/backend-go/indexer/internal/model"
	"github.com/721tools/backend-go/indexer/pkg/blockchain/model"
	"github.com/721tools/backend-go/indexer/pkg/utils/hex"
	"github.com/721tools/backend-go/indexer/pkg/utils/log16"
)

var log = log16.NewLogger("module", "service")

type ContractIface interface {
	Sync(ctx context.Context, contract *model2.Contract) error
	GetContractType(ctx context.Context, address hex.Hex) model.ContractType
}

type BlockIface interface {
	ClearBlocks(ctx context.Context, blockNumbers []uint64) error
	LatestDBBlockHeight(ctx context.Context) uint64
	DBNotDoneHeights(ctx context.Context) []uint64
	GetTimeStampByHeights(ctx context.Context, blockNumbers []uint64) []model2.OriginBlock
	FindContractInBlock(ctx context.Context, contractAddress hex.Hex, height uint64) bool
	Done(ctx context.Context, heights []math.HexOrDecimal64) error
}

// AssetIface 资产接口

type TokenInfo struct {
	Height      uint64     `json:"height"`
	Address     hex.Hex    `json:"address"`
	IsStandard  bool       `json:"isStandard"`
	Name        string     `json:"name"`
	Symbol      string     `json:"symbol"`
	Decimal     uint64     `json:"decimal"`
	TotalSupply hex.BigInt `json:"totalSupply"`
}

type Tx struct {
	From        hex.Hex    `json:"from"`
	To          hex.Hex    `json:"to"`
	TxHash      hex.Hex    `json:"tx_hash"`
	Height      uint64     `json:"height"`
	Address     hex.Hex    `json:"address"`
	Type        uint64     `json:"type"`
	Value       hex.BigInt `json:"value"`
	Name        string     `json:"name"`
	Symbol      string     `json:"symbol"`
	Decimal     uint64     `json:"decimal"`
	TotalSupply hex.BigInt `json:"totalSupply"`
}

type ERC20Balance struct {
	Token     hex.Hex    `json:"token"`
	TokenInfo TokenInfo  `json:"tokenInfo"`
	Amount    hex.BigInt `json:"amount"`
}

type ERC721Balance struct {
	Token     hex.Hex    `json:"token"`
	TokenID   hex.BigInt `json:"tokenID"`
	TokenInfo TokenInfo  `json:"tokenInfo"`
}

type AssetIface interface {
	GetERC20BalanceByAccount(ctx context.Context, wallet hex.Hex) ([]ERC20Balance, error)
	GetERC721BalanceByAccount(ctx context.Context, wallet hex.Hex) ([]ERC721Balance, error)
}

type AccountIface interface {
	GetAcountBalance(ctx context.Context, address hex.Hex, balance *[]TokenBalance) error
	GetAcountTx(ctx context.Context, address hex.Hex, balance *[]TokenBalance) error
}
