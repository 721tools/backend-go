package hex

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type EthParamHex []byte

func (e EthParamHex) String() string {
	return hexutil.Encode(e)
}

func (e EthParamHex) Hex() string {
	return e.String()
}

func (e EthParamHex) Bytes() []byte {
	return e
}

func (e *EthParamHex) ToHex(value any) (err error) {
	switch v := value.(type) {
	case []byte:
		*e = v
		return
	case string:
		if len(v) > 2 {
			*e, err = hexutil.Decode(v)
		}
		return
	default:
		return fmt.Errorf("can not convert %T to HexStr", value)
	}
}

func (e *EthParamHex) ToAddress() string {
	bytes := e.Bytes()
	hex := EthParamHex(bytes[len(bytes)-20:])
	return hex.Hex()
}

func (e *EthParamHex) ToHexInt() Hex {
	leftZeroes := common.TrimLeftZeroes(e.Bytes())
	if len(leftZeroes) != 0 {
		return leftZeroes
	}
	return IntToHex(0)
}

func (e *EthParamHex) ToBigInt() BigInt {
	return HexstrToBigInt(e.Hex())
}

func (e *EthParamHex) ToDec() int64 {
	zeroes := common.TrimLeftZeroes(e.Bytes())
	str := common.Bytes2Hex(zeroes)
	num, _ := strconv.ParseInt(str, 16, 64)
	return num
}

func (e *EthParamHex) ToBool() bool {
	zeroes := common.TrimLeftZeroes(e.Bytes())
	var t big.Int
	bytes := t.SetBytes(zeroes)
	return bytes.Int64() > 0
}
