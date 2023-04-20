// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"net/http"
	"time"

	"github.com/721tools/stream-api-go/sdk"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan Message

	// Joined Topics
	topics map[string]bool

	// The latest response timestamp
	livingUntil time.Time
}

// isLiving returns true if client's the latest response time is nearly
func (c *Client) isLiving() bool {
	if time.Since(c.livingUntil) > time.Second*60 {
		return false
	}
	return true
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		msg := &Message{}
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Warn("debug IsUnexpectedCloseError", "error", err)
				break
			}
			log.Info("debug readPump", "msg decode error", err, "msg is", msg)
			continue
		}
		switch msg.Event {
		case "phx_join":
			slug := msg.getSlugNameFromTopic()
			c.topics[slug] = true
		case "phx_leave":
			slug := msg.getSlugNameFromTopic()
			if _, ok := c.topics[slug]; ok {
				delete(c.topics, slug)
			}
		case "heartbeat":
			c.livingUntil = time.Now()
		default:
			log.Info("debug ws", "unknown evnet type, msg is", msg)
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Info("debug writePump", "get out of writePump", "close connection")
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			slug := message.getSlugNameFromTopic()
			// c.topic [* or slug]
			_, ok1 := c.topics[sdk.ITEM_FEACH_ALL]
			if _, ok := c.topics[slug]; !ok1 && !ok {
				continue
			}
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteJSON(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				c.conn.WriteJSON(<-c.send)
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Info("debug servews", "req", r)

	for k, v := range r.Header {
		log.Info("debug servews header", "k", k, "v", v)
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Info("debug ServeWs", "err", err)
		return
	}
	client := &Client{
		hub:         hub,
		conn:        conn,
		topics:      make(map[string]bool),
		send:        make(chan Message, 10),
		livingUntil: time.Now(),
	}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
