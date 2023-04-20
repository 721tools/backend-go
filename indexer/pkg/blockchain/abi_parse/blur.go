/**

Blur Market: https://etherscan.io/address/0x031aa05da8bf778dfc36d8d25ca68cbb2fc447c6#code

Event Log: https://etherscan.io/tx/0xe6c57fc9d9683fe7dc40036a0bba4d970bf0fe91fc3a2a087eb6e735018854b1#eventlog

// SPDX-License-Identifier: MIT
pragma solidity 0.8.17;

enum Side { Buy, Sell }
enum SignatureVersion { Single, Bulk }
enum AssetType { ERC721, ERC1155 }

struct Fee {
    uint16 rate;
    address payable recipient;
}

struct Order {
    address trader;
    Side side;
    address matchingPolicy;
    address collection;
    uint256 tokenId;
    uint256 amount;
    address paymentToken;
    uint256 price;
    uint256 listingTime;
    uint256 expirationTime;
    Fee[] fees;
    uint256 salt;
    bytes extraParams;
}

struct Input {
    Order order;
    uint8 v;
    bytes32 r;
    bytes32 s;
    bytes extraSignature;
    SignatureVersion signatureVersion;
    uint256 blockNumber;
}

emit OrdersMatched(
		buy.order.trader,
		sell.order.trader,
		sell.order,
		sellHash,
		buy.order,
		buyHash
);

tx:
https://etherscan.io/tx/0xe167bf9d1412fd762e5f87f4922eda5cc52aac659ffcd410ddcef188a0a05c89#eventlog

event log:
Address
0x000000000000ad05ccc4f10045630fb830b95127
Topics
0 0x61cbb2a3dee0b6064c2e681aadd61677fb4ef319f0b547508d495626f5a62f64
1 0x000000000000000000000000d93ab3b614e3de617032d2fd953db8f3c5bbd26e - seller address
2 0x00000000000000000000000039da41747a83aee658334415666f3ef92dd0d541 - buyer address
Data
0000000000000000000000000000000000000000000000000000000000000080 - 126/16=8
596c2193805d233c59b153644ebaf6793716878b085b62f6e153bd9609a6f4b6 - hash
00000000000000000000000000000000000000000000000000000000000002a0 - 672/16=42
75d2f8bef01fb92e1a146000580f1b77eab09e0d735e46b5749bfb753a254dc9 - hash
000000000000000000000000d93ab3b614e3de617032d2fd953db8f3c5bbd26e - sell.order.trader
0000000000000000000000000000000000000000000000000000000000000001 - side
00000000000000000000000000000000006411739da1c40b106f8511de5d1fac
0000000000000000000000006f9eb87f5a5638a3424c68ffae824608671f4ea6 - collection
000000000000000000000000000000000000000000000000000000000000061c - tokenid
0000000000000000000000000000000000000000000000000000000000000001 - amount
0000000000000000000000000000000000000000000000000000000000000000 - payment token
0000000000000000000000000000000000000000000000000049e57d63540000 - price
000000000000000000000000000000000000000000000000000000006351e9e8 - listingTime
00000000000000000000000000000000000000000000000000000000635b2468 - expirationTime
00000000000000000000000000000000000000000000000000000000000001a0 - fees[] 0x1a0 = 416
00000000000000000000000000000000e10f73da65326b78d105044b62284ab3
0000000000000000000000000000000000000000000000000000000000000200
0000000000000000000000000000000000000000000000000000000000000001
0000000000000000000000000000000000000000000000000000000000000320
000000000000000000000000aa84dd76a799c33880d61855650f0dd4f5f17b3c - recipient
0000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000039da41747a83aee658334415666f3ef92dd0d541 - buy.order.trader
0000000000000000000000000000000000000000000000000000000000000000 - side
00000000000000000000000000000000006411739da1c40b106f8511de5d1fac -
0000000000000000000000006f9eb87f5a5638a3424c68ffae824608671f4ea6 - collection
000000000000000000000000000000000000000000000000000000000000061c - tokenid
0000000000000000000000000000000000000000000000000000000000000001 - amount
0000000000000000000000000000000000000000000000000000000000000000 - payment token
0000000000000000000000000000000000000000000000000049e57d63540000 - price
000000000000000000000000000000000000000000000000000000006352799e - listingTime
00000000000000000000000000000000000000000000000000000000635295be - expirationTime
00000000000000000000000000000000000000000000000000000000000001a0 - fees[] 0x1a0 = 416
0000000000000000000000000000000099aaf67f2433d3f82274ee5f7a73b2dd
00000000000000000000000000000000000000000000000000000000000001c0
0000000000000000000000000000000000000000000000000000000000000000
0000000000000000000000000000000000000000000000000000000000000000


tx: https://etherscan.io/tx/0x0f44e24dc33f83caaa626bca090abef0bdb30298d53b460487a830149741a0ec#eventlog
Address
0x000000000000ad05ccc4f10045630fb830b95127
Topics
0 0x61cbb2a3dee0b6064c2e681aadd61677fb4ef319f0b547508d495626f5a62f64
1 0x000000000000000000000000dd2dadbd2d2512bebdb380f2ca02f61373d13ff3 - seller
2 0x000000000000000000000000f8799cdbaa97b050eca2996b4cdeb64169ad7dff - buyer
Data
0000000000000000000000000000000000000000000000000000000000000080 - 128/32=4
47aaf1dddca685b76006724e7eccf91dd30995c1d0394de38f04df449c2c8641 - hash
0000000000000000000000000000000000000000000000000000000000000260 - 608/32=19
56df5c39ebff9dd3b5ccd9f3b9be7b91c1ab13ce8a657b6308f77cdab9a36b6c - hash
000000000000000000000000dd2dadbd2d2512bebdb380f2ca02f61373d13ff3 - sell.order.trader
0000000000000000000000000000000000000000000000000000000000000001 - side
00000000000000000000000000000000006411739da1c40b106f8511de5d1fac
00000000000000000000000033c6eec1723b12c46732f7ab41398de45641fa42 - collection
0000000000000000000000000000000000000000000000000000000000001848 - tokenid
0000000000000000000000000000000000000000000000000000000000000001 - amount
0000000000000000000000000000000000000000000000000000000000000000 - payment token
00000000000000000000000000000000000000000000000003dba786ef8f0000 - price
0000000000000000000000000000000000000000000000000000000063528b09 - listingTime
000000000000000000000000000000000000000000000000000000006353dc8a
00000000000000000000000000000000000000000000000000000000000001a0 - fees[] 0x1a0 = 416
000000000000000000000000000000001b8489855cec409ffe17d1b484b5309d
00000000000000000000000000000000000000000000000000000000000001c0
0000000000000000000000000000000000000000000000000000000000000000
0000000000000000000000000000000000000000000000000000000000000000
000000000000000000000000f8799cdbaa97b050eca2996b4cdeb64169ad7dff - buy.order.trader
0000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000006411739da1c40b106f8511de5d1fac
00000000000000000000000033c6eec1723b12c46732f7ab41398de45641fa42
0000000000000000000000000000000000000000000000000000000000001848
0000000000000000000000000000000000000000000000000000000000000001
0000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000003dba786ef8f0000
000000000000000000000000000000000000000000000000000000006352a29f
000000000000000000000000000000000000000000000000000000006352bebf
00000000000000000000000000000000000000000000000000000000000001a0
00000000000000000000000000000000f0cacba1c8e7cd22d7955d0c4552f9e9
00000000000000000000000000000000000000000000000000000000000001c0
0000000000000000000000000000000000000000000000000000000000000000
0000000000000000000000000000000000000000000000000000000000000000

**/

package abi_parse

import (
	"github.com/721tools/backend-go/indexer/pkg/blockchain/alg"
	"github.com/721tools/backend-go/indexer/pkg/utils/hex"
)

type Blur struct {
}

func NewBlur() *Blur {
	return &Blur{}
}

func (s *Blur) Address() []hex.Hex {
	return make([]hex.Hex, 0)
}

func (s *Blur) Tag() string {
	return "Blur"
}

func (s *Blur) Methods() map[string]JieHandle {
	return nil
}

func (s *Blur) EventLogs() map[string]JieHandle {
	return map[string]JieHandle{
		alg.OrdersMatchedSigId.ToString(): s.ordersMatchedEvent,
	}
}

func (s *Blur) ordersMatchedEvent(hex hex.Hex) (string, Args) {
	params := dealWithEventLogInput(hex)
	params_len := len(params)

	seller_offset := params[2].ToDec()/32 + 2
	buyer_offset := params[4].ToDec()/32 + 2

	log.Info("debug ordersMatchedEvent params parse", "data", params, "len", params_len)

	if params_len != 36 && params_len != 38 {
		log.Warn("ordersMatchedEvent event parse error", "data", params, "len", params_len)
		return "", nil
	}

	args := make(Args)
	args["plateform"] = 1 // blur

	seller_list_time := params[seller_offset+8].ToDec()
	buyer_list_time := params[buyer_offset+8].ToDec()
	log.Info("debug seller list time", "seller_list_time", seller_list_time, "buyer_list_time", buyer_list_time)
	args["seller"] = params[0].ToAddress()
	args["buyer"] = params[1].ToAddress()

	direction := 0 // buy
	if params[seller_offset+1].ToDec() == 0 {
		direction = 1
	}
	args["direction"] = direction
	nfts := make([]NFTItem, 0)

	nft := &NFTItem{
		Token:      params[buyer_offset+3].ToAddress(),
		Identifier: params[buyer_offset+4].ToBigInt(),
		Amount:     params[buyer_offset+5].ToBigInt(),
	}

	args["nfts"] = append(nfts, *nft)
	args["volumeETH"] = params[13].ToBigInt()

	log.Info("debug volume", "args volumeETH", args["volumeETH"], "args", args)
	return "OrdersMatched", args
}
