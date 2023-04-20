package abi_parse

import (
	"context"
	"testing"

	"github.com/721tools/backend-go/indexer/pkg/blockchain/client"
	"github.com/721tools/backend-go/indexer/pkg/utils/hex"
	"github.com/stretchr/testify/assert"
)

func TestNewERC721_TransferFrom(t *testing.T) {
	ethClient := client.NewClient("https://mainnet-eth.token.im", client.EvmClient)
	tx, err := ethClient.TxByHash(context.Background(), hex.HexstrToHex("0x82d38a4b63234e69b5e642ffd9f62bba446764593d6fbde6ac49dcac54c72142"))
	assert.NoError(t, err)
	t.Log("tx", tx)
	erc721 := NewERC721()
	methodsHandle := erc721.Methods()
	eventsHandle := erc721.EventLogs()
	handle, ok := methodsHandle[tx.GetMethodID().HexStr()]
	assert.True(t, ok)
	methodName, methodArgs := handle(tx.GetInputBytes())
	assert.Equal(t, "transferFrom", methodName)
	t.Log("method_args", methodArgs)
	for _, rawLog := range tx.Receipt.Logs {
		topic0 := rawLog.GetTopic0()
		handle, ok = eventsHandle[topic0.HexStr()]
		assert.True(t, ok)
		eventName, eventArgs := handle(rawLog.GetLogBytes())
		t.Log("event_name", eventName)
		t.Log("event_args", eventArgs)
		toHex := eventArgs.ToHex()
		t.Log("to Hex", toHex.HexStr())
		args := HexToArgs(toHex)
		t.Log("un_event_args", args)
		if eventName == "Transfer" {
			hexstrToHex := hex.HexstrToHex(args["amount"].(string))
			t.Log("sss", hexstrToHex.HexStr())

		}
	}
}

func TestNewERC721_SafeTransferFrom(t *testing.T) {
	ethClient := client.NewClient("https://main-rpc.linkpool.io/", client.EvmClient)
	tx, err := ethClient.TxByHash(context.Background(), hex.HexstrToHex("0xa1044d33e8e931451e514267b668c8658cdb1625ecda2eada73140889805a514"))
	assert.NoError(t, err)
	t.Log("tx", tx)
	erc721 := NewERC721()
	methodsHandle := erc721.Methods()
	eventsHandle := erc721.EventLogs()
	handle, ok := methodsHandle[tx.GetMethodID().HexStr()]
	assert.True(t, ok)
	methodName, methodArgs := handle(tx.GetInputBytes())
	assert.Equal(t, "safeTransferFrom", methodName)
	t.Log("method_args", methodArgs)
	for _, rawLog := range tx.Receipt.Logs {
		topic0 := rawLog.GetTopic0()
		handle, ok := eventsHandle[topic0.HexStr()]
		assert.True(t, ok)
		eventName, eventArgs := handle(rawLog.GetLogBytes())
		t.Log("event_name", eventName)
		t.Log("event_args", eventArgs)
	}
}

func TestNewERC721_Mint(t *testing.T) {
	ethClient := client.NewClient("https://main-rpc.linkpool.io/", client.EvmClient)
	tx, err := ethClient.TxByHash(context.Background(), hex.HexstrToHex("0x0db970df33b959aac2fec1f03ae48458ca499bb4a0dd60107edc017e046484c5"))
	assert.NoError(t, err)
	t.Log("tx", tx)
	erc721 := NewERC721()
	methodsHandle := erc721.Methods()
	eventsHandle := erc721.EventLogs()
	handle, ok := methodsHandle[tx.GetMethodID().HexStr()]
	assert.True(t, ok)
	methodName, methodArgs := handle(tx.GetInputBytes())
	assert.Equal(t, "mint", methodName)
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
