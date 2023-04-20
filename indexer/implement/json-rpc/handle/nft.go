package handle

import (
	"context"
	"encoding/json"

	"github.com/721tools/backend-go/indexer/implement/hotness"
	"github.com/osamingo/jsonrpc/v2"
)

type Collection struct {
	//	Slug string `json: slug`
	TokenAddress string `json:"token_address,omitempty"`
}

type CollectionsResp struct {
	Collections []Collection `json:"collections,omitempty"`
}

func GetTheMostPopularCollections(c context.Context, rawReq *json.RawMessage) (interface{}, *jsonrpc.Error) {
	resp := CollectionsResp{}
	res := hotness.GetHotness()

	for _, address := range res {
		col := &Collection{TokenAddress: address}
		resp.Collections = append(resp.Collections, *col)
	}

	return resp, nil
}
