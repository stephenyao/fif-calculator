package repository

import (
	"fif-clacultor/internal/model"
	"github.com/jmoiron/sqlx"
)

type TradeRepository interface {
	Insert(trade *model.Trade) error
	Update(trade *model.Trade) error
	DeleteByID(id int) error
	GetByID(id int) (*model.Trade, error)
	GetAll() ([]*model.Trade, error)
	GetBySymbol(symbol string) ([]*model.Trade, error)
}

type SQLTradeRepository struct {
	DB *sqlx.DB
}

func NewTradeRepository(db *sqlx.DB) TradeRepository {
	return &SQLTradeRepository{
		DB: db,
	}
}

func (r *SQLTradeRepository) Insert(trade *model.Trade) error {
	_, err := r.DB.NamedExec(`
		INSERT INTO trades (symbol, buy_date, quantity, price, currency, action)
		VALUES (:symbol, :buy_date, :quantity, :price, :currency, :action)
		`, &trade)
	return err
}

func (r *SQLTradeRepository) GetAll() ([]*model.Trade, error) {
	var trades []*model.Trade
	err := r.DB.Select(&trades, `SELECT * FROM trades ORDER BY buy_date DESC`)
	return trades, err
}

func (r *SQLTradeRepository) GetByID(id int) (*model.Trade, error) {
	var trade model.Trade
	err := r.DB.Get(&trade, "SELECT * FROM trades WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &trade, nil
}

func (r *SQLTradeRepository) GetBySymbol(symbol string) ([]*model.Trade, error) {
	var trades []*model.Trade
	err := r.DB.Select(&trades, "SELECT * FROM trades WHERE symbol = ?", symbol)
	if err != nil {
		return nil, err
	}
	return trades, nil
}

func (r *SQLTradeRepository) DeleteByID(id int) error {
	_, err := r.DB.Exec("DELETE FROM trades WHERE id = ?", id)
	return err
}

func (r *SQLTradeRepository) Update(trade *model.Trade) error {
	_, err := r.DB.Exec(`
		UPDATE trades 
		SET symbol = ?, buy_date = ?, quantity = ?, price = ?, currency = ?, action = ?
		WHERE id = ?
	`, trade.Symbol, trade.BuyDate, trade.Quantity, trade.Price, trade.Currency, trade.Action, trade.ID)
	return err
}
