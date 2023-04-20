package model

import (
	"github.com/721tools/backend-go/index/pkg/utils/hex"
)

type RawInternalTx struct {
	From  hex.Hex
	To    hex.Hex
	Value uint64
}
