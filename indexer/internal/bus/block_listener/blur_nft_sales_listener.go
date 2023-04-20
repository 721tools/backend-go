package block_listener

import (
	"context"
	"fmt"
	"math/big"

	model2 "github.com/721tools/backend-go/index/internal/model"
	"github.com/721tools/backend-go/index/pkg/blockchain/abi_parse"
	"github.com/721tools/backend-go/index/pkg/consts"
	"github.com/721tools/backend-go/index/pkg/db"
	"github.com/721tools/backend-go/index/pkg/mq"
	"github.com/721tools/backend-go/index/pkg/utils/hex"
	"xorm.io/xorm"
)

type BlurSalesListener struct {
	jieData *abi_parse.JieData
	xorm    *xorm.Engine
	mq      *mq.RedisMQ
}

func NewBlurSalesListener() *BlurSalesListener {
	return &BlurSalesListener{
		jieData: abi_parse.NewJieData(),
		xorm:    db.GetDBEngine(),
		mq:      mq.GetMQ(),
	}
}

func (t *BlurSalesListener) Handle(event *Event) (err error) {
	ctx := context.Background()
	nft_sales := make([]model2.NFTSales, 0)

	for _, tx := range event.Block.Txs {
		log.Info("debug tx", "tx is", tx.TxHash)
		ts := event.Block.Timestamp
		for _, rawLog := range tx.RawTxReceipt.Logs {
			logName, logArgs := t.jieData.JieEventLogsWithTag("Blur", rawLog.GetTopic0(), rawLog.GetLogBytes())
			if logName != "" {
				nft_len := 1
				nfts := logArgs["nfts"].([]abi_parse.NFTItem)
				var priceETH big.Int
				volumeETH := big.Int(logArgs["volumeETH"].(hex.BigInt))
				priceETH.Div(&volumeETH, big.NewInt(int64(nft_len)))

				// blur swap
				// if seller address = 0x39da41747a83aeE658334415666f3EF92DD0D541
				// then buyer = eof and seller = params[1]
				buyer := logArgs["buyer"].(string)
				seller := logArgs["seller"].(string)
				log.Info("debug", "buyer", buyer, "seller", seller)
				if seller == "0x39da41747a83aee658334415666f3ef92dd0d541" { // swap mode
					seller = buyer
					buyer = tx.From.String()
				} else if buyer == "0x39da41747a83aee658334415666f3ef92dd0d541" {
					buyer = tx.From.String()
				}
				log.Info("debug", "buyer", buyer, "seller", seller)
				for _, nft := range nfts {
					sale := model2.NFTSales{
						Height:    uint64(rawLog.BlockNumber),
						Timestamp: ts,
						TxHash:    rawLog.TransactionHash.HexStr(),
						LogIndex:  uint(rawLog.LogIndex),
						Plateform: uint(consts.Blur),
						IsBundle:  0,
						Address:   nft.Token,
						TokenId:   nft.Identifier.String(),
						Amount:    nft.Amount.String(),
						Buyer:     buyer,
						Seller:    seller,
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
