package handle

import (
	"context"
	"encoding/json"
	"time"

	"github.com/721tools/backend-go/index/internal/service"
	"github.com/osamingo/jsonrpc/v2"
)

type HeightReq struct {
	Heights []uint64 `json:"heights"`
}

type TimeStampResp struct {
	Height    uint64    `json:"height"`
	Timestamp time.Time `json:"timestamp"`
}

func Ts(c context.Context, rawReq *json.RawMessage) (interface{}, *jsonrpc.Error) {
	var heightReq *HeightReq
	if err := json.Unmarshal(*rawReq, &heightReq); err != nil || len(heightReq.Heights) > 100 {
		return nil, &jsonrpc.Error{Code: jsonrpc.ErrorCodeInvalidRequest, Message: "req invalid", Data: nil}
	}

	resp := []TimeStampResp{}
	db := service.NewBlockService()
	blocks := db.GetTimeStampByHeights(context.Background(), heightReq.Heights)

	log.Info("debug rpc handle get ts api", "heightreq", heightReq, "blocks", blocks)

	for _, ix := range blocks {
		res := &TimeStampResp{}
		res.Height = ix.Height
		res.Timestamp = ix.Timestamp
		resp = append(resp, *res)
	}
	return resp, nil
}
