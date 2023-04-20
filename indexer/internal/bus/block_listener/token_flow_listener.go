package block_listener

import (
	"context"

	model2 "github.com/721tools/backend-go/indexer/internal/model"
	"github.com/721tools/backend-go/indexer/internal/service"
	"github.com/721tools/backend-go/indexer/pkg/blockchain/abi_parse"
	"github.com/721tools/backend-go/indexer/pkg/consts"
	"github.com/721tools/backend-go/indexer/pkg/db"
	"github.com/721tools/backend-go/indexer/pkg/utils/hex"
	"xorm.io/xorm"
)

type TokenFlowListener struct {
	svc     service.ContractIface
	jieData *abi_parse.JieData
	xorm    *xorm.Engine
}

func NewTokenFlowListener() *TokenFlowListener {
	return &TokenFlowListener{svc: service.NewContract(), jieData: abi_parse.NewJieData(), xorm: db.GetDBEngine()}
}

func (t *TokenFlowListener) Handle(event *Event) (err error) {

	tokenFlows := make([]model2.TokenFlow, 0)
	for _, tx := range event.Block.Txs {
		for _, rawLog := range tx.RawTxReceipt.Logs {
			var eventSig hex.Hex
			eventSig = rawLog.GetTopic0()
			oLog := &model2.OriginReceiptLog{
				Height:       uint64(tx.Height),
				TxHash:       tx.TxHash,
				LogIndex:     uint64(rawLog.LogIndex),
				Address:      rawLog.Address,
				EventSig:     eventSig,
				EventName:    nil,
				EventPayload: nil,
			}
			contractType := t.svc.GetContractType(context.Background(), rawLog.Address)
			if contractType.IsERC721() {
				logName, logArgs := t.jieData.JieEventLogsWithTag("ERC721", rawLog.GetTopic0(), rawLog.GetLogBytes())
				updateLog(tx.TxHash, rawLog.LogIndex, []byte(logName), logArgs.ToHex())
				// 如果存在token flow, 处理对应数据
				oLog.EventName = []byte(logName)
				oLog.EventPayload = logArgs.ToHex()
				if rawLog.IsTokenFlow() {
					flow := oLog.ToTokenFlow()
					tokenFlows = append(tokenFlows, flow)
				}
			} else if contractType.IsERC20() {
				logName, logArgs := t.jieData.JieEventLogsWithTag("ERC20", rawLog.GetTopic0(), rawLog.GetLogBytes())
				updateLog(tx.TxHash, rawLog.LogIndex, []byte(logName), logArgs.ToHex())
				oLog.EventName = []byte(logName)
				oLog.EventPayload = logArgs.ToHex()
				// 如果存在token flow, 处理对应数据
				if rawLog.IsTokenFlow() {
					flow := oLog.ToTokenFlow()
					tokenFlows = append(tokenFlows, flow)
				}
			}
		}

		if tx.Value.GreaterZero() && tx.Input.EqualTo(hex.HexstrToHex("0x")) {
			tokenFlows = append(tokenFlows, model2.TokenFlow{
				Height:   tx.Height,
				TxHash:   tx.TxHash,
				LogIndex: -1,
				Address:  consts.NativeToken,
				From:     tx.From,
				To:       tx.To,
				Value:    tx.Value.Bytes(),
				Amount:   tx.Value.Bytes(),
			})
		}
	}

	// sync db
	if len(tokenFlows) == 0 {
		return nil
	}

	session := t.xorm.NewSession()
	err = session.Begin()
	if err = session.Begin(); err != nil {
		log.Error("start session failed", "err", err)
		return err
	}
	defer session.Close()
	explode := explodeG(tokenFlows, tokenFlowsChunkSize)
	for _, arr := range explode {
		if _, err = session.InsertMulti(&arr); err != nil {
			log.Error("insert token flow a error", "err", err)
			return err
		}
	}

	return session.Commit()
}
