package repository

import (
	"fif-calculator/internal/model"
	"github.com/jmoiron/sqlx"
)

type TradeRepository interface {
	Insert(userID string, trade *model.Trade) error
	Update(userID string, trade *model.Trade) error
	DeleteByID(id int, userID string) error
	GetByID(id int, userID string) (*model.Trade, error)
	GetAllByAscendingDate(userID string) ([]model.Trade, error)
	GetByHoldingID(id int, userID string) ([]model.Trade, error)
}

type SQLTradeRepository struct {
	DB *sqlx.DB
}

func NewTradeRepository(db *sqlx.DB) TradeRepository {
	return &SQLTradeRepository{
		DB: db,
	}
}

func (r *SQLTradeRepository) Insert(userID string, trade *model.Trade) error {
	// Use an INSERT ... SELECT pattern to ensure the holding belongs to the user
	_, err := r.DB.Exec(`
		INSERT INTO trades (buy_date, quantity, price, exchange_rate, action, holding_id)
		SELECT ?, ?, ?, ?, ?, h.id
		FROM holdings h
		WHERE h.id = ? AND h.user_id = ?
	`,
		trade.BuyDate,
		trade.Quantity,
		trade.Price,
		trade.ExchangeRate,
		trade.Action,
		trade.HoldingID,
		userID,
	)

	return err
}

func (r *SQLTradeRepository) GetAllByAscendingDate(userID string) ([]model.Trade, error) {
	var trades []model.Trade

	err := r.DB.Select(&trades, `
		SELECT 
			trades.id,
			trades.buy_date,
			trades.quantity,
			trades.price,
			trades.exchange_rate,
			trades.action,			
			trades.holding_id,
			holdings.currency AS currency,
			holdings.name AS holding_name,
			holdings.symbol AS symbol
		FROM trades
		JOIN holdings ON trades.holding_id = holdings.id
		WHERE holdings.user_id = ?
		ORDER BY trades.buy_date ASC
	`, userID)

	return trades, err
}

func (r *SQLTradeRepository) GetByID(id int, userID string) (*model.Trade, error) {
	var trade model.Trade
	err := r.DB.Get(&trade, `
		SELECT 
			trades.id,
			trades.buy_date,
			trades.quantity,
			trades.price,
			trades.exchange_rate,
			trades.action,			
			trades.holding_id,
			holdings.currency AS currency,
			holdings.name AS holding_name,
			holdings.symbol AS symbol
		FROM trades
		JOIN holdings ON trades.holding_id = holdings.id
		WHERE trades.id = ? AND holdings.user_id = ?
	`, id, userID)

	if err != nil {
		return nil, err
	}
	return &trade, nil
}

func (r *SQLTradeRepository) DeleteByID(tradeID int, userID string) error {
	_, err := r.DB.Exec(`
		DELETE FROM trades
		WHERE id = ?
		AND holding_id IN (
			SELECT id FROM holdings WHERE user_id = ?
		)
	`, tradeID, userID)
	return err
}

func (r *SQLTradeRepository) Update(userID string, trade *model.Trade) error {
	_, err := r.DB.Exec(`
		UPDATE trades 
		SET buy_date = ?, quantity = ?, price = ?,  exchange_rate = ?, action = ?
		WHERE id = ?
		  AND holding_id IN (
			  SELECT id FROM holdings WHERE user_id = ?
		  )
	`,
		trade.BuyDate,
		trade.Quantity,
		trade.Price,
		trade.ExchangeRate,
		trade.Action,
		trade.ID,
		userID,
	)
	return err
}

func (r *SQLTradeRepository) GetByHoldingID(holdingID int, userID string) ([]model.Trade, error) {
	var trades []model.Trade
	err := r.DB.Select(&trades, `
		SELECT 
			trades.id,
			trades.buy_date,
			trades.quantity,
			trades.price,
			trades.action,
			trades.holding_id,
			holdings.currency AS currency,
			holdings.name AS holding_name,
			holdings.symbol AS symbol
		FROM trades
		JOIN holdings ON trades.holding_id = holdings.id
		WHERE trades.holding_id = ?
		  AND holdings.user_id = ?
		ORDER BY trades.buy_date ASC
	`, holdingID, userID)

	if err != nil {
		return nil, err
	}
	return trades, nil
}
