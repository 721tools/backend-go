package opensealistener

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/721tools/backend-go/indexer/pkg/consts"
	"github.com/721tools/backend-go/indexer/pkg/mq"
	"github.com/721tools/stream-api-go/sdk"
)

func Init(ctx context.Context, chain, key string) {

	rdb := mq.GetMQ()

	chainId := sdk.MAIN_NET
	if chain == "goerli" {
		chainId = sdk.TEST_NET
	}

	ns := sdk.NewNotifyService(chainId, key)

	ns.Subscribe("*", sdk.ITEM_LISTED, func(msg *sdk.Message) error {

		payload := msg.Payload.Payload.(*sdk.ItemListedRes)
		fmt.Printf("ITEM_LISTED: msg.Payload.Payload: %v\n", payload.Item.Metadata.ImageUrl)

		if !payload.IsPrivate &&
			payload.PaymentToken.Symbol == "ETH" &&
			len(payload.BasePrice) >= 17 {
			if (chainId == sdk.MAIN_NET && payload.Item.Chain.Name == "ethereum") ||
				(chainId == sdk.TEST_NET && payload.Item.Chain.Name == chain) {
				item, _ := json.Marshal(msg)
				fmt.Printf("ITEM_LISTED: msg.Payload.Payload: %v\n", string(item))
				rdb.Publish(ctx, consts.OPENSEA_ETH_ORDER_LISTING, string(item))

				if id := to(payload.Item.NFTId, consts.NFT_LIST); id != "" {
					rdb.Publish(ctx, consts.NFT_BEHAVIOR, id)
				}
			}
		}

		return nil
	})

	ns.Subscribe("*", sdk.ITEM_RECEIVED_BID, func(msg *sdk.Message) error {
		payload := msg.Payload.Payload.(*sdk.ItemReceivedBidRes)
		//				fmt.Printf("ITEM_RECEIVED_BID: msg.Payload.Payload: %v\n", payload.BasePrice)
		if (chainId == sdk.MAIN_NET && payload.Item.Chain.Name == "ethereum") ||
			(chainId == sdk.TEST_NET && payload.Item.Chain.Name == chain) &&
				payload.PaymentToken.Symbol == "WETH" &&
				len(payload.BasePrice) >= 17 &&
				payload.Quantity == 1 {
			item, _ := json.Marshal(msg)
			fmt.Printf("ITEM_RECEIVED_BID: msg.Payload.Payload: %v\n", string(item))
			rdb.Publish(ctx, consts.OPENSEA_ETH_ORDER_OFFER, string(item))
		}
		return nil
	})

	ns.Subscribe("*", sdk.COLLECTION_OFFER, func(msg *sdk.Message) error {
		payload := msg.Payload.Payload.(*sdk.CollectionOfferRes)

		if payload.PaymentToken.Symbol == "WETH" && len(payload.BasePrice) >= 17 {
			item, _ := json.Marshal(msg)
			fmt.Printf("COLLECTION_OFFER: msg.Payload.Payload: %v\n", string(item))
			rdb.Publish(ctx, consts.OPENSEA_ETH_COLLECTION_OFFER, string(item))
		}
		return nil
	})

	ns.Subscribe("*", sdk.TRAIT_OFFER, func(msg *sdk.Message) error {
		payload := msg.Payload.Payload.(*sdk.TraitOfferRes)

		if payload.PaymentToken.Symbol == "WETH" && len(payload.BasePrice) >= 17 {
			item, _ := json.Marshal(msg)
			fmt.Printf("TRAIT_OFFER: msg.Payload.Payload: %v\n", string(item))
			rdb.Publish(ctx, consts.OPENSEA_ETH_COLLECTION_OFFER, string(item))
		}
		return nil
	})

	ns.Start()
}
