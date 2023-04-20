package alg

import (
	"github.com/721tools/backend-go/indexer/pkg/utils/hex"
	"github.com/ethereum/go-ethereum/common/hexutil"
	sha32 "golang.org/x/crypto/sha3"
)

type SigID hex.Hex

func NewSigID(sigStr string) SigID {
	var h hex.Hex
	_ = h.ToHex(sigStr)
	return SigID(h)
}

func (s SigID) ToString() string {
	return hexutil.Encode(s)
}

func (s SigID) ToHex() hex.Hex {
	return hex.Hex(s)
}

func (s SigID) NoOxPrefix() string {
	return hex.Hex(s).NoPrefixHex()
}

// MethodSig
// example:
// methodAbi: deposit(uint256,uint256)
// return: 0xe2bbb158
func MethodSig(methodAbi string) SigID {
	sha3 := sha32.NewLegacyKeccak256()
	sha3.Write([]byte(methodAbi))
	return sha3.Sum(nil)[:4]
}

// EventSig
// example:
// eventAbi: Transfer(address,address,uint256)
// return: 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef
func EventSig(eventAbi string) SigID {
	sha3 := sha32.NewLegacyKeccak256()
	sha3.Write([]byte(eventAbi))
	return sha3.Sum(nil)
}
