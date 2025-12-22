package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var symbols = []string{"AAPL", "AMZN", "TSLA", "GOOGL"}

var (
	broadcast    = make(chan *BroadcastMessage)
	tempCandles  = make(map[string]*TempCandle)
	mu           sync.Mutex

	clientConns = make(map[*websocket.Conn]string)
	clientsMu   sync.Mutex
)

func connectToFinnhub() *websocket.Conn {
	key := os.Getenv("API_KEY")
	ws, _, err := websocket.DefaultDialer.Dial(
		fmt.Sprintf("wss://ws.finnhub.io?token=%s", key),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range symbols {
		msg, _ := json.Marshal(map[string]string{
			"type":   "subscribe",
			"symbol": s,
		})
		ws.WriteMessage(websocket.TextMessage, msg)
	}

	return ws
}

func handleFinnhubMessages(ws *websocket.Conn, db *gorm.DB) {
	for {
		var msg FinnhubMessage
		if err := ws.ReadJSON(&msg); err != nil {
			continue
		}

		if msg.Type == "trade" {
			for _, trade := range msg.Data {
				processTradeData(&trade, db)
			}
		}
	}
}

func processTradeData(trade *TradeData, db *gorm.DB) {
	mu.Lock()
	defer mu.Unlock()

	symbol := trade.Symbol
	price := trade.Price
	volume := float64(trade.Volume)
	ts := time.UnixMilli(trade.Timestamp)

	candle, exists := tempCandles[symbol]

	if !exists || ts.After(candle.CloseTime) {
		if exists {
			final := candle.toCandle()
			db.Create(final)
			broadcast <- &BroadcastMessage{Closed, final}
		}

		candle = &TempCandle{
			Symbol:     symbol,
			OpenTime:  ts,
			CloseTime: ts.Truncate(time.Minute).Add(time.Minute),
			OpenPrice: price,
			ClosePrice: price,
			HighPrice: price,
			LowPrice:  price,
			Volume:    volume,
		}
		tempCandles[symbol] = candle
		return
	}

	candle.ClosePrice = price
	candle.Volume += volume
	if price > candle.HighPrice {
		candle.HighPrice = price
	}
	if price < candle.LowPrice {
		candle.LowPrice = price
	}

	broadcast <- &BroadcastMessage{Live, candle.toCandle()}
}
