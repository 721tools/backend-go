package mq

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMq(t *testing.T) {
	rdb := NewMQ(context.TODO(), "redis://:@cachecluster.rfkqxz.ng.0001.ape1.cache.amazonaws.com:6379/0")
	assert.NotNil(t, rdb)
}

func TestSubscribe(t *testing.T) {
	rdb := NewMQ(context.TODO(), "redis://:@localhost:6379/0")
	assert.NotNil(t, rdb)

	go rdb.Subscribe(context.TODO(), "OPENSEA-ETH-ORDER-LISTING")

	for {
		msg := rdb.Pop()
		fmt.Println(msg.Channel, msg.Payload)
	}
}
