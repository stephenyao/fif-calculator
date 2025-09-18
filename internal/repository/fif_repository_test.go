package repository

import (
	"github.com/jmoiron/sqlx"
	"reflect"
	"testing"
	"time"
)

func TestFIFRepository(t *testing.T) {
	type seedFn func(t *testing.T, db *sqlx.DB)

	cutoff := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

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
					Quantity: 150,
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
					Quantity: 150,
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
				2: {
					Quantity: 0,
					Name:     "Apple",
					Symbol:   "APPL",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			tc.seed(t, db)
			repo := NewFIFSQLRepository(db)
			got := repo.GetHoldingQuantities(tc.holdingIDs, tc.upUntil)

			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
