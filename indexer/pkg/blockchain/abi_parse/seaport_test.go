package abi_parse

import (
	"context"
	"testing"

	"github.com/721tools/backend-go/index/pkg/blockchain/client"
	"github.com/721tools/backend-go/index/pkg/utils/hex"
	"github.com/stretchr/testify/assert"
)

// https://etherscan.io/tx/0x08484c11c73f8bdab01f2f172eb5dece1f7601be9da3c96c072312fe275e6e39#eventlog
func TestFulfillBasicOrder(t *testing.T) {

	ethClient := client.NewClient("wss://green-small-mountain.discover.quiknode.pro/70510f2528a6c57bcb376fda51073da9aceb214f/", client.EvmClient)
	tx, err := ethClient.TxByHash(context.Background(), hex.HexstrToHex("0x08484c11c73f8bdab01f2f172eb5dece1f7601be9da3c96c072312fe275e6e39"))
	assert.NoError(t, err)
	sea := NewSeaPort()
	eventsHandle := sea.EventLogs()

	for _, rawLog := range tx.Receipt.Logs {
		topic0 := rawLog.GetTopic0()
		handle, ok := eventsHandle[topic0.HexStr()]
		if !ok {
			continue
		}
		assert.True(t, ok)
		eventName, eventArgs := handle(rawLog.GetLogBytes())
		t.Log("event_name", eventName)
		t.Log("event_args", eventArgs)
	}
}

// https://etherscan.io/tx/0x058ae6ca60ed92f3c8197c5871468553b91ba6a1bd83fc39a4021e27cb7d1d94#eventlog
func TestFulfillAvailableAdvancedOrders(t *testing.T) {

	ethClient := client.NewClient("wss://green-small-mountain.discover.quiknode.pro/70510f2528a6c57bcb376fda51073da9aceb214f/", client.EvmClient)
	tx, err := ethClient.TxByHash(context.Background(), hex.HexstrToHex("0x058ae6ca60ed92f3c8197c5871468553b91ba6a1bd83fc39a4021e27cb7d1d94"))
	assert.NoError(t, err)
	sea := NewSeaPort()
	eventsHandle := sea.EventLogs()

	for _, rawLog := range tx.Receipt.Logs {
		topic0 := rawLog.GetTopic0()
		handle, ok := eventsHandle[topic0.HexStr()]
		if !ok {
			continue
		}
		assert.True(t, ok)
		eventName, eventArgs := handle(rawLog.GetLogBytes())
		t.Log("event_name", eventName)
		t.Log("event_args", eventArgs)
	}
}

// https://etherscan.io/tx/0x101239b9ef78582a3d628685f99f596266cff0c25e1c1e868d763428d0c18013#eventlog
func TestFulfillAdvancedOrder(t *testing.T) {

	ethClient := client.NewClient("wss://green-small-mountain.discover.quiknode.pro/70510f2528a6c57bcb376fda51073da9aceb214f/", client.EvmClient)
	tx, err := ethClient.TxByHash(context.Background(), hex.HexstrToHex("0x101239b9ef78582a3d628685f99f596266cff0c25e1c1e868d763428d0c18013"))
	assert.NoError(t, err)
	sea := NewSeaPort()
	eventsHandle := sea.EventLogs()

	for _, rawLog := range tx.Receipt.Logs {
		topic0 := rawLog.GetTopic0()
		handle, ok := eventsHandle[topic0.HexStr()]
		if !ok {
			continue
		}
		assert.True(t, ok)
		eventName, eventArgs := handle(rawLog.GetLogBytes())
		t.Log("event_name", eventName)
		t.Log("event_args", eventArgs)
	}
}

// https://etherscan.io/tx/0x8b75877d48c15aa728089a73100545bb1fd49967144e3cec3e0a25f221beadf8#eventlog
// TODO: 未发生实际交易操作，暂时过滤掉这类
func TestFulfillOrder(t *testing.T) {

	ethClient := client.NewClient("wss://green-small-mountain.discover.quiknode.pro/70510f2528a6c57bcb376fda51073da9aceb214f/", client.EvmClient)
	tx, err := ethClient.TxByHash(context.Background(), hex.HexstrToHex("0x8b75877d48c15aa728089a73100545bb1fd49967144e3cec3e0a25f221beadf8"))
	assert.NoError(t, err)
	sea := NewSeaPort()
	eventsHandle := sea.EventLogs()

	for _, rawLog := range tx.Receipt.Logs {
		topic0 := rawLog.GetTopic0()
		handle, ok := eventsHandle[topic0.HexStr()]
		if !ok {
			continue
		}
		assert.True(t, ok)
		eventName, eventArgs := handle(rawLog.GetLogBytes())
		t.Log("event_name", eventName)
		t.Log("event_args", eventArgs)
	}
}

// https://etherscan.io/tx/0xd94a9d5a0db792d413425ed9140cf7a83cfc04234c028dfbe5836f287b8a6dc8
func TestFulfillAdvancedOrderBid(t *testing.T) {

	ethClient := client.NewClient("wss://green-small-mountain.discover.quiknode.pro/70510f2528a6c57bcb376fda51073da9aceb214f/", client.EvmClient)
	tx, err := ethClient.TxByHash(context.Background(), hex.HexstrToHex("0x835603465C401B1450514572E10AAA394C56CBB9FF18332A73D5F0BDFF1BBBE9"))
	assert.NoError(t, err)
	sea := NewSeaPort()
	eventsHandle := sea.EventLogs()

	for _, rawLog := range tx.Receipt.Logs {

		topic0 := rawLog.GetTopic0()
		handle, ok := eventsHandle[topic0.HexStr()]
		if !ok {
			continue
		}
		assert.True(t, ok)
		eventName, eventArgs := handle(rawLog.GetLogBytes())

		if tx.Receipt.From.EqualTo(hex.HexstrToHex(eventArgs["fulfiller"].(string))) {
			t.Log("tx.Receipt.From", tx.Receipt.From, "fullfiller is", eventArgs["fulfiller"].(string))
		}
		t.Log("event_name", eventName)
		t.Log("event_args", eventArgs)
	}
}
