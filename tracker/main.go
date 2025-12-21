package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var (
	symbols = []string{"AAPL", "AMZN", "TSLA", "GOOGL", "NFLX", "PYPL"}

	temCandles = make(map[string]*TempCandle)
	mu sync.Mutex

)


func main() {
	env := os.Getenv()
	db := DBConnection(env)

	//connect to finhub websckt
	finhubWSCon := connectToFinhub(env)
	defer finhubWSCon.Close(env)
	//handle finhub ws incoming msgs
	go handleFinhubMesgs(finhubWSCon,db)
	//broadcast candle updates to all clients


	//endpoints
	//conn to ws
	//fetch all past candels for all symbols
	//fetch past candels for single symbol
	 
}

func connectToFinhub(env *Env) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("wss://ws.finnhub.io?token=%s", env.API_KEY), nil)
	if err != nil {
		panic(err)
	}

	for _, s := range symbols {
		msg, _ := json.Marshal(map[string]interface{}{"type": "subscribe", "symbol": s})
		ws.WriteMessage(websocket.TextMessage, msg)
	}

	return ws
}

func handleFinhubMesgs(ws *websocket.Conn, db *gorm.DB) {
	finnhubMessage := &FinnhubMessage{}

	for {
		if err := ws.ReadJSON(finnhubMessage); err != nil {
			fmt.Println("Error reading message: ",err)
			continue
		}

		//only try to process the message data if type trade
		if (finnhubMessage.Type == "trade" ) {
			for _,trade := range finnhubMessage.Data {
				//proricess trade data
				processTradeData(trade,db)
			}
		}
	}
}

func processTradeData( trade *TradeData,db *gorm.DB) {
	//prorext go rountines
	mu.Lock()
	defer mu.Unlock()
}