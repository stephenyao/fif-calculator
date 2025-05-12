package repository

import (
	"fif-clacultor/internal/model"
	"github.com/jmoiron/sqlx"
)

type TradeRepository struct {
	DB *sqlx.DB
}

func NewTradeRepository(db *sqlx.DB) *TradeRepository {
	return &TradeRepository{
		DB: db,
	}
}

func (r *TradeRepository) Insert(trade *model.Trade) error {
	_, err := r.DB.NamedExec(`
		INSERT INTO trades (symbol, buy_date, quantity, price, currency, action)
		VALUES (:symbol, :buy_date, :quantity, :price, :currency, :action)
		`, &trade)
	return err
}
