package handle

import (
	"context"
	"encoding/json"

	"github.com/osamingo/jsonrpc/v2"
)

// https://etherscan.io/address/0xb38827497daf7f28261910e33e22219de087c8f5#code
// 批量查询订单状态
// req: {
//  hash(order1),
//  ...
// }
// res:
// hash(order1) - true | false
func BatchCheckOrderStatus2(ctx context.Context, rawReq *json.RawMessage) (interface{}, *jsonrpc.Error) {
	return nil, &jsonrpc.Error{}
}

// 批量获取地址 nonce 值，用于判断订单是否可执行
// req: [
//	{address:0x112, nonce:0}
// ]
// res: [
//	{address:0x112, nonce:1}
// ]
func BatchGetAccountCounter2(ctx context.Context, rawReq *json.RawMessage) (interface{}, *jsonrpc.Error) {
	return nil, &jsonrpc.Error{}
}
