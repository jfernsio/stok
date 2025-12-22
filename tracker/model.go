package main

import "time"

type FinnhubMessage struct {
	Type string      `json:"type"`
	Data []TradeData `json:"data"`
}

type TradeData struct {
	Symbol    string  `json:"s"`
	Price     float64 `json:"p"`
	Volume    int     `json:"v"`
	Timestamp int64   `json:"t"`
}

type Candle struct {
	ID        uint      `gorm:"primaryKey"`
	Symbol    string    `json:"symbol"`
	OpenTime  time.Time `json:"openTime"`
	CloseTime time.Time `json:"closeTime"`
	Open      float64   `json:"open"`
	Close     float64   `json:"close"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Volume    float64   `json:"volume"`
	CreatedAt time.Time
}

type TempCandle struct {
	Symbol     string
	OpenTime   time.Time
	CloseTime  time.Time
	OpenPrice  float64
	ClosePrice float64
	HighPrice  float64
	LowPrice   float64
	Volume     float64
}

func (t *TempCandle) toCandle() *Candle {
	return &Candle{
		Symbol:    t.Symbol,
		OpenTime:  t.OpenTime,
		CloseTime: t.CloseTime,
		Open:      t.OpenPrice,
		Close:     t.ClosePrice,
		High:      t.HighPrice,
		Low:       t.LowPrice,
		Volume:    t.Volume,
	}
}

type UpdateType string

const (
	Live   UpdateType = "LIVE"
	Closed UpdateType = "CLOSED"
)

type BroadcastMessage struct {
	UpdateType UpdateType `json:"type"`
	Candle     *Candle    `json:"candle"`
}
