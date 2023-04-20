package alg

import (
	"testing"

	"github.com/721tools/backend-go/indexer/pkg/utils/alg"
	"github.com/stretchr/testify/assert"
)

func TestMethodSig(t *testing.T) {
	sig := MethodSig("deposit(uint256,uint256)")
	assert.Equal(t, "0xe2bbb158", sig.ToString())
}

func TestEventSig(t *testing.T) {
	assert.Equal(t, "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", TransferEventSigID.ToString())
}

func TestEventSig2(t *testing.T) {
	assert.Equal(t, "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925", ApprovalEventSigID.ToString())
}

func TestMethodSig2(t *testing.T) {
	//nameMethods   := []hex.Hex{alg.Keccak256("name()"), alg.Keccak256("NAME()"), alg.Keccak256("GetName()")}
	name := alg.Keccak256("name()")
	t.Log(name.NoPrefixHex())
}
