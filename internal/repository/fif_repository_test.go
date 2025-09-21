package repository

import (
	"reflect"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
)

func TestFIFRepository(t *testing.T) {
	type seedFn func(t *testing.T, db *sqlx.DB)

	cutoff := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	t.Run("GetHoldingQuantities", func(t *testing.T) {
		testCases := []struct {
			name       string
			seed       seedFn
			holdingIDs []HoldingID
			upUntil    time.Time
			want       map[HoldingID]FIFHoldingQuantity
		}{
			{
				name: "single holding",
				seed: func(t *testing.T, db *sqlx.DB) {
					insertHoldings(t, db)
					insertTradesSingleHolding(t, db)
				},
				holdingIDs: []HoldingID{1, 2},
				upUntil:    cutoff,
				want: map[HoldingID]FIFHoldingQuantity{
					1: {
						Quantity: 150,
						Name:     "Google",
						Symbol:   "GOOG",
					},
					2: {
						Quantity: 0,
						Name:     "Apple",
						Symbol:   "APPL",
					},
				},
			},
			{
				name: "only sell trades",
				seed: func(t *testing.T, db *sqlx.DB) {
					insertHoldings(t, db)
					insertOnlySellTrades(t, db)
				},
				holdingIDs: []HoldingID{1, 2},
				want: map[HoldingID]FIFHoldingQuantity{
					1: {
						Quantity: 0,
						Name:     "Google",
						Symbol:   "GOOG",
					},
					2: {
						Quantity: 0,
						Name:     "Apple",
						Symbol:   "APPL",
					},
				},
			},
			{
				name: "multiple holdings",
				seed: func(t *testing.T, db *sqlx.DB) {
					insertHoldings(t, db)
					insertTradesMultipleHoldings(t, db)
				},
				holdingIDs: []HoldingID{1, 2},
				upUntil:    cutoff,
				want: map[HoldingID]FIFHoldingQuantity{
					1: {
						Quantity: 200,
						Name:     "Google",
						Symbol:   "GOOG",
					},
					2: {
						Quantity: 50,
						Name:     "Apple",
						Symbol:   "APPL",
					},
				},
			},
			{
				name: "multiple holdings querying 1 holding",
				seed: func(t *testing.T, db *sqlx.DB) {
					insertHoldings(t, db)
					insertTradesMultipleHoldings(t, db)
				},
				holdingIDs: []HoldingID{1},
				upUntil:    cutoff,
				want: map[HoldingID]FIFHoldingQuantity{
					1: {
						Quantity: 200,
						Name:     "Google",
						Symbol:   "GOOG",
					},
				},
			},
			{
				name: "sell quantity greater than buy quantity",
				seed: func(t *testing.T, db *sqlx.DB) {
					insertHoldings(t, db)
					insertTradesNegativeQuantity(t, db)
				},
				holdingIDs: []HoldingID{1, 2},
				upUntil:    cutoff,
				want: map[HoldingID]FIFHoldingQuantity{
					1: {
						Quantity: 0,
						Name:     "Google",
						Symbol:   "GOOG",
					},
					2: {
						Quantity: 0,
						Name:     "Apple",
						Symbol:   "APPL",
					},
				},
			},
			{
				name: "no holding ids",
				seed: func(t *testing.T, db *sqlx.DB) {
					insertHoldings(t, db)
					insertTradesNegativeQuantity(t, db)
				},
				holdingIDs: []HoldingID{},
				upUntil:    cutoff,
				want:       map[HoldingID]FIFHoldingQuantity{},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				db := setupTestDB(t)
				tc.seed(t, db)
				repo := NewFIFSQLRepository(db)
				got, err := repo.GetHoldingQuantities(tc.holdingIDs, tc.upUntil)

				if err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(tc.want, got) {
					t.Errorf("got %v, want %v", got, tc.want)
				}
			})
		}
	})

	t.Run("GetTrades", func(t *testing.T) {
		testCases := []struct {
			name       string
			seed       seedFn
			holdingIDs []HoldingID
			startDate  time.Time
			endDate    time.Time
			want       map[HoldingID][]FIFTradeActivity
		}{
			{
				name: "fetches trades within dates",
				seed: func(t *testing.T, db *sqlx.DB) {
					insertHoldings(t, db)
					insertSingleTrade(t, db)
				},
				holdingIDs: []HoldingID{1, 2},
				startDate:  time.Unix(0, 0),
				endDate:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				want: map[HoldingID][]FIFTradeActivity{
					1: {
						{
							Date:         time.Date(2024, 8, 8, 0, 0, 0, 0, time.UTC),
							Action:       "sell",
							Quantity:     100,
							Price:        100,
							ExchangeRate: 1.6,
							HoldingID:    1,
						},
					},
				},
			}, {
				name: "does not fetch trades outside of dates",
				seed: func(t *testing.T, db *sqlx.DB) {
					insertHoldings(t, db)
					insertTradesOutsideDate(t, db)
				},
				holdingIDs: []HoldingID{1, 2},
				startDate:  time.Date(2024, 7, 2, 0, 0, 0, 0, time.UTC),
				endDate:    time.Date(2024, 8, 31, 0, 0, 0, 0, time.UTC),
				want: map[HoldingID][]FIFTradeActivity{
					1: {
						{
							Date:         time.Date(2024, 8, 8, 0, 0, 0, 0, time.UTC),
							Action:       "sell",
							Quantity:     100,
							Price:        100,
							ExchangeRate: 1.6,
							HoldingID:    1,
						},
					},
				},
			}, {
				name: "fetches trades inclusive of dates",
				seed: func(t *testing.T, db *sqlx.DB) {
					insertHoldings(t, db)
					insertTradesOutsideDate(t, db)
				},
				holdingIDs: []HoldingID{1, 2},
				startDate:  time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
				endDate:    time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC),
				want: map[HoldingID][]FIFTradeActivity{
					1: {
						{
							Date:         time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
							Action:       "sell",
							Quantity:     100,
							Price:        100,
							ExchangeRate: 1.6,
							HoldingID:    1,
						},
						{
							Date:         time.Date(2024, 8, 8, 0, 0, 0, 0, time.UTC),
							Action:       "sell",
							Quantity:     100,
							Price:        100,
							ExchangeRate: 1.6,
							HoldingID:    1,
						},
						{
							Date:         time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC),
							Action:       "sell",
							Quantity:     100,
							Price:        100,
							ExchangeRate: 1.6,
							HoldingID:    1,
						},
					},
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				db := setupTestDB(t)
				tc.seed(t, db)
				repo := NewFIFSQLRepository(db)
				got, err := repo.GetTrades(tc.holdingIDs, tc.startDate, tc.endDate)
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(tc.want, got) {
					t.Errorf("got %v, want %v", got, tc.want)
				}
			})
		}
	})
}
