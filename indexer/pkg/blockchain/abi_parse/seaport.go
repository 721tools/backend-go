/**

https://etherscan.io/tx/0x08484c11c73f8bdab01f2f172eb5dece1f7601be9da3c96c072312fe275e6e39#eventlog

Address
0x00000000006c3852cbef3e08e8df289169ede581
Name
OrderFulfilled(bytes32,address,address,address,(uint8,address,uint256,uint256)[],(uint8,address,uint256,uint256,address)[])
Topics
0 0x9d9af8e38d66c62e2c12f0225249fd9d721c54b83f48d9352c97c6cacdcb6f31 - OrderFulfilled event signature
1 0x000000000000000000000000be0fe4781098762978853301ae7ac0e42de0113c - offerer
2 0x000000000000000000000000004c00500000ad104d7dbd00e3ae0a5c00560c00 - zone

Data
 *  - 0x00: 2e182f97b32a78309378f539abb276b8f7deaf501700a821f87f6893c590d901 - orderHash
 *  - 0x20: 0000000000000000000000008d8890235639aa0715afb141fd2943f19e101e78 - fulfiller
 *  - 0x40: 0000000000000000000000000000000000000000000000000000000000000080 - offer offset (0x80)
 *  - 0x60: 0000000000000000000000000000000000000000000000000000000000000120 - consideration offset (0x120)
 *  - 0x80: 0000000000000000000000000000000000000000000000000000000000000001 - offer.length (1)
 *  - 0xa0: 0000000000000000000000000000000000000000000000000000000000000002 - offerItemType
 *  - 0xc0: 00000000000000000000000013aecd04c22fda1cafa38dd83831feb34ce9d8f8 - offerToken
 *  - 0xe0: 0000000000000000000000000000000000000000000000000000000000000083 - offerIdentifier
 *  - 0x100: 0000000000000000000000000000000000000000000000000000000000000001- offerAmount
 *  - 0x120: 0000000000000000000000000000000000000000000000000000000000000003- consideration.length
0000000000000000000000000000000000000000000000000000000000000000 - considerationItemType[0]
0000000000000000000000000000000000000000000000000000000000000000 - considerationToken[0]
0000000000000000000000000000000000000000000000000000000000000000 - considerationIdentifier[0]
000000000000000000000000000000000000000000000000007265bab9c48000 - considerationAmount[0]
000000000000000000000000be0fe4781098762978853301ae7ac0e42de0113c - considerationRecipient[0]

0000000000000000000000000000000000000000000000000000000000000000 - considerationItemType[1]
0000000000000000000000000000000000000000000000000000000000000000
0000000000000000000000000000000000000000000000000000000000000000
000000000000000000000000000000000000000000000000000344bc31318000
0000000000000000000000008de9c5a032463c561423387a9648c5c7bcc5bc90

0000000000000000000000000000000000000000000000000000000000000000 - considerationItemType[2]
0000000000000000000000000000000000000000000000000000000000000000
0000000000000000000000000000000000000000000000000000000000000000
000000000000000000000000000000000000000000000000000d12f0c4c60000
000000000000000000000000c41d68b025713f85baa75f127f6911859b014d42

* OrderFulfilled event data structrue:
*   event OrderFulfilled(
*     bytes32 orderHash,
*     address indexed offerer,
*     address indexed zone,
*     address fulfiller,
*     SpentItem[] offer,
*       > (itemType, token, id, amount)
*     ReceivedItem[] consideration
*       > (itemType, token, id, amount, recipient)
*   )
* topic0 - OrderFulfilled event signature
* topic1 - offerer
* topic2 - zone
* data:
 *  - 0x00: orderHash
 *  - 0x20: fulfiller
 *  - 0x40: offer offset (0x80)
 *  - 0x60: consideration offset (0x120)
 *  - 0x80: offer.length (1)
 *  - 0xa0: offerItemType
 *  - 0xc0: offerToken
 *  - 0xe0: offerIdentifier
 *  - 0x100: offerAmount
 *  - 0x120: consideration.length (1 + additionalRecipients.length)
 *  - 0x140: considerationItemType
 *  - 0x160: considerationToken
 *  - 0x180: considerationIdentifier
 *  - 0x1a0: considerationAmount
 *  - 0x1c0: considerationRecipient
**/

package abi_parse

import (
	"math/big"

	"github.com/721tools/backend-go/indexer/pkg/blockchain/alg"
	"github.com/721tools/backend-go/indexer/pkg/utils/hex"
)

type SeaPort struct {
}

func NewSeaPort() *SeaPort {
	return &SeaPort{}
}

func (s *SeaPort) Address() []hex.Hex {
	return make([]hex.Hex, 0)
}

func (s *SeaPort) Tag() string {
	return "SeaPort"
}

func (s *SeaPort) Methods() map[string]JieHandle {
	return nil
}

func (s *SeaPort) EventLogs() map[string]JieHandle {
	return map[string]JieHandle{
		alg.OrderFulfilledSigId.ToString(): s.orderFulfilledEvent,
	}
}

type offerer struct {
	ItemType   int64
	Token      string
	Identifier hex.BigInt
	Amount     hex.BigInt
}

type consideration struct {
	offerer
	RationRecipient string
}

type money_item struct {
	Amount hex.BigInt
}

// Constants ported from Seaport contracts
// See https://github.com/ProjectOpenSea/seaport/blob/main/contracts/lib/ConsiderationEnums.sol#L116
type SeaportItemType int32

const (
	NATIVE SeaportItemType = iota // value -> 0
	ERC_20
	ERC_721
	ERC1155
	ERC721_WITH_CRITERIA
	ERC1155_WITH_CRITERIA
)

func isERC721Item(thisType SeaportItemType) bool {
	if thisType == ERC_721 {
		return true
	}
	return false
}

// todo: erc1155 support
func isNFTItem(thisType SeaportItemType) bool {
	if isERC721Item(thisType) == true {
		return true
	}

	return false
}

func isMoney(thisType SeaportItemType) bool {
	if thisType == NATIVE || thisType == ERC_20 {
		return true
	}

	return false
}

func getOrderETHVolume(orders []money_item) *big.Int {
	volumeETH := big.NewInt(0)
	for _, o := range orders {
		amount := big.NewInt(0).SetBytes(o.Amount.Bytes())
		volumeETH.Add(volumeETH, amount)
	}

	return volumeETH
}

// test case:
/*
- bid with weth erc1155
https://etherscan.io/tx/0x0fe99f91df7644952b4eb9a906878aa34cc91f65a0314b92abe3cee8e17901ca#eventlog
*/
func (s *SeaPort) orderFulfilledEvent(hex hex.Hex) (string, Args) {
	params := dealWithEventLogInput(hex)

	offer_offset := params[4].ToDec()/32 + 2         // 128/32=4
	consideration_offset := params[5].ToDec()/32 + 2 // 288/32=9

	offer_len := params[offer_offset].ToDec()
	consideration_len := params[consideration_offset].ToDec()
	params_len := len(params)

	//	log.Info("debug OrderFulfilled event params parse", "data", params, "len", len(params))
	// 4+2+offer_len*1*4+consideration_len*3*5
	if params_len < 1 || consideration_len < 2 ||
		int64(params_len) != (8+offer_len*4+consideration_len*5) {
		log.Warn("OrderFulfilled event parse error", "data", params)
		return "", nil
	}

	args := make(Args)
	args["buyer"] = ""
	args["seller"] = ""
	args["nfts"] = make([]NFTItem, 0)
	args["volumeETH"] = ""
	args["plateform"] = 0 // opensea

	offerer_items := make([]offerer, 0)
	consideration_items := make([]consideration, 0)

	// offerer
	for i := 0; i < int(offer_len); i++ {
		item := offerer{}

		start_index := int(offer_offset) + i*4
		item.ItemType = params[start_index+1].ToDec()
		item.Token = params[start_index+2].ToAddress()
		item.Identifier = params[start_index+3].ToBigInt()
		item.Amount = params[start_index+4].ToBigInt()

		offerer_items = append(offerer_items, item)
	}

	// consideration
	for i := 0; i < int(consideration_len); i++ {
		item := consideration{}

		start_index := int(consideration_offset) + i*5

		item.ItemType = params[start_index+1].ToDec()
		item.Token = params[start_index+2].ToAddress()
		item.Identifier = params[start_index+3].ToBigInt()
		item.Amount = params[start_index+4].ToBigInt()
		item.RationRecipient = params[start_index+5].ToAddress()

		consideration_items = append(consideration_items, item)
	}

	//	log.Info("debug offerer & consideration", "offerer", offerer_items, "consideration", consideration_items)

	if isMoney(SeaportItemType(offerer_items[0].ItemType)) {
		// this is a bid order, offerer = buyer
		args["buyer"] = params[0].ToAddress()
		args["seller"] = params[3].ToAddress()
		args["direction"] = 0

		money_items := make([]money_item, 0)
		for _, o := range offerer_items {
			if isMoney(SeaportItemType(o.ItemType)) {
				item := &money_item{}
				item.Amount = o.Amount
				money_items = append(money_items, *item)
			}
		}

		args["volumeETH"] = getOrderETHVolume(money_items)

		nft_items := make([]NFTItem, 0)
		for _, o := range consideration_items {
			if isERC721Item(SeaportItemType(o.ItemType)) {
				item := &NFTItem{
					Token:      o.Token,
					Identifier: o.Identifier,
					Amount:     o.Amount,
				}
				nft_items = append(nft_items, *item)
			}
		}

		args["nfts"] = nft_items

	} else {
		// this is a ask order, offerer_items = nft
		args["buyer"] = params[3].ToAddress()
		args["seller"] = params[0].ToAddress()
		args["direction"] = 1

		money_items := make([]money_item, 0)
		for _, o := range consideration_items {
			if isMoney(SeaportItemType(o.ItemType)) {
				item := &money_item{}
				item.Amount = o.Amount
				money_items = append(money_items, *item)
			}
		}

		args["volumeETH"] = getOrderETHVolume(money_items)

		nft_items := make([]NFTItem, 0)
		for _, o := range offerer_items {
			if isERC721Item(SeaportItemType(o.ItemType)) {
				item := &NFTItem{
					Token:      o.Token,
					Identifier: o.Identifier,
					Amount:     o.Amount,
				}
				nft_items = append(nft_items, *item)
			}
		}

		args["nfts"] = nft_items
	}

	//	log.Info("debug volume", "args volumeETH", args["volumeETH"], "args", args)

	return "OrderFulfilled", args
}
