package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

func WSHandler(w http.ResponseWriter, r *http.Request) {
	up := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := up.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	_, msg, err := conn.ReadMessage()
	if err != nil {
		return
	}

	clientsMu.Lock()
	clientConns[conn] = string(msg)
	clientsMu.Unlock()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			clientsMu.Lock()
			delete(clientConns, conn)
			clientsMu.Unlock()
			conn.Close()
			break
		}
	}
}

func StocksHistoryHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var candles []Candle
	db.Order("open_time desc").Find(&candles)
	json.NewEncoder(w).Encode(candles)
}

func CandlesHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	symbol := r.URL.Query().Get("symbol")
	var candles []Candle
	db.Where("symbol = ?", symbol).Order("open_time desc").Find(&candles)
	json.NewEncoder(w).Encode(candles)
}
