package abi_parse

import (
	"context"
	"testing"

	"github.com/721tools/backend-go/indexer/pkg/blockchain/client"
	"github.com/721tools/backend-go/indexer/pkg/utils/hex"
	"github.com/stretchr/testify/assert"
)

func TestNewERC20_Transfer(t *testing.T) {

	ethClient := client.NewClient("https://main-rpc.linkpool.io/", client.EvmClient)
	tx, err := ethClient.TxByHash(context.Background(), hex.HexstrToHex("0x8bdcb88022f4a27b35418c4a664f2e94cf0ba2b12c89da072e923f08e1ce5426"))
	assert.NoError(t, err)

	t.Log("tx", tx)
	erc20 := NewERC20()
	methodsHandle := erc20.Methods()
	eventsHandle := erc20.EventLogs()
	handle, ok := methodsHandle[tx.GetMethodID().HexStr()]
	assert.True(t, ok)
	methodName, methodArgs := handle(tx.GetInputBytes())
	assert.Equal(t, "transferFrom", methodName)
	t.Log("method_args", methodArgs)
	for _, rawLog := range tx.Receipt.Logs {
		topic0 := rawLog.GetTopic0()
		handle, ok = eventsHandle[topic0.HexStr()]
		if !ok {
			continue
		}
		assert.True(t, ok)
		eventName, eventArgs := handle(rawLog.GetLogBytes())
		t.Log("event_name", eventName)
		t.Log("event_args", eventArgs)
	}
}
