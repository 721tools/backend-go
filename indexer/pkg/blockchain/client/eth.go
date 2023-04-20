package client

import (
	"context"
	"fmt"

	"github.com/721tools/backend-go/indexer/pkg/blockchain/alg"
	"github.com/721tools/backend-go/indexer/pkg/blockchain/model"
	"github.com/721tools/backend-go/indexer/pkg/blockchain/pool"
	"github.com/721tools/backend-go/indexer/pkg/utils/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/rpc"
)

type EthRpcClient struct {
	endpoint string
	pool     pool.EthRpcPool
}

func NewEthRpcClient(endpoint string, opts ...pool.Option) Client {
	rpcPool, err := pool.NewEthRPCPool(endpoint, opts...)
	if err != nil {
		panic(fmt.Errorf("new Eth RPC client failed, err is %w", err))
	}
	client := &EthRpcClient{endpoint: endpoint, pool: rpcPool}
	return client
}

// GetChainID 获取chain network ID
func (e *EthRpcClient) GetChainID(ctx context.Context) (chainID math.HexOrDecimal64, err error) {
	err = e.pool.Run(func(client *rpc.Client) error {
		return client.Call(&chainID, "eth_chainId")
	})
	return
}

// Balance native token balance
func (e *EthRpcClient) Balance(ctx context.Context, address hex.Hex) (balance math.HexOrDecimal256, err error) {
	err = e.pool.Run(func(client *rpc.Client) error {
		return client.Call(&balance, "eth_getBalance", address, "latest")
	})
	return
}

func (e *EthRpcClient) ERC20Balances(ctx context.Context, query []*QueryTokenBalance) (err error) {
	batch := make([]rpc.BatchElem, 0)
	methodID := alg.ERC20BalanceOfMethod.ToHex()
	result := make([]string, len(query))
	for idx := range query {
		var methodCall []byte
		methodCall = append(methodCall, methodID.Bytes()...)
		methodCall = append(methodCall, common.LeftPadBytes(query[idx].WalletAddress.Bytes(), 32)...)
		args := EthCallParam{
			To:   query[idx].TokenAddress.HexStr(),
			Data: hexutil.Encode(methodCall),
			Gas:  "0x30000",
		}
		batch = append(batch, rpc.BatchElem{
			Method: "eth_call",
			Args:   []interface{}{args, "latest"},
			Result: &result[idx],
		})
	}
	if err = e.pool.BatchEthCall(batch, true); err != nil {
		return
	}
	for idx := range query {
		query[idx].Balance = hex.HexstrToBigInt(result[idx])
	}
	return nil
}

func (e *EthRpcClient) ERC721Balances(ctx context.Context, query []*QueryTokenBalance) (err error) {
	batch := make([]rpc.BatchElem, 0)
	methodID := alg.ERC721OwnerOfMethod.ToHex()
	result := make([]string, len(query))
	for idx := range query {
		var methodCall []byte
		methodCall = append(methodCall, methodID.Bytes()...)
		methodCall = append(methodCall, common.LeftPadBytes(query[idx].TokenID.Bytes(), 32)...)
		args := EthCallParam{
			To:   query[idx].TokenAddress.HexStr(),
			Data: hexutil.Encode(methodCall),
			Gas:  "0x30000",
		}
		batch = append(batch, rpc.BatchElem{
			Method: "eth_call",
			Args:   []interface{}{args, "latest"},
			Result: &result[idx],
		})
	}

	if err = e.pool.BatchEthCall(batch, false); err != nil {
		return
	}
	for idx := range query {
		bytes := hex.HexstrToHex(result[idx]).Bytes()
		if len(bytes) > 20 {
			address := hex.Hex(bytes[len(bytes)-20:])
			if query[idx].WalletAddress.EqualTo(address) {
				query[idx].Balance = hex.IntstrToBigInt("1")
				continue
			}
		}
		query[idx].Balance = hex.IntstrToBigInt("0")
	}
	return nil
}

func (e *EthRpcClient) ERC1155Balances(ctx context.Context, query []*QueryTokenBalance) (err error) {
	batch := make([]rpc.BatchElem, 0)
	methodID := alg.ERC1155BalanceOfMethod.ToHex()
	result := make([]string, len(query))
	for idx := range query {
		var methodCall []byte
		methodCall = append(methodCall, methodID.Bytes()...)
		methodCall = append(methodCall, common.LeftPadBytes(query[idx].WalletAddress.Bytes(), 32)...)
		methodCall = append(methodCall, common.LeftPadBytes(query[idx].TokenID.Bytes(), 32)...)
		args := EthCallParam{
			To:   query[idx].TokenAddress.HexStr(),
			Data: hexutil.Encode(methodCall),
			Gas:  "0x30000",
		}
		batch = append(batch, rpc.BatchElem{
			Method: "eth_call",
			Args:   []interface{}{args, "latest"},
			Result: &result[idx],
		})
	}

	if err = e.pool.BatchEthCall(batch, true); err != nil {
		return
	}
	for idx := range query {
		query[idx].Balance = hex.HexstrToBigInt(result[idx])
	}
	return nil
}

// LatestBlockHeight 获取节点最新的块高
func (e *EthRpcClient) LatestBlockHeight(ctx context.Context) (height math.HexOrDecimal64, err error) {
	err = e.pool.Run(func(client *rpc.Client) error {
		return client.Call(&height, "eth_blockNumber")
	})
	return
}

// GasPrice 查询gas price, return the current price per gas in wei.
func (e *EthRpcClient) GasPrice(ctx context.Context) (price math.HexOrDecimal64, err error) {
	err = e.pool.Run(func(client *rpc.Client) error {
		return client.Call(&price, "eth_gasPrice")
	})
	return
}

// BlockByHeight 根据height 查询block info
func (e *EthRpcClient) BlockByHeight(ctx context.Context, height math.HexOrDecimal64, skip bool) (block model.RawBlock, err error) {
	err = e.pool.Run(func(client *rpc.Client) error {
		return client.Call(&block, "eth_getBlockByNumber", height, true)
	})
	if err != nil {
		return
	}

	if len(block.Transactions) == 0 {
		return
	}

	elements := make([]rpc.BatchElem, 0)
	for idx := range block.Transactions {
		block.Transactions[idx].Timestamp = block.Timestamp
		block.Transactions[idx].Receipt = new(model.RawTxReceipt)
		elements = append(elements, rpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []interface{}{block.Transactions[idx].Hash},
			Result: block.Transactions[idx].Receipt,
		})
	}
	if err = e.pool.BatchEthCall(elements, true); err != nil {
		return
	}
	for idx := range block.Transactions {
		if block.Transactions[idx].Receipt.TransactionHash != nil && !block.Transactions[idx].Receipt.TransactionHash.EqualTo(block.Transactions[idx].Hash) {
			err = fmt.Errorf("retrieve receipt failed, block height is %d and tx hash is %s", height, block.Transactions[idx].Hash)
			return
		}
	}
	return
}

// BlocksByHeight 批量获取区块信息
func (e *EthRpcClient) BlocksByHeight(ctx context.Context, heights []math.HexOrDecimal64, skip bool) (blocks []model.RawBlock, err error) {
	if len(heights) == 0 {
		return
	}
	var batchBlockElements []rpc.BatchElem
	blocks = make([]model.RawBlock, len(heights))
	for heightIdx, height := range heights {
		batchBlockElements = append(batchBlockElements, rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{height, true},
			Result: &blocks[heightIdx]})
	}
	log.Info("batch eth call blocks", "element_len", len(batchBlockElements))
	if err = e.pool.BatchEthCall(batchBlockElements, true); err != nil {
		return
	}
	// check blocks
	for heightIdx, height := range heights {
		if height != blocks[heightIdx].Number {
			return blocks, fmt.Errorf("target is %d, real is %d", height, blocks[heightIdx].Number)
		}
	}
	// get receipts
	var batchReceiptElements []rpc.BatchElem
	for blockIdx := range blocks {
		for txIdx := range blocks[blockIdx].Transactions {
			blocks[blockIdx].Transactions[txIdx].Timestamp = blocks[blockIdx].Timestamp
			blocks[blockIdx].Transactions[txIdx].Receipt = new(model.RawTxReceipt)
			batchReceiptElements = append(batchReceiptElements, rpc.BatchElem{
				Method: "eth_getTransactionReceipt",
				Args:   []interface{}{blocks[blockIdx].Transactions[txIdx].Hash},
				Result: blocks[blockIdx].Transactions[txIdx].Receipt,
			})
		}
	}
	if len(batchReceiptElements) == 0 {
		return
	}

	log.Info("batch eth call tx receipts", "heights", heights, "element_len", len(batchReceiptElements))
	if err = e.pool.BatchEthCall(batchReceiptElements, true); err != nil {
		return
	}
	// check receipts & merge receipt to transaction
	for blockIdx := range blocks {
		for txIdx := range blocks[blockIdx].Transactions {
			if blocks[blockIdx].Transactions[txIdx].Receipt == nil || !blocks[blockIdx].Transactions[txIdx].Receipt.TransactionHash.EqualTo(blocks[blockIdx].Transactions[txIdx].Hash) {
				err = fmt.Errorf("retrieve receipt failed, block height is %d and tx hash is %s", blocks[blockIdx].Number, blocks[blockIdx].Transactions[txIdx].Hash)
				return
			}
		}
	}
	return
}

// TxByHash 根据hash 获取交易数据
func (e *EthRpcClient) TxByHash(ctx context.Context, hash hex.Hex) (tx model.RawTx, err error) {
	var receipt model.RawTxReceipt
	var batchElements []rpc.BatchElem
	batchElements = append(batchElements, rpc.BatchElem{
		Method: "eth_getTransactionReceipt",
		Args:   []interface{}{hash},
		Result: &receipt,
	})
	batchElements = append(batchElements, rpc.BatchElem{
		Method: "eth_getTransactionByHash",
		Args:   []interface{}{hash},
		Result: &tx,
	})
	if err = e.pool.BatchEthCall(batchElements, true); err != nil {
		return
	}
	if !receipt.TransactionHash.EqualTo(tx.Hash) {
		err = fmt.Errorf("retrieve receipt failed, tx hash is %s and receipt hash is %s", tx.Hash, receipt.TransactionHash)
	}
	tx.Receipt = &receipt
	return
}

// TxsByHash 根据hashes获取tx数据
func (e *EthRpcClient) TxsByHash(ctx context.Context, hashes []hex.Hex) (txs []model.RawTx, err error) {
	txs = make([]model.RawTx, len(hashes))
	var batchElements []rpc.BatchElem
	for idx := range hashes {
		batchElements = append(batchElements, rpc.BatchElem{
			Method: "eth_getTransactionByHash",
			Args:   []interface{}{hashes[idx]},
			Result: &txs[idx],
		})
		txs[idx].Receipt = new(model.RawTxReceipt)
		batchElements = append(batchElements, rpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []interface{}{hashes[idx]},
			Result: &txs[idx].Receipt,
		})
	}
	if err = e.pool.BatchEthCall(batchElements, true); err != nil {
		return
	}
	for idx := range txs {
		if !txs[idx].Hash.EqualTo(txs[idx].Receipt.TransactionHash) {
			err = fmt.Errorf("idx receit hash not equal tx hash, receit hash is %s, tx hash is %s", txs[idx].Receipt.TransactionHash, txs[idx].Hash)
			return
		}
	}
	return
}

// EthCall 执行内部的eth call
func (e *EthRpcClient) EthCall(ctx context.Context, to hex.Hex, data string, gas string) (result string, err error) {
	err = e.pool.Run(func(client *rpc.Client) error {
		return client.Call(&result, "eth_call", EthCallParam{To: to.HexStr(), Data: data, Gas: gas}, "latest")
	})
	return
}

// BatchEthCall 批量执行eth call
func (e *EthRpcClient) BatchEthCall(ctx context.Context, batch []rpc.BatchElem, allOk bool) error {
	return e.pool.BatchEthCall(batch, allOk)
}

// EthCode 获取address 的code信息
func (e *EthRpcClient) EthCode(ctx context.Context, address hex.Hex) (code string, err error) {
	err = e.pool.Run(func(client *rpc.Client) error {
		return client.Call(&code, "eth_getCode", address.HexStr(), "latest")
	})
	return
}

// GetContract 获取合约信息
func (e *EthRpcClient) GetContract(ctx context.Context, address hex.Hex) (rawContract model.RawContract, err error) {
	rawContract.Address = address
	rawContract.Type = model.Contract
	code, err := e.EthCode(ctx, address)
	if err != nil {
		log.Err(ctx, "get EthCode failed", "address", address, "err", err)
		return
	}
	if len(code) <= 2 {
		rawContract.Type = model.EOA
		return
	}
	if e.erc1155Stander(ctx, address) {
		rawContract.Type = model.ERC1155
		return
	}
	if e.erc721Stander(ctx, address) {
		rawContract.Type = model.ERC721
		rawContract.Name, err = e.name(ctx, address)
		if err != nil {
			log.Warn("eth_call name failed", "address", address, "err", err)
			err = nil
			return
		}
		rawContract.Symbol, err = e.symbol(ctx, address)
		if err != nil {
			log.Warn("eth_call symbol failed", "address", address, "err", err)
			err = nil
			return
		}

		rawContract.TotalSupply, err = e.totalSupply(ctx, address)
		if err != nil {
			log.Warn("eth_call TotalSupply failed", "address", address, "err", err)
			err = nil
			return
		}
		return
	}

	isStander, name, symbol, decimal, supply := e.erc20Stander(ctx, address, code)
	if !isStander {
		return
	}
	rawContract.Type = model.ERC20
	rawContract.Name = name
	rawContract.Symbol = symbol
	rawContract.Decimals = decimal
	rawContract.TotalSupply = supply
	return
}

// GetRpcClient 获取连接池对象
func (e *EthRpcClient) GetRpcClient(ctx context.Context, client *rpc.Client, err error) {
	client, err = e.pool.GetClient()
	return
}

// PutRpcClient 放回连接池
func (e *EthRpcClient) PutRpcClient(ctx context.Context, client *rpc.Client) {
	e.pool.PutClient(client)
}

// Release 释放所有链接池
func (e *EthRpcClient) Release() {
	e.pool.ReleaseAll()
}

// tokenURL 通过eth call tokenURI()获取合约的 tokenURL 信息
func (e *EthRpcClient) TokenURL(ctx context.Context, address, token_id hex.Hex) (url string, err error) {
	var methodCall []byte

	methodID := alg.TokenURIMethods.ToHex()
	methodCall = append(methodCall, methodID.Bytes()...)
	methodCall = append(methodCall, common.LeftPadBytes(token_id.Bytes(), 32)...)

	url, err = e.EthCall(ctx, address, hexutil.Encode(methodCall), "0x30000")
	if url == "0x" || err != nil {
		return "", err
	}
	return hex.TrimHexStrAndDecodeToStr(url), err
}

// name 通过eth call 获取合约的name信息
func (e *EthRpcClient) name(ctx context.Context, address hex.Hex) (name string, err error) {
	for _, item := range alg.ContractNameMethods {
		name, err = e.EthCall(ctx, address, item.ToString(), "0x30000")
		if name == "0x" || err != nil {
			continue
		}
		return hex.TrimHexStrAndDecodeToStr(name), err
	}
	return "", err
}

// symbol 通过eth call symbol()获取合约的symbol信息
func (e *EthRpcClient) symbol(ctx context.Context, address hex.Hex) (symbol string, err error) {
	for _, item := range alg.SymbolMethods {
		symbol, err = e.EthCall(ctx, address, item.ToString(), "0x30000")
		if symbol == "0x" || err != nil {
			continue
		}
		return hex.TrimHexStrAndDecodeToStr(symbol), err
	}
	return "", err
}

// totalSupply 通过eth call supply()获取total supply信息
func (e *EthRpcClient) totalSupply(ctx context.Context, address hex.Hex) (totalSupply math.HexOrDecimal256, err error) {
	result, err := e.EthCall(ctx, address, alg.TotalSuppleMethod.ToString(), "0x30000")
	if err != nil || result == "0" {
		return math.HexOrDecimal256{}, fmt.Errorf("totalSupply is empty")
	}
	err = totalSupply.UnmarshalText([]byte(result))
	if err != nil {
		return math.HexOrDecimal256{}, err
	}
	return totalSupply, nil
}

// decimal 通过eth call decimal() 获取decimal信息
func (e *EthRpcClient) decimal(ctx context.Context, address hex.Hex) (decimal math.HexOrDecimal64, err error) {
	for _, item := range alg.DecimalMethods {
		decimalString, err := e.EthCall(ctx, address, item.ToString(), "0x10000")
		if decimalString == "0x" || err != nil {
			continue
		}
		if err = decimal.UnmarshalText([]byte(decimalString)); err != nil {
			continue
		}
		return decimal, nil
	}
	return 0, fmt.Errorf("decimal is empty")
}

// erc721Stander erc721标准
// 参考文档： https://docs.openzeppelin.com/contracts/4.x/api/utils#IERC165-supportsInterface-bytes4-
func (e *EthRpcClient) erc721Stander(ctx context.Context, address hex.Hex) bool {
	var result math.HexOrDecimal64
	var data []byte
	data = append(data, alg.SupportsInterfaceMethod...)
	data = append(data, common.RightPadBytes(hex.HexstrToHex("0x80ac58cd"), 32)...)
	err := e.pool.Run(func(client *rpc.Client) error {
		return client.Call(&result, "eth_call", EthCallParam{
			To:   address.HexStr(),
			Data: hex.Hex(data).HexStr(),
			Gas:  "0x30000"}, "latest")
	})
	if err != nil {
		return false
	}
	return result > 0
}

// erc1155Stander erc1155标准
// 参考文档： https://docs.openzeppelin.com/contracts/4.x/api/utils#IERC165-supportsInterface-bytes4-
func (e *EthRpcClient) erc1155Stander(ctx context.Context, address hex.Hex) bool {
	var result math.HexOrDecimal64
	var data []byte
	data = append(data, alg.SupportsInterfaceMethod...)
	data = append(data, common.RightPadBytes(hex.HexstrToHex("0xd9b67a26"), 32)...)
	err := e.pool.Run(func(client *rpc.Client) error {
		return client.Call(&result, "eth_call", EthCallParam{
			To:   address.HexStr(),
			Data: hex.Hex(data).HexStr(),
			Gas:  "0x30000"}, "latest")
	})
	if err != nil {
		return false
	}
	return result > 0
}

func (e *EthRpcClient) erc20Stander(ctx context.Context, address hex.Hex, code string) (isErc20 bool, name string, symbol string, decimal math.HexOrDecimal64, supply math.HexOrDecimal256) {
	//if !strings.Contains(code, alg.TransferMethod.NoOxPrefix()) || !strings.Contains(code, alg.TransferFromMethod.NoOxPrefix()) {
	//	return false, "", "", math.HexOrDecimal64(0), *math.NewHexOrDecimal256(0)
	//}
	name, err := e.name(ctx, address)
	if err != nil {
		return false, "", "", math.HexOrDecimal64(0), *math.NewHexOrDecimal256(0)
	}
	symbol, err = e.symbol(ctx, address)
	if err != nil {
		return false, "", "", math.HexOrDecimal64(0), *math.NewHexOrDecimal256(0)
	}
	decimal, err = e.decimal(ctx, address)
	if err != nil {
		return false, "", "", math.HexOrDecimal64(0), *math.NewHexOrDecimal256(0)
	}
	supply, err = e.totalSupply(ctx, address)
	if err != nil {
		return false, "", "", math.HexOrDecimal64(0), *math.NewHexOrDecimal256(0)
	}
	return true, name, symbol, decimal, supply
}
