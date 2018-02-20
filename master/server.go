// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package master

import (
	"log"
	"net/http"
	"sort"

	"github.com/satori/go.uuid"

	"github.com/DankBotList/Courier/messaging"
)

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan *messaging.Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// authKey to authenticate against.
	authKey string

	// Collected message UUID's so as to not repeat messages.
	collectedMessages []uuid.UUID
}

func NewHub(authKey string) *Hub {
	return &Hub{
		broadcast:         make(chan *messaging.Message),
		register:          make(chan *Client),
		unregister:        make(chan *Client),
		clients:           make(map[*Client]bool),
		authKey:           authKey,
		collectedMessages: make([]uuid.UUID, 100),
	}
}

// addMessageID add message id's keeping the length of the slice at 100, SENSITIVE!!!!!
func (h *Hub) addMessageID(id uuid.UUID) {
	if len(h.collectedMessages) > 100 {
		h.collectedMessages = h.collectedMessages[len(h.collectedMessages)-97 : len(h.collectedMessages)]
	}
	h.collectedMessages = append(h.collectedMessages, id)
}

func (h *Hub) run() {
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
			i := sort.Search(len(h.collectedMessages), func(i int) bool { return h.collectedMessages[i] == message.ID })
			if i < len(h.collectedMessages) && h.collectedMessages[i] == message.ID {
				continue
			} else {
				h.addMessageID(message.ID)
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
}

// serveWs handles websocket requests from the peer.
func (h *Hub) ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: h, conn: conn, send: make(chan interface{}, 256)}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump(h.authKey)
}
