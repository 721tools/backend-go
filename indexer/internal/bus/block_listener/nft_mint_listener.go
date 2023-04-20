package block_listener

import (
	"context"
	"fmt"

	model2 "github.com/721tools/backend-go/index/internal/model"
	"github.com/721tools/backend-go/index/internal/service"
	"github.com/721tools/backend-go/index/pkg/blockchain/abi_parse"
	"github.com/721tools/backend-go/index/pkg/blockchain/client"
	"github.com/721tools/backend-go/index/pkg/consts"
	"github.com/721tools/backend-go/index/pkg/db"
	"github.com/721tools/backend-go/index/pkg/mq"
	"github.com/721tools/backend-go/index/pkg/utils/hex"
	"xorm.io/xorm"
)

type NFTMintListener struct {
	svc     service.ContractIface
	client  client.Client
	jieData *abi_parse.JieData
	xorm    *xorm.Engine
	mq      *mq.RedisMQ
}

func NewNFTMintListener() *NFTMintListener {
	return &NFTMintListener{
		svc:     service.NewContract(),
		client:  client.GetClient(),
		jieData: abi_parse.NewJieData(),
		xorm:    db.GetDBEngine(),
		mq:      mq.GetMQ(),
	}
}

func (t *NFTMintListener) Handle(event *Event) (err error) {
	ctx := context.Background()
	nftMintFlows := make([]model2.NFTMint, 0)
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
			contractType := t.svc.GetContractType(ctx, rawLog.Address)
			if contractType.IsERC721() {
				logName, logArgs := t.jieData.JieEventLogsWithTag("ERC721", rawLog.GetTopic0(), rawLog.GetLogBytes())
				// 如果存在token flow, 处理对应数据
				oLog.EventName = []byte(logName)
				oLog.EventPayload = logArgs.ToHex()
				if rawLog.IsTokenFlow() {
					flow := oLog.ToTokenFlow()
					// nft_mint
					if flow.IsMint() {
						token := flow.ToMintNFTToken(tx.From)
						token.TokenURL, _ = t.client.TokenURL(ctx, token.Address, token.TokenID)
						if len(token.TokenURL) > 128 {
							token.TokenURL = ""
						}
						nftMintFlows = append(nftMintFlows, token)
					}
				}
			}
		}
	}

	// sync db
	if len(nftMintFlows) == 0 {
		return nil
	}

	session := t.xorm.NewSession()
	err = session.Begin()
	if err = session.Begin(); err != nil {
		log.Error("start session failed", "err", err)
		return err
	}
	defer session.Close()

	explode2 := explodeG(nftMintFlows, nftMintChunkSize)
	for _, arr := range explode2 {
		if _, err = session.InsertMulti(&arr); err != nil {
			log.Error("insert nft mint occur a error", "err", err, "arr", arr)
			return err
		}
	}

	for _, n := range nftMintFlows {
		id := fmt.Sprintf("%s/%s/%d", n.Address, hex.HexstrToBigInt(n.TokenID.HexStr()), consts.NFT_MINT)
		t.mq.Publish(ctx, consts.NFT_BEHAVIOR, id)
	}
	return session.Commit()
}
