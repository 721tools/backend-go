package alg

import (
	"github.com/721tools/backend-go/indexer/pkg/utils/hex"
	"golang.org/x/crypto/sha3"
)

func Keccak256(str string) hex.Hex {
	sha := sha3.NewLegacyKeccak256()
	sha.Write([]byte(str))
	return hex.Hex(sha.Sum(nil)[0:4])
}
