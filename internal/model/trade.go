package model

import "time"

type Trade struct {
	ID       int       `db:"id"`
	Symbol   string    `db:"symbol"`
	BuyDate  time.Time `db:"buy_date"`
	Quantity float64   `db:"quantity"`
	Price    float64   `db:"price"`
	Currency string    `db:"currency"`
	Action   string    `db:"action"`
}
