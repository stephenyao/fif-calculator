package repository

import (
	"reflect"
	"testing"
	"time"
)

func TestFIFRepository(t *testing.T) {
	db := setupTestDB(t)
	seedDB(t, db)

	t.Run("Get HoldingQuantities", func(t *testing.T) {
		repository := NewFIFSQLRepository(db)

		got := repository.GetHoldingQuantities(
			[]HoldingID{
				1, 2,
			},
			time.Now(),
		)

		var want map[HoldingID]FIFHoldingQuantity = map[HoldingID]FIFHoldingQuantity{
			1: {
				Quantity: 50,
				Name:     "Google",
				Symbol:   "GOOG",
			},
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
