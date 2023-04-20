package bus

import (
	"sync"
	"time"

	"github.com/721tools/backend-go/indexer/implement/bus/block_listener"
	model2 "github.com/721tools/backend-go/indexer/implement/model"
	"github.com/721tools/backend-go/indexer/implement/service"
	"github.com/721tools/backend-go/indexer/pkg/blockchain/abi_parse"
	"github.com/721tools/backend-go/indexer/pkg/blockchain/model"
	"github.com/721tools/backend-go/indexer/pkg/utils/hex"
	"github.com/721tools/backend-go/indexer/pkg/utils/log16"
)

var log = log16.NewLogger("module", "bus")

type Convert struct {
	svc       service.ContractIface
	rawBlocks []model.RawBlock
	jieData   *abi_parse.JieData
}

func NewConvert() *Convert {
	return &Convert{
		rawBlocks: make([]model.RawBlock, 0),
		jieData:   abi_parse.NewJieData(),
		svc:       service.NewContract()}
}

func (c *Convert) SetRawBlocks(rawBlocks []model.RawBlock) {
	c.rawBlocks = append(c.rawBlocks, rawBlocks...)
}

func (c *Convert) DispatchEvents() {
	// 处理 raw block data
	blockEvents := c.extractBlocks()
	log.Info("extract blocks", "len", len(blockEvents))
	wg := sync.WaitGroup{}

	for _, event := range blockEvents {
		log.Info("range event add", "blockid is", event.Block.Height)
		wg.Add(1)
		go func(event *block_listener.Event) {
			defer wg.Done()
			DispatchRawBlockData(event)
			log.Info("range event done", "blockid is", event.Block.Height)
		}(event)
	}
	wg.Wait()
}

func (c *Convert) extractBlocks() []*block_listener.Event {
	events := make([]*block_listener.Event, 0)
	for blockIdx := range c.rawBlocks {
		block := model2.OriginBlock{
			Height:     uint64(c.rawBlocks[blockIdx].Number),
			BlockHash:  c.rawBlocks[blockIdx].Hash,
			ParentHash: c.rawBlocks[blockIdx].ParentHash,
			GasLimit:   uint64(c.rawBlocks[blockIdx].GasLimit),
			GasUsed:    uint64(c.rawBlocks[blockIdx].GasUsed),
			Timestamp:  time.Unix(int64(c.rawBlocks[blockIdx].Timestamp), 0),
			CreatedAt:  time.Now().UTC(),
			Done:       false,
			Txs:        make([]model2.OriginTx, 0),
			RawBlock:   &c.rawBlocks[blockIdx],
		}
		for _, rawTx := range c.rawBlocks[blockIdx].Transactions {
			methodName, methodArgs := c.jieData.JieMethod(rawTx.GetMethodID(), rawTx.GetInputBytes())
			tx := model2.OriginTx{
				TxHash:          rawTx.Hash,
				TxIdx:           uint64(rawTx.TransactionIndex),
				Height:          uint64(c.rawBlocks[blockIdx].Number),
				From:            rawTx.From,
				To:              rawTx.To,
				Value:           hex.BigInt(rawTx.Value),
				Status:          uint64(rawTx.Receipt.Status),
				MethodID:        rawTx.GetMethodID(),
				MethodName:      []byte(methodName),
				MethodArgs:      methodArgs.ToHex(),
				ContractAddress: rawTx.Receipt.ContractAddress,
				GasPrice:        hex.BigInt(rawTx.GasPrice),
				GasLimit:        hex.BigInt(rawTx.Gas),
				GasUsed:         hex.BigInt(rawTx.Receipt.GasUsed),
				Nonce:           uint64(rawTx.Nonce),
				Timestamp:       time.Unix(int64(c.rawBlocks[blockIdx].Timestamp), 0),
				ReceiptLogs:     make([]model2.OriginReceiptLog, 0),
				Input:           rawTx.Input,
			}
			for _, rawLog := range rawTx.Receipt.Logs {
				tx.ReceiptLogs = append(tx.ReceiptLogs, model2.OriginReceiptLog{
					Height:       uint64(c.rawBlocks[blockIdx].Number),
					TxHash:       rawTx.Hash,
					LogIndex:     uint64(rawLog.LogIndex),
					Address:      rawLog.Address,
					EventSig:     rawLog.GetTopic0(),
					EventName:    nil,
					EventPayload: nil,
				})
			}
			tx.RawTxReceipt = rawTx.Receipt
			block.Txs = append(block.Txs, tx)
		}
		events = append(events, block_listener.NewEvent(&block))
	}
	return events
}
