package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func corsHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		h(w, r)
	}
}

func main() {
	_ = godotenv.Load()

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	db := DBConnection()

	// Finnhub WS
	finnhubWS := connectToFinnhub()
	defer finnhubWS.Close()

	go handleFinnhubMessages(finnhubWS, db)
	go broadcastUpdates()

	http.HandleFunc("/ws", WSHandler)
	http.HandleFunc("/stocks-history", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		StocksHistoryHandler(w, r, db)
	}))
	http.HandleFunc("/stock-candles", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		CandlesHandler(w, r, db)
	}))

	log.Println("🚀 Server running on port", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
