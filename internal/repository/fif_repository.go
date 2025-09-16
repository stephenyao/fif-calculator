package repository

import "time"

type HoldingID int

type FIFRepository interface {
	GetHoldingQuantities(holdingsIDs []HoldingID, upUntil time.Time) map[HoldingID]FIFHoldingQuantity
	GetTrades(holdingsIDs []HoldingID, startDate time.Time, endDate time.Time) map[HoldingID][]FIFTradeActivity
}

type FIFHoldingQuantity struct {
	Quantity float64
	Name     string
	Symbol   string
}

type FIFTradeActivity struct {
	Date         time.Time
	Action       string
	Quantity     float64
	Price        float64
	ExchangeRate float64
	AmountInNZD  float64
}
