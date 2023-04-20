package abi_parse

import (
	"github.com/721tools/backend-go/indexer/pkg/utils/hex"
	"github.com/721tools/backend-go/indexer/pkg/utils/log16"
)

type JieHandle func(hex hex.Hex) (string, Args)

type Parse interface {
	Tag() string
	Address() []hex.Hex
	Methods() map[string]JieHandle
	EventLogs() map[string]JieHandle
}

type NFTItem struct {
	Token      string
	Identifier hex.BigInt
	Amount     hex.BigInt
}

var log = log16.NewLogger("module", "abi_parse")
