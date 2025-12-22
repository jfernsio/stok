package main

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

func broadcastUpdates() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	var latest *BroadcastMessage

	for {
		select {
		case msg := <-broadcast:
			if msg.UpdateType == Closed {
				sendToClients(msg)
			} else {
				latest = msg
			}

		case <-ticker.C:
			if latest != nil {
				sendToClients(latest)
				latest = nil
			}
		}
	}
}

func sendToClients(msg *BroadcastMessage) {
	data, _ := json.Marshal(msg)

	clientsMu.Lock()
	defer clientsMu.Unlock()

	for conn, sym := range clientConns {
		if sym == msg.Candle.Symbol {
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				conn.Close()
				delete(clientConns, conn)
			}
		}
	}
}
