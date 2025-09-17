package repository

import (
	"fif-calculator/internal/constants"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type HoldingID int

type FIFRepository interface {
	GetHoldingQuantities(holdingsIDs []HoldingID, upUntil time.Time) map[HoldingID]FIFHoldingQuantity
	GetTrades(holdingsIDs []HoldingID, startDate time.Time, endDate time.Time) map[HoldingID][]FIFTradeActivity
}

type FIFHoldingQuantity struct {
	Quantity float64
	Name     string
	Symbol   string
}

type FIFTradeActivity struct {
	Date         time.Time
	Action       string  `db:"action"`
	Quantity     float64 `db:"quantity"`
	Price        float64 `db:"price"`
	ExchangeRate float64 `db:"exchange_rate"`
	HoldingID    int     `db:"holding_id"`
	AmountInNZD  float64
}

type FIFSQLRepository struct {
	db *sqlx.DB
}

func NewFIFSQLRepository(db *sqlx.DB) FIFRepository {
	return &FIFSQLRepository{db: db}
}

func (r FIFSQLRepository) GetHoldingQuantities(holdingsIDs []HoldingID, upUntil time.Time) map[HoldingID]FIFHoldingQuantity {
	var activities []FIFTradeActivity
	var results = make(map[HoldingID]FIFHoldingQuantity)

	timeStr := upUntil.Format(time.DateOnly)

	base := `
		SELECT 
		    action, quantity, price, exchange_rate, holding_id
		FROM
		    trades
		WHERE buy_date <= ?
			AND holding_id IN (?)
		ORDER BY buy_date ASC;
	`

	query, args, err := sqlx.In(base, timeStr, holdingsIDs)
	if err != nil {
		panic(err)
	}

	// Rebind for the current driver (SQLite, MySQL, Postgres use different placeholders)
	query = r.db.Rebind(query)

	// Execute
	if err := r.db.Select(&activities, query, args...); err != nil {
		panic(err)
	}

	var holdingQuantitiesMap = make(map[HoldingID]float64)

	for _, activity := range activities {
		switch activity.Action {
		case constants.Buy:
			holdingQuantitiesMap[HoldingID(activity.HoldingID)] += activity.Quantity
		case constants.Sell:
			newQty := holdingQuantitiesMap[HoldingID(activity.HoldingID)] - activity.Quantity
			holdingQuantitiesMap[HoldingID(activity.HoldingID)] -= max(0, newQty)
		}
	}

	for _, id := range holdingsIDs {
		results[id] = FIFHoldingQuantity{
			Quantity: holdingQuantitiesMap[HoldingID(id)],
			Name:     "test",
			Symbol:   "test",
		}
	}

	// Debug print
	fmt.Printf("%+v\n", activities)

	return results
}

func (r FIFSQLRepository) GetTrades(holdingsIDs []HoldingID, startDate time.Time, endDate time.Time) map[HoldingID][]FIFTradeActivity {
	return make(map[HoldingID][]FIFTradeActivity)
}
