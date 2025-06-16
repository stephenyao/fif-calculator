package costbasisservice

import (
	"fif-calculator/internal/constants"
	"fif-calculator/internal/model"
	"testing"
	"time"
)

var untilDate = time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)

func TestCostBasisBySymbol(t *testing.T) {
	t.Run("buy and sell orders of exact same quantity", func(t *testing.T) {
		var trades []model.Trade = []model.Trade{
			{0, "XYZ", "2024-01-01", 15, 100, "USD", constants.Buy},
			{0, "XYZ", "2024-02-01", 10, 100, "USD", constants.Buy},
			{0, "XYZ", "2024-03-01", 30, 300, "USD", constants.Sell},
		}

		costBasisBySymbol := CostBasisBySymbol(trades, untilDate)
		got, ok := costBasisBySymbol["XYZ"]
		costBasis := got.CostBasis

		if !ok {
			t.Errorf("not ok")
		}

		if costBasis != 0 {
			t.Errorf("got %f, want 0", costBasis)
		}
	})

	t.Run("buy orders is greater than sell order quantity", func(t *testing.T) {
		var trades []model.Trade = []model.Trade{
			{0, "XYZ", "2024-01-01", 10, 100, "USD", constants.Buy},
			{0, "XYZ", "2024-02-01", 20, 100, "USD", constants.Buy},
			{0, "XYZ", "2024-03-01", 10, 300, "USD", constants.Sell},
		}

		costBasisBySymbol := CostBasisBySymbol(trades, untilDate)
		got, ok := costBasisBySymbol["XYZ"]
		var expected float64 = 2000
		costBasis := got.CostBasis

		if !ok {
			t.Errorf("not ok")
		}

		if costBasis != expected {
			t.Errorf("got %f, want %f", costBasis, expected)
		}
	})

	t.Run("buy orders is greater than sell order with overflow to next buy order", func(t *testing.T) {
		var trades []model.Trade = []model.Trade{
			{0, "XYZ", "2024-01-01", 15, 110, "USD", constants.Buy},
			{0, "XYZ", "2024-02-01", 10, 110, "USD", constants.Buy},
			{0, "XYZ", "2024-03-01", 5, 300, "USD", constants.Sell},
		}

		costBasisBySymbol := CostBasisBySymbol(trades, untilDate)
		got, ok := costBasisBySymbol["XYZ"]
		var expected float64 = 2200
		costBasis := got.CostBasis

		if !ok {
			t.Errorf("not ok")
		}

		if costBasis != expected {
			t.Errorf("got %f, want %f", costBasis, expected)
		}
	})

	t.Run("buy orders is greater than sell order with overflow to second buy order", func(t *testing.T) {
		var trades []model.Trade = []model.Trade{
			{0, "XYZ", "2024-01-01", 10, 100, "USD", constants.Buy},
			{0, "XYZ", "2024-02-01", 10, 110, "USD", constants.Buy},
			{0, "XYZ", "2024-03-01", 15, 300, "USD", constants.Sell},
		}

		costBasisBySymbol := CostBasisBySymbol(trades, untilDate)
		got, ok := costBasisBySymbol["XYZ"]
		var expected float64 = 550
		costBasis := got.CostBasis

		if !ok {
			t.Errorf("not ok")
		}

		if costBasis != expected {
			t.Errorf("got %f, want %f", costBasis, expected)
		}
	})

	t.Run("Multiple buy and sell orders should use FIFO", func(t *testing.T) {
		var trades []model.Trade = []model.Trade{
			{0, "XYZ", "2024-01-01", 10, 100, "USD", constants.Buy},
			{0, "XYZ", "2024-02-01", 20, 100, "USD", constants.Buy},
			{0, "XYZ", "2024-03-01", 5, 300, "USD", constants.Sell},
			{0, "XYZ", "2024-04-01", 10, 100, "USD", constants.Buy},
			{0, "XYZ", "2024-05-01", 20, 100, "USD", constants.Buy},
			{0, "XYZ", "2024-06-01", 10, 300, "USD", constants.Sell},
		}

		var expected float64 = 15*100 + 10*100 + 20*100
		costBasisBySymbol := CostBasisBySymbol(trades, untilDate)
		got, ok := costBasisBySymbol["XYZ"]
		costBasis := got.CostBasis

		if !ok {
			t.Errorf("not ok")
		}

		if costBasis != expected {
			t.Errorf("got %f, want %f", costBasis, expected)
		}
	})

	t.Run("First order is sell order should not be counted", func(t *testing.T) {
		var trades []model.Trade = []model.Trade{
			{0, "XYZ", "2024-01-01", 10, 100, "USD", constants.Sell},
			{0, "XYZ", "2024-02-01", 20, 100, "USD", constants.Buy},
		}

		var expected float64 = 20 * 100
		costBasisBySymbol := CostBasisBySymbol(trades, untilDate)
		got, ok := costBasisBySymbol["XYZ"]
		costBasis := got.CostBasis

		if !ok {
			t.Errorf("not ok")
		}

		if costBasis != expected {
			t.Errorf("got %f, want %f", costBasis, expected)
		}
	})

	t.Run("All sell orders", func(t *testing.T) {
		var trades []model.Trade = []model.Trade{
			{0, "XYZ", "2024-01-01", 10, 100, "USD", constants.Sell},
			{0, "XYZ", "2024-02-01", 20, 100, "USD", constants.Sell},
		}

		var expected float64 = 0
		costBasisBySymbol := CostBasisBySymbol(trades, untilDate)
		got, ok := costBasisBySymbol["XYZ"]
		costBasis := got.CostBasis

		if !ok {
			t.Errorf("not ok")
		}

		if costBasis != expected {
			t.Errorf("got %f, want %f", costBasis, expected)
		}
	})

	t.Run("Two symbols", func(t *testing.T) {
		var trades []model.Trade = []model.Trade{
			{0, "XYZ", "2024/01/01", 10, 100, "USD", constants.Buy},
			{0, "XYZ", "2024/02/01", 20, 100, "USD", constants.Buy},
			{0, "XYZ", "2024/03/01", 5, 300, "USD", constants.Sell},
			{0, "GGG", "2024/01/01", 10, 100, "USD", constants.Buy},
			{0, "GGG", "2024/02/01", 20, 100, "USD", constants.Buy},
			{0, "GGG", "2024/03/01", 5, 300, "USD", constants.Sell},
		}

		costBasisBySymbol := CostBasisBySymbol(trades, untilDate)

		testCases := []struct {
			got  float64
			want float64
		}{
			{got: costBasisBySymbol["XYZ"].CostBasis, want: 2500},
			{got: costBasisBySymbol["GGG"].CostBasis, want: 2500},
		}

		for _, tc := range testCases {
			if tc.got != tc.want {
				t.Errorf("got %f, want %f", tc.got, tc.want)
			}
		}
	})

	t.Run("Trades outside of date range", func(t *testing.T) {
		var trades []model.Trade = []model.Trade{
			{0, "XYZ", "2024-01-01", 10, 100, "USD", constants.Buy},
			{0, "XYZ", "2024-02-01", 20, 100, "USD", constants.Buy},
			{0, "XYZ", "2024-03-01", 5, 300, "USD", constants.Sell},
			{0, "XYZ", "2024-04-01", 10, 100, "USD", constants.Buy},
			{0, "XYZ", "2024-05-01", 20, 100, "USD", constants.Buy},
			{0, "XYZ", "2024-06-01", 10, 300, "USD", constants.Sell},
		}

		until := time.Date(2024, 5, 2, 0, 0, 0, 0, time.UTC)

		costBasisBySymbol := CostBasisBySymbol(trades, until)

		got, ok := costBasisBySymbol["XYZ"]
		costBasis := got.CostBasis
		var expected float64 = 15*100 + 10*100 + 30*100

		if !ok {
			t.Errorf("not ok")
		}

		if costBasis != expected {
			t.Errorf("got %f, want %f", costBasis, expected)
		}
	})
}
