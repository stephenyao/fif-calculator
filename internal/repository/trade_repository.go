package repository

import (
	"fif-calculator/internal/model"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type TradeRepository interface {
	Insert(trade *model.Trade) error
	Update(trade *model.Trade) error
	DeleteByID(id int) error
	GetByID(id int) (*model.Trade, error)
	GetAll() ([]model.Trade, error)
	GetAllByAscendingDate() ([]model.Trade, error)
	GetByHoldingID(id int) ([]model.Trade, error)
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
		INSERT INTO trades (symbol, buy_date, quantity, price, currency, action, holding_id)
		VALUES (:symbol, :buy_date, :quantity, :price, :currency, :action, :holding_id)
	`, trade)
	return err
}

func (r *SQLTradeRepository) GetAll() ([]model.Trade, error) {
	var trades []model.Trade

	query := fmt.Sprintf("SELECT * FROM trades ORDER BY buy_date DESC")
	err := r.DB.Select(&trades, query)
	return trades, err
}

func (r *SQLTradeRepository) GetAllByAscendingDate() ([]model.Trade, error) {
	var trades []model.Trade

	query := fmt.Sprintf("SELECT * FROM trades ORDER BY buy_date ASC")
	err := r.DB.Select(&trades, query)
	return trades, err
}

func (r *SQLTradeRepository) GetByID(id int) (*model.Trade, error) {
	var trade model.Trade
	err := r.DB.Get(&trade, `
		SELECT 
			trades.id,
			trades.symbol,
			trades.buy_date,
			trades.quantity,
			trades.price,
			trades.currency,
			trades.action,
			trades.holding_id,
			holdings.name AS holding_name
		FROM trades
		JOIN holdings ON trades.holding_id = holdings.id
		WHERE trades.id = ?
	`, id)

	if err != nil {
		return nil, err
	}
	return &trade, nil
}

func (r *SQLTradeRepository) DeleteByID(id int) error {
	_, err := r.DB.Exec("DELETE FROM trades WHERE id = ?", id)
	return err
}

func (r *SQLTradeRepository) Update(trade *model.Trade) error {
	_, err := r.DB.Exec(`
		UPDATE trades 
		SET symbol = ?, buy_date = ?, quantity = ?, price = ?, currency = ?, action = ?, holding_id = ?
		WHERE id = ?
	`, trade.Symbol, trade.BuyDate, trade.Quantity, trade.Price, trade.Currency, trade.Action, trade.HoldingID, trade.ID)
	return err
}

func (r *SQLTradeRepository) GetByHoldingID(id int) ([]model.Trade, error) {
	var trades []model.Trade
	err := r.DB.Select(&trades, `
	SELECT 
		trades.id,
		trades.symbol,
		trades.buy_date,
		trades.quantity,
		trades.price,
		trades.currency,
		trades.action,
		trades.holding_id,
		holdings.name AS holding_name
	FROM trades
	JOIN holdings ON trades.holding_id = holdings.id
	WHERE trades.holding_id = ?
	ORDER BY trades.buy_date ASC
`, id)

	if err != nil {
		return nil, err
	}
	return trades, nil
}
