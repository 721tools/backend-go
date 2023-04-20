package client

import (
	"context"
	"fmt"

	"github.com/721tools/backend-go/index/pkg/blockchain/model"
	"github.com/721tools/backend-go/index/pkg/blockchain/pool"
	"github.com/721tools/backend-go/index/pkg/utils/hex"
	"github.com/721tools/backend-go/index/pkg/utils/log16"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	log = log16.NewLogger("module", "eth_client")
)

type Client interface {
	GetChainID(ctx context.Context) (chainID math.HexOrDecimal64, err error)
	Balance(ctx context.Context, address hex.Hex) (balance math.HexOrDecimal256, err error)
	ERC20Balances(ctx context.Context, query []*QueryTokenBalance) (err error)
	ERC721Balances(ctx context.Context, query []*QueryTokenBalance) (err error)
	ERC1155Balances(ctx context.Context, query []*QueryTokenBalance) (err error)
	LatestBlockHeight(ctx context.Context) (height math.HexOrDecimal64, err error)
	GasPrice(ctx context.Context) (price math.HexOrDecimal64, err error)
	BlockByHeight(ctx context.Context, height math.HexOrDecimal64, skip bool) (block model.RawBlock, err error)
	BlocksByHeight(ctx context.Context, heights []math.HexOrDecimal64, skip bool) (blocks []model.RawBlock, err error)
	TxByHash(ctx context.Context, hash hex.Hex) (tx model.RawTx, err error)
	TxsByHash(ctx context.Context, hashes []hex.Hex) (txs []model.RawTx, err error)
	GetContract(ctx context.Context, address hex.Hex) (rawContract model.RawContract, err error)
	EthCall(ctx context.Context, to hex.Hex, data string, gas string) (result string, err error)
	BatchEthCall(ctx context.Context, batch []rpc.BatchElem, allOk bool) error
	EthCode(ctx context.Context, address hex.Hex) (code string, err error)
	GetRpcClient(ctx context.Context, client *rpc.Client, err error)
	PutRpcClient(ctx context.Context, client *rpc.Client)
	Release()
	TokenURL(ctx context.Context, address, token_id hex.Hex) (url string, err error)
}

type QueryTokenBalance struct {
	WalletAddress hex.Hex    `json:"walletAddress"`
	TokenAddress  hex.Hex    `json:"token_address"`
	TokenID       hex.BigInt `json:"token_id"`
	Balance       hex.BigInt `json:"balance"`
	Name          string     `json:"name"`
	Symbol        string     `json:"symbol"`
	Decimal       uint64     `json:"decimal"`
}

type EthCallParam struct {
	To   string `json:"to"`
	Data string `json:"data"`
	Gas  string `json:"gas"`
}

type EvmClientType int

const (
	EvmClient  = EvmClientType(0)
	TestClient = EvmClientType(1)
)

var client Client

func NewClient(rpc string, clientType EvmClientType, opts ...pool.Option) Client {
	if clientType == EvmClient {
		client = NewEthRpcClient(rpc, opts...)
	}
	return client
}

func GetClient() Client {
	if client == nil {
		panic(fmt.Errorf("client not init, please check it"))
	}
	return client
}
