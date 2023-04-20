package abi_parse

import (
	"encoding/json"

	"github.com/721tools/backend-go/indexer/pkg/utils/hex"
)

type Args map[string]interface{}

func (a Args) ToHex() hex.Hex {
	marshal, _ := json.Marshal(a)
	return marshal
}

func HexToArgs(hex hex.Hex) Args {
	a := make(Args)
	_ = json.Unmarshal(hex, &a)
	return a
}
