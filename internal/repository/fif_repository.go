package repository

import (
	"fif-calculator/internal/constants"
	"time"

	"github.com/jmoiron/sqlx"
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
	Date         time.Time
	Action       string  `db:"action"`
	Quantity     float64 `db:"quantity"`
	Price        float64 `db:"price"`
	ExchangeRate float64 `db:"exchange_rate"`
	HoldingID    int     `db:"holding_id"`
}

type holdingRecord struct {
	Id     int    `db:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
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
	holdings, err := r.fetchHoldings(holdingsIDs)
	activities, err := r.fetchActivities(holdingsIDs, upUntil)

	if err != nil {
		panic(err)
	}

	holdingQuantitiesMap := make(map[HoldingID]float64)

	for _, activity := range activities {
		holdingID := HoldingID(activity.HoldingID)
		holdingActivities[holdingID] = true
		switch activity.Action {
		case constants.Buy:
			holdingQuantitiesMap[holdingID] += activity.Quantity
		case constants.Sell:
			newQty := max(0, holdingQuantitiesMap[holdingID]-activity.Quantity)
			holdingQuantitiesMap[holdingID] = newQty
		}
	}

	for _, id := range holdingsIDs {
		holdingID := HoldingID(id)
		results[id] = FIFHoldingQuantity{
			Quantity: holdingQuantitiesMap[holdingID],
			Name:     holdings[holdingID].Name,
			Symbol:   holdings[holdingID].Symbol,
		}
	}

	return results
}

func (r FIFSQLRepository) GetTrades(holdingsIDs []HoldingID, startDate time.Time, endDate time.Time) map[HoldingID][]FIFTradeActivity {
	return make(map[HoldingID][]FIFTradeActivity)
}

func (r FIFSQLRepository) fetchHoldings(holdingIDs []HoldingID) (map[HoldingID]holdingRecord, error) {
	results := make(map[HoldingID]holdingRecord)

	query := `
		SELECT holdings.id, holdings.name, holdings.symbol
		FROM holdings
		WHERE holdings.id IN (?)
	`

	query, args, err := sqlx.In(query, holdingIDs)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)

	var holdings []holdingRecord

	if err := r.db.Select(&holdings, query, args...); err != nil {
		return nil, err
	}

	for _, holding := range holdings {
		results[HoldingID(holding.Id)] = holding
	}

	return results, nil
}

func (r FIFSQLRepository) fetchActivities(holdingsIDs []HoldingID, upUntil time.Time) ([]fifHoldingActivityRecord, error) {
	timeStr := upUntil.Format(time.DateOnly)

	base := `
		SELECT 
		    trades.action, trades.quantity, trades.price, trades.exchange_rate, trades.holding_id		    
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
