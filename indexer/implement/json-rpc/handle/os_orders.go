package handle

import (
	"context"
	"encoding/json"

	"github.com/osamingo/jsonrpc/v2"
)

// https://etherscan.io/address/0x00000000006c3852cbef3e08e8df289169ede581#readContract
// 批量查询订单状态
// req: {
//  hash(order1),
//  ...
// }
// res:
// hash(order1) - true | false
func BatchCheckOrderStatus(ctx context.Context, rawReq *json.RawMessage) (interface{}, *jsonrpc.Error) {
	return nil, &jsonrpc.Error{}
}

// 批量获取地址 nonce 值，用于判断订单是否可执行
// req: [
//	{address:0x112, nonce:0}
// ]
// res: [
//	{address:0x112, nonce:1}
// ]
func BatchGetAccountCounter(ctx context.Context, rawReq *json.RawMessage) (interface{}, *jsonrpc.Error) {
	return nil, &jsonrpc.Error{}
}
