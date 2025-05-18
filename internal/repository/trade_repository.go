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

func (r *TradeRepository) GetAll() ([]*model.Trade, error) {
	var trades []*model.Trade
	err := r.DB.Select(&trades, `SELECT * FROM trades ORDER BY buy_date DESC`)
	return trades, err
}

func (r *TradeRepository) GetByID(id int) (*model.Trade, error) {
	var trade model.Trade
	err := r.DB.Get(&trade, "SELECT * FROM trades WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &trade, nil
}

func (r *TradeRepository) DeleteByID(id int) error {
	_, err := r.DB.Exec("DELETE FROM trades WHERE id = ?", id)
	return err
}
