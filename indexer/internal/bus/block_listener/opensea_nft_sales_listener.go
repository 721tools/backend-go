package block_listener

import (
	"context"
	"fmt"
	"math/big"

	model2 "github.com/721tools/backend-go/indexer/internal/model"
	"github.com/721tools/backend-go/indexer/pkg/blockchain/abi_parse"
	"github.com/721tools/backend-go/indexer/pkg/consts"
	"github.com/721tools/backend-go/indexer/pkg/db"
	"github.com/721tools/backend-go/indexer/pkg/mq"
	"xorm.io/xorm"
)

type NFTSalesListener struct {
	jieData *abi_parse.JieData
	xorm    *xorm.Engine
	mq      *mq.RedisMQ
}

func NewNFTSalesListener() *NFTSalesListener {
	return &NFTSalesListener{
		jieData: abi_parse.NewJieData(),
		xorm:    db.GetDBEngine(),
		mq:      mq.GetMQ(),
	}
}

func (t *NFTSalesListener) Handle(event *Event) (err error) {

	ctx := context.Background()

	nft_sales := make([]model2.NFTSales, 0)

	for _, tx := range event.Block.Txs {
		//		log.Info("debug tx", "tx is", tx.TxHash)
		ts := event.Block.Timestamp
		for _, rawLog := range tx.RawTxReceipt.Logs {

			logName, logArgs := t.jieData.JieEventLogsWithTag("SeaPort", rawLog.GetTopic0(), rawLog.GetLogBytes())
			if logName != "" {

				nfts := logArgs["nfts"].([]abi_parse.NFTItem)

				var isBundle uint
				nft_len := len(nfts)
				if nft_len > 1 {
					isBundle = 1
				} else {
					isBundle = 0
				}

				if nft_len < 1 {
					continue
				}
				var priceETH big.Int
				priceETH.Div(logArgs["volumeETH"].(*big.Int), big.NewInt(int64(nft_len)))

				for _, nft := range nfts {
					sale := model2.NFTSales{
						Height:    uint64(rawLog.BlockNumber),
						Timestamp: ts,
						TxHash:    rawLog.TransactionHash.HexStr(),
						LogIndex:  uint(rawLog.LogIndex),
						Plateform: uint(consts.OPENSEA),
						IsBundle:  isBundle,
						Address:   nft.Token,
						TokenId:   nft.Identifier.String(),
						Amount:    nft.Amount.String(),
						Buyer:     logArgs["buyer"].(string),
						Seller:    logArgs["seller"].(string),
						PriceETH:  priceETH.String(),
						Direction: uint(logArgs["direction"].(int)),
					}
					nft_sales = append(nft_sales, sale)
				}
			}
		}
	}

	// nft sales
	if len(nft_sales) == 0 {
		return nil
	}

	for _, n := range nft_sales {
		id := fmt.Sprintf("%s/%s/%d", n.Address, n.TokenId, consts.NFT_SALE)
		t.mq.Publish(ctx, consts.NFT_BEHAVIOR, id)
	}

	session := t.xorm.NewSession()
	err = session.Begin()
	if err = session.Begin(); err != nil {
		log.Error("start session failed", "err", err)
		return err
	}
	defer session.Close()
	explode := explodeG(nft_sales, nftSalesChunkSize)
	for _, arr := range explode {
		if _, err = session.InsertMulti(&arr); err != nil {
			return err
		}
	}
	return session.Commit()
}
