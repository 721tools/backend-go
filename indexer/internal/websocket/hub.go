// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/721tools/backend-go/indexer/internal/service"
	"github.com/721tools/backend-go/indexer/pkg/consts"
	"github.com/721tools/backend-go/indexer/pkg/mq"
	"github.com/721tools/backend-go/indexer/pkg/utils/hex"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the mq.
	broadcast chan Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) listenToMQ() {
	mq := mq.GetMQ()
	srv := service.NewNFTItemService()
	var payload Message
	var err error

	go mq.Subscribe(context.TODO(), consts.OPENSEA_ETH_ORDER_LISTING)

	for {
		pop := mq.Pop()

		err = json.Unmarshal([]byte(pop.Payload), &payload)
		if err != nil {
			log.Warn("debug ws", "cant Unmarshal payload", pop.Payload)
			continue
		}

		res := payload.Payload.Payload.(*ItemListedRes)
		NFTId := res.Item.NFTId
		s := strings.Split(NFTId, "/")
		contractAddress := hex.HexstrToHex(s[1])
		tokenId := hex.IntstrToBigInt(s[2]).Bytes()
		if item := srv.GetNFTExtraDataByContractAndTokenId(contractAddress, tokenId); item != nil {
			res.Item.Metadata.FloorPrice = item.Col.FloorPrice
			res.Item.Metadata.Rank = item.Item.TraitsRank
			res.Item.Metadata.Verified = item.Col.Verified
			res.Item.Metadata.TotalSupply = item.Col.TotalSupply
			if res.Item.Metadata.ImageUrl == "" {
				res.Item.Metadata.ImageUrl = item.Item.ImageUrl
			}
		}

		h.broadcast <- payload
	}
}

func (h *Hub) Run() {

	go h.listenToMQ()

	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
