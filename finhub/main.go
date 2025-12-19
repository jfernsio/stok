package main

import (
    "context"
    "encoding/json"
    "fmt"

    finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
)

func main() {
    cfg := finnhub.NewConfiguration()
    cfg.AddDefaultHeader("X-Finnhub-Token", "d52p241r01qggm5t4hn0d52p241r01qggm5t4hng")

    client := finnhub.NewAPIClient(cfg).DefaultApi

    candles, _, err := client.Quote(context.Background()).
        Symbol("AAPL").
      
     
        Execute()

    if err != nil {
        panic(err)
    }

    b, err := json.MarshalIndent(candles, "", "  ")
    if err != nil {
        panic(err)
    }

	insider,_,err := client.InsiderTransactions(context.Background()). Symbol("AAPL").	Execute()
	 if err != nil {
        panic(err)
    }

    c, err := json.MarshalIndent(insider, "", "  ")
    if err != nil {
        panic(err)
    }
    fmt.Println(string(b))
	fmt.Println("insider",string(c))

}
