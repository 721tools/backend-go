package consts

import "github.com/721tools/backend-go/index/pkg/utils/hex"

var (
	NativeToken            = hex.HexstrToHex("0x0000000000000000000000000000000000000000")
	NativeTokenType        = 1
	NativeTokenName        = "ETH"
	NativeTokenSymbol      = "ETH"
	NativeTokenDecimal     = 18
	NativeTokenTotalSupply = hex.IntstrToBigInt("100000000")
)
