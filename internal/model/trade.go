package model

type Trade struct {
	ID           int     `db:"id"`
	BuyDate      string  `db:"buy_date"`
	Quantity     float64 `db:"quantity"`
	Price        float64 `db:"price"`
	ExchangeRate float64 `db:"exchange_rate"`
	Action       string  `db:"action"`
	HoldingID    int     `db:"holding_id"`
	Currency     string  `db:"currency"`
	HoldingName  string  `db:"holding_name"`
	Symbol       string  `db:"symbol"`
}

func (t Trade) PriceInNZD() float64 {
	return t.Price * t.ExchangeRate
}
