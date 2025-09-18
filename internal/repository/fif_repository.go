package repository

import (
	"fif-calculator/internal/constants"
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

type fifHoldingActivityRecord struct {
	Date          time.Time
	Action        string  `db:"action"`
	Quantity      float64 `db:"quantity"`
	Price         float64 `db:"price"`
	ExchangeRate  float64 `db:"exchange_rate"`
	HoldingID     int     `db:"holding_id"`
	HoldingName   string  `db:"name"`
	HoldingSymbol string  `db:"symbol"`
}

type FIFSQLRepository struct {
	db *sqlx.DB
}

func NewFIFSQLRepository(db *sqlx.DB) FIFRepository {
	return &FIFSQLRepository{db: db}
}

func (r FIFSQLRepository) GetHoldingQuantities(holdingsIDs []HoldingID, upUntil time.Time) map[HoldingID]FIFHoldingQuantity {
	results := make(map[HoldingID]FIFHoldingQuantity)
	holdingActivities := make(map[HoldingID]bool)
	activities, err := r.fetchActivities(holdingsIDs, upUntil)

	if err != nil {
		panic(err)
	}

	holdingQuantitiesMap := make(map[HoldingID]float64)
	holdingInfo := make(map[HoldingID]struct {
		name   string
		symbol string
	})

	for _, activity := range activities {
		holdingID := HoldingID(activity.HoldingID)
		holdingActivities[holdingID] = true
		switch activity.Action {
		case constants.Buy:
			holdingQuantitiesMap[holdingID] += activity.Quantity
			holdingInfo[holdingID] = struct {
				name   string
				symbol string
			}{name: activity.HoldingName, symbol: activity.HoldingSymbol}
		case constants.Sell:
			newQty := holdingQuantitiesMap[holdingID] - activity.Quantity
			holdingQuantitiesMap[holdingID] -= newQty
			holdingQuantitiesMap[holdingID] = max(0, newQty)
		}
	}

	for _, id := range holdingsIDs {
		holdingID := HoldingID(id)

		// If the holding had no activities then do not return a quantity
		if _, ok := holdingActivities[holdingID]; !ok {
			continue
		}

		results[id] = FIFHoldingQuantity{
			Quantity: holdingQuantitiesMap[holdingID],
			Name:     holdingInfo[holdingID].name,
			Symbol:   holdingInfo[holdingID].symbol,
		}
	}

	return results
}

func (r FIFSQLRepository) GetTrades(holdingsIDs []HoldingID, startDate time.Time, endDate time.Time) map[HoldingID][]FIFTradeActivity {
	return make(map[HoldingID][]FIFTradeActivity)
}

func (r FIFSQLRepository) fetchActivities(holdingsIDs []HoldingID, upUntil time.Time) ([]fifHoldingActivityRecord, error) {
	timeStr := upUntil.Format(time.DateOnly)

	base := `
		SELECT 
		    trades.action, trades.quantity, trades.price, trades.exchange_rate, trades.holding_id, holdings.name, holdings.symbol		    
		FROM
		    trades
		JOIN holdings ON trades.holding_id = holdings.id
		WHERE buy_date <= ? AND holding_id IN (?)
		ORDER BY buy_date ASC;
	`

	query, args, err := sqlx.In(base, timeStr, holdingsIDs)
	if err != nil {
		return nil, err
	}

	// Rebind for the current driver (SQLite, MySQL, Postgres use different placeholders)
	query = r.db.Rebind(query)

	var activities []fifHoldingActivityRecord
	// Execute
	if err := r.db.Select(&activities, query, args...); err != nil {
		return nil, err
	}

	return activities, nil
}
