package model

import "time"

type HoldingRecord struct {
	ID        int       `db:"id"`
	UserID    string    `db:"user_id"`
	Name      string    `db:"name"`
	Symbol    string    `db:"symbol"`
	Currency  string    `db:"currency"`
	CreatedAt time.Time `db:"created_at"`
}
