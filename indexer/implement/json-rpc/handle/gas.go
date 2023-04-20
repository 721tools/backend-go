package handle

import (
	"context"
	"encoding/json"

	"github.com/721tools/backend-go/indexer/implement/bus/block_listener"
	"github.com/721tools/backend-go/indexer/pkg/consts"
	"github.com/721tools/backend-go/indexer/pkg/mq"
	"github.com/osamingo/jsonrpc/v2"
)

type GasNowResp struct {
	block_listener.GasPerBlock
}

func GASNOW(c context.Context, rawReq *json.RawMessage) (interface{}, *jsonrpc.Error) {
	resp := GasNowResp{}
	msg, err := mq.GetMQ().Get(context.TODO(), consts.GAS_PRICE_NOW_KEY)
	if err != nil {
		log.Warn("get from redis error", "msg", msg)
	}

	err = json.Unmarshal([]byte(msg), &resp)
	if err != nil {
		log.Warn("get gas error", "cant Unmarshal msg", msg)
	}

	return resp, nil
}
