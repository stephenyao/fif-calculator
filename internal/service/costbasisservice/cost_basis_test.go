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

func TestMaxCostBasisDuringYear(t *testing.T) {
	end := time.Date(2025, 3, 31, 0, 0, 0, 0, time.UTC)

	t.Run("never exceeds 50000", func(t *testing.T) {
		trades := []model.Trade{
			{Symbol: "XYZ", BuyDate: "2024-04-15", Quantity: 10, Price: 1000, Action: constants.Buy}, // $10,000
			{Symbol: "XYZ", BuyDate: "2024-06-01", Quantity: 20, Price: 1000, Action: constants.Buy}, // $20,000
			{Symbol: "XYZ", BuyDate: "2024-07-01", Quantity: 10, Price: 1000, Action: constants.Sell},
		}

		got := MaxCostBasisDuringYear(trades, end)
		want := 30000.0
		if got != want {
			t.Errorf("got %f, want %f", got, want)
		}
	})

	t.Run("exceeds 50000 on a single day", func(t *testing.T) {
		trades := []model.Trade{
			{Symbol: "XYZ", BuyDate: "2024-04-10", Quantity: 30, Price: 1000, Action: constants.Buy}, // $30,000
			{Symbol: "XYZ", BuyDate: "2024-04-11", Quantity: 30, Price: 1000, Action: constants.Buy}, // $60,000
			{Symbol: "XYZ", BuyDate: "2024-04-12", Quantity: 40, Price: 1000, Action: constants.Sell},
		}

		got := MaxCostBasisDuringYear(trades, end)
		want := 60000.0
		if got != want {
			t.Errorf("got %f, want %f", got, want)
		}
	})

	t.Run("crosses 50000 then sells down", func(t *testing.T) {
		trades := []model.Trade{
			{Symbol: "XYZ", BuyDate: "2024-05-01", Quantity: 60, Price: 1000, Action: constants.Buy},  // $60,000
			{Symbol: "XYZ", BuyDate: "2024-06-01", Quantity: 20, Price: 1000, Action: constants.Sell}, // $40,000
		}

		got := MaxCostBasisDuringYear(trades, end)
		want := 60000.0
		if got != want {
			t.Errorf("got %f, want %f", got, want)
		}
	})

	t.Run("buys and sells to hover below 50000", func(t *testing.T) {
		trades := []model.Trade{
			{Symbol: "XYZ", BuyDate: "2024-04-05", Quantity: 25, Price: 1000, Action: constants.Buy},  // $25,000
			{Symbol: "XYZ", BuyDate: "2024-05-01", Quantity: 10, Price: 1000, Action: constants.Buy},  // $35,000
			{Symbol: "XYZ", BuyDate: "2024-06-01", Quantity: 5, Price: 1000, Action: constants.Sell},  // $30,000
			{Symbol: "XYZ", BuyDate: "2024-07-01", Quantity: 15, Price: 1000, Action: constants.Buy},  // $45,000
			{Symbol: "XYZ", BuyDate: "2024-08-01", Quantity: 10, Price: 1000, Action: constants.Sell}, // $35,000
		}

		got := MaxCostBasisDuringYear(trades, end)
		want := 45000.0
		if got != want {
			t.Errorf("got %f, want %f", got, want)
		}
	})

	t.Run("trades outside of date range are ignored", func(t *testing.T) {
		trades := []model.Trade{
			{Symbol: "XYZ", BuyDate: "2023-12-01", Quantity: 60, Price: 1000, Action: constants.Buy}, // Ignored
			{Symbol: "XYZ", BuyDate: "2025-04-01", Quantity: 60, Price: 1000, Action: constants.Buy}, // Ignored
		}

		got := MaxCostBasisDuringYear(trades, end)
		var want float64 = 60000
		if got != want {
			t.Errorf("got %f, want %f", got, want)
		}
	})
}
