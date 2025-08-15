package viewmodel

import (
	"fif-calculator/internal/model"
	"fmt"
)

type TradeDetailsViewModel struct {
	HoldingName  string
	Symbol       string
	Date         string
	Quantity     string
	Price        string
	Currency     string
	ExchangeRate string
	Action       string
}

func NewTradeDetailsViewModel(trade model.Trade) TradeDetailsViewModel {
	return TradeDetailsViewModel{
		HoldingName:  trade.HoldingName,
		Symbol:       trade.Symbol,
		Date:         trade.BuyDate,
		Quantity:     fmt.Sprintf("%2.f", trade.Quantity),
		Price:        fmt.Sprintf("%.2f", trade.Price),
		Currency:     trade.Currency,
		ExchangeRate: fmt.Sprintf("%.2f", trade.ExchangeRate),
		Action:       trade.Action,
	}
}
