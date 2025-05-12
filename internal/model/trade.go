package model

type Trade struct {
	ID       int     `db:"id"`
	Symbol   string  `db:"symbol"`
	BuyDate  string  `db:"buy_date"`
	Quantity float64 `db:"quantity"`
	Price    float64 `db:"price"`
	Currency string  `db:"currency"`
	Action   string  `db:"action"`
}
