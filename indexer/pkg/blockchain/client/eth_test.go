package client

import (
	"context"
	"testing"

	"github.com/721tools/backend-go/index/pkg/blockchain/model"
	"github.com/721tools/backend-go/index/pkg/utils/hex"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/assert"
)

func TestEthRpcClient_LatestBlockHeight(t *testing.T) {
	client := NewEthRpcClient("https://api.mycryptoapi.com/eth")
	defer client.Release()
	height, err := client.LatestBlockHeight(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, height)
	t.Log("latest block height", height)
}

func TestEthRpcClient_Balance(t *testing.T) {
	client := NewEthRpcClient("https://api.mycryptoapi.com/eth")
	defer client.Release()
	var address hex.Hex
	err := address.ToHex("0x1a3F275b9Af71D597219899151140a0049DB557b")
	assert.NoError(t, err)
	balance, err := client.Balance(context.Background(), address)
	assert.NoError(t, err)
	t.Log("balance", balance)
	text, err := balance.MarshalText()
	assert.NoError(t, err)
	t.Log("balance text", string(text))
}

func TestEthRpcClient_GasPrice(t *testing.T) {
	client := NewEthRpcClient("https://api.mycryptoapi.com/eth")
	defer client.Release()
	price, err := client.GasPrice(context.Background())
	assert.NoError(t, err)
	t.Log("price", price)
}

func TestEthRpcClient_GetChainID(t *testing.T) {
	client := NewEthRpcClient("https://cloudflare-eth.com")
	defer client.Release()
	id, err := client.GetChainID(context.Background())
	assert.NoError(t, err)
	t.Log("chain_id", id, uint64(id))
}

func TestEthRpcClient_BlockByHeight(t *testing.T) {
	client := NewEthRpcClient("https://cloudflare-eth.com")
	defer client.Release()
	var blockNumber math.HexOrDecimal64 = 2000000
	height, err := client.BlockByHeight(context.Background(), blockNumber, true)
	assert.NoError(t, err)
	t.Log("height", blockNumber, "info", height)
}

func TestEthRpcClient_BlocksByHeight(t *testing.T) {
	client := NewEthRpcClient("https://wild-empty-darkness.discover.quiknode.pro/87866012a3b481a1851cbb91c9687c44b2b55eaa/")
	defer client.Release()
	var batch []math.HexOrDecimal64
	batch = append(batch, 15139880)
	blocks, err := client.BlocksByHeight(context.Background(), batch, true)
	assert.NoError(t, err)
	for idx := range blocks {
		t.Log("height", blocks[idx].Number)
		for txIdx := range blocks[idx].Transactions {
			t.Log("tx", blocks[idx].Transactions[txIdx].Hash, "receipt gas used", blocks[idx].Transactions[txIdx].Receipt.GasUsed)
		}
	}
}

func TestEthRpcClient_TxByHash(t *testing.T) {
	client := NewEthRpcClient("https://api.mycryptoapi.com/eth")
	defer client.Release()
	var hex hex.Hex
	err := hex.ToHex("0xe394cd682186087c659c338d7b896b5a3146465ec46051297cec8eb359d3106c")
	assert.NoError(t, err)
	tx, err := client.TxByHash(context.Background(), hex)
	assert.NoError(t, err)
	t.Log("tx_hash", tx.Hash, "receipt", tx.Receipt.TransactionHash)
}

func TestEthRpcClient_TxsByHash(t *testing.T) {
	var hexes []hex.Hex
	var add1 hex.Hex
	err := add1.ToHex("0xadbebd1c39af0e6c48daf47ed515900e1296dcebde225c702339cedf58bca5f0")
	assert.NoError(t, err)
	var add2 hex.Hex
	err = add2.ToHex("0xe394cd682186087c659c338d7b896b5a3146465ec46051297cec8eb359d3106c")
	assert.NoError(t, err)
	hexes = append(hexes, add1, add2)
	client := NewEthRpcClient("https://api.mycryptoapi.com/eth")
	defer client.Release()
	txs, err := client.TxsByHash(context.Background(), hexes)
	assert.NoError(t, err)
	for idx := range txs {
		t.Log("idx", idx, "hash", txs[idx].Hash, "receipt hash", txs[idx].Receipt.TransactionHash)
	}
}

func TestEthRpcClient_GetContractERC20(t *testing.T) {
	client := NewEthRpcClient("https://api.mycryptoapi.com/eth")
	defer client.Release()
	var hex hex.Hex
	// https://etherscan.io/token/0xdac17f958d2ee523a2206206994597c13d831ec7
	err := hex.ToHex("0xdac17f958d2ee523a2206206994597c13d831ec7")
	assert.NoError(t, err)
	contract, err := client.GetContract(context.Background(), hex)
	assert.NoError(t, err)
	assert.Equal(t, model.ERC20, contract.Type)
	assert.Equal(t, "Tether USD", contract.Name)
	t.Log("token info", contract)
}

func TestEthRpcClient_GetContractSelfDestructAsEOA(t *testing.T) {
	client := NewEthRpcClient("https://api.mycryptoapi.com/eth")
	defer client.Release()
	var hex hex.Hex
	// https://etherscan.io/address/0xcbce61316759d807c474441952ce41985bbc5a40#code
	err := hex.ToHex("0xcbce61316759d807c474441952ce41985bbc5a40")
	assert.NoError(t, err)
	contract, err := client.GetContract(context.Background(), hex)
	assert.NoError(t, err)
	assert.Equal(t, model.EOA, contract.Type)
	t.Log("token info", contract)
}

func TestEthRpcClient_GetContractERC721(t *testing.T) {
	client := NewEthRpcClient("https://api.mycryptoapi.com/eth")
	defer client.Release()
	var hex hex.Hex
	// https://etherscan.io/token/0x0326b0688d9869a19388312df6805d1d72aab7bc
	err := hex.ToHex("0x57f1887a8bf19b14fc0df6fd9b2acc9af147ea85")
	assert.NoError(t, err)
	contract, err := client.GetContract(context.Background(), hex)
	assert.NoError(t, err)
	assert.Equal(t, model.ERC721, contract.Type)
	t.Log("token info", contract)
}

func TestEthRpcClient_ERC20Balances(t *testing.T) {
	client := NewEthRpcClient("https://api.mycryptoapi.com/eth")
	defer client.Release()
	balances := make([]*QueryTokenBalance, 0)
	balances = append(balances, &QueryTokenBalance{
		WalletAddress: hex.HexstrToHex("0xf977814e90da44bfa03b6295a0616a897441acec"),
		TokenAddress:  hex.HexstrToHex("0xdAC17F958D2ee523a2206206994597C13D831ec7"),
	}, &QueryTokenBalance{
		WalletAddress: hex.HexstrToHex("0xf977814e90da44bfa03b6295a0616a897441acec"),
		TokenAddress:  hex.HexstrToHex("0xB8c77482e45F1F44dE1745F52C74426C631bDD52"),
	})
	err := client.ERC20Balances(context.Background(), balances)
	assert.NoError(t, err)
	for idx := range balances {
		t.Log("idx:", idx, " token: ", balances[idx].TokenAddress, " wallet: ", balances[idx].WalletAddress, " balances:", balances[idx].Balance)
	}

}

func TestEthRpcClient_ERC721Balances(t *testing.T) {
	client := NewEthRpcClient("https://green-small-mountain.discover.quiknode.pro/70510f2528a6c57bcb376fda51073da9aceb214f/")
	defer client.Release()
	balances := make([]*QueryTokenBalance, 0)
	balances = append(balances, &QueryTokenBalance{
		WalletAddress: hex.HexstrToHex("0x8d8890235639AA0715aFb141fd2943f19E101e78"),
		TokenID:       hex.IntstrToBigInt("161"),
		TokenAddress:  hex.HexstrToHex("0x413788815307ce10428c624ed20998e2c35bbf2d"),
	}, &QueryTokenBalance{
		WalletAddress: hex.HexstrToHex("0x8d8890235639AA0715aFb141fd2943f19E101e78"),
		TokenID:       hex.IntstrToBigInt("784"),
		TokenAddress:  hex.HexstrToHex("0x667d28ca8a8f4391fe13c92d36e60c7615d2f8db"),
	})
	err := client.ERC721Balances(context.Background(), balances)
	assert.NoError(t, err)
	for idx := range balances {
		t.Log("idx:", idx, " token addrs: ", balances[idx].TokenAddress, " token id: ", balances[idx].TokenID, " wallet: ", balances[idx].WalletAddress, " balances:", balances[idx].Balance)
	}
}

func TestEthRpcClient_ERC1155Balances(t *testing.T) {
	client := NewEthRpcClient("https://api.mycryptoapi.com/eth")
	defer client.Release()
	balances := make([]*QueryTokenBalance, 0)
	balances = append(balances, &QueryTokenBalance{
		WalletAddress: hex.HexstrToHex("0x3bebeb0b42046eb8250e3d4f345cdb227a2ba5a0"),
		TokenAddress:  hex.HexstrToHex("0x495f947276749ce646f68ac8c248420045cb7b5e"),
		TokenID:       hex.IntstrToBigInt("63406052020414813134660098691514919229802284462760201759499335706523780251649"),
	}, &QueryTokenBalance{
		WalletAddress: hex.HexstrToHex("0x24499d606f42b06e94a6bedee6d094b97bfd236e"),
		TokenAddress:  hex.HexstrToHex("0x495f947276749ce646f68ac8c248420045cb7b5e"),
		TokenID:       hex.IntstrToBigInt("10667520178378707509196689156321488950583476902291880253708054277517725402112"),
	})
	err := client.ERC1155Balances(context.Background(), balances)
	assert.NoError(t, err)
	for idx := range balances {
		t.Log("idx:", idx, " token: ", balances[idx].TokenAddress, " wallet: ", balances[idx].WalletAddress, " balances:", balances[idx].Balance)
	}
}

func TestEthRpcClient_GetContract_EOA(t *testing.T) {
	client := NewEthRpcClient("https://api.mycryptoapi.com/eth")
	defer client.Release()
	var hex hex.Hex
	// https://etherscan.io/token/0xee93f5ddfb69418d12ddf41f9400fc6545ec9042
	err := hex.ToHex("0xee93f5ddfb69418d12ddf41f9400fc6545ec9042")
	assert.NoError(t, err)
	contract, err := client.GetContract(context.Background(), hex)
	assert.NoError(t, err)
	t.Log("token info", contract)
}

func TestEthRpcClient_GetContract(t *testing.T) {
	client := NewEthRpcClient("https://api.mycryptoapi.com/eth")
	contract, err := client.GetContract(context.Background(), hex.HexstrToHex("0xccc441ac31f02cd96c153db6fd5fe0a2f4e6a68d "))
	assert.NoError(t, err)
	t.Log("contract info", contract)
}

func TestNewEthRpcClient_GetTx(t *testing.T) {
	client := NewEthRpcClient("https://api.mycryptoapi.com/eth")
	tx, err := client.TxByHash(context.Background(), hex.HexstrToHex("0xbe03790872e51ef0ffe1b5d741bdaa09b4e158a579f721da0725ace53b55b87f "))
	assert.NoError(t, err)
	t.Log("tx info", tx)
	t.Log("receipt", tx.Receipt.TransactionHash)
}

func TestNewEthRpcClient_GetTokenURL(t *testing.T) {
	client := NewEthRpcClient("https://api.mycryptoapi.com/eth")
	url, err := client.TokenURL(context.Background(), hex.HexstrToHex("0x79fcdef22feed20eddacbb2587640e45491b757f"), hex.HexstrToHex("0x6f")) // 111
	t.Log("receipt", url)
	assert.NoError(t, err)
	assert.Equal(t, "ipfs://QmWiQE65tmpYzcokCheQmng2DCM33DEhjXcPB6PanwpAZo/111", url)
}

// go test -count=1 -v -run TestEthRpcClient_BlocksByHeight  github.com/721tools/backend-go/index/pkg/blockchain/client
