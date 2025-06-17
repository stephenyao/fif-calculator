package fifservice

import (
	"fif-calculator/internal/constants"
	"fif-calculator/internal/model"
	"testing"
	"time"
)

func mustDate(s string) time.Time {
	d, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return d
}

func TestComputeHoldingsBetween(t *testing.T) {
	start := mustDate("2024-04-01")
	end := mustDate("2025-03-31")

	t.Run("only buy trades before start and end", func(t *testing.T) {
		trades := []model.Trade{
			{Symbol: "XYZ", BuyDate: "2024-03-01", Quantity: 10, Action: constants.Buy},
			{Symbol: "XYZ", BuyDate: "2024-03-15", Quantity: 5, Action: constants.Buy},
		}

		holdings, err := ComputeHoldingsBetween(trades, start, end)
		if err != nil {
			t.Fatal(err)
		}

		if len(holdings) != 1 {
			t.Fatalf("expected 1 holding, got %d", len(holdings))
		}
		h := holdings[0]
		if h.QuantityStart != 15 || h.QuantityEnd != 15 {
			t.Errorf("got start %f, end %f; want 15.0 both", h.QuantityStart, h.QuantityEnd)
		}

		if h.NumberOfTrades != 0 {
			t.Errorf("got number of trades; want 0, got %d", h.NumberOfTrades)
		}
	})

	t.Run("buy and sell before and after start", func(t *testing.T) {
		trades := []model.Trade{
			{Symbol: "XYZ", BuyDate: "2024-03-01", Quantity: 10, Action: constants.Buy},
			{Symbol: "XYZ", BuyDate: "2024-04-10", Quantity: 5, Action: constants.Buy},
			{Symbol: "XYZ", BuyDate: "2024-05-01", Quantity: 4, Action: constants.Sell},
		}

		holdings, err := ComputeHoldingsBetween(trades, start, end)
		if err != nil {
			t.Fatal(err)
		}

		if len(holdings) != 1 {
			t.Fatalf("expected 1 holding, got %d", len(holdings))
		}
		h := holdings[0]
		if h.QuantityStart != 10 {
			t.Errorf("got start %f, want 10", h.QuantityStart)
		}
		if h.QuantityEnd != 11 {
			t.Errorf("got end %f, want 11", h.QuantityEnd)
		}
		if h.NumberOfTrades != 2 {
			t.Errorf("got number of trades; want 2, got %d", h.NumberOfTrades)
		}
	})

	t.Run("negative quantity should be skipped", func(t *testing.T) {
		trades := []model.Trade{
			{Symbol: "XYZ", BuyDate: "2024-03-01", Quantity: 5, Action: constants.Sell},
		}

		holdings, err := ComputeHoldingsBetween(trades, start, end)
		if err != nil {
			t.Fatal(err)
		}
		if len(holdings) != 0 {
			t.Errorf("expected 0 holdings due to negative quantity, got %d", len(holdings))
		}
	})

	t.Run("empty trade list", func(t *testing.T) {
		holdings, err := ComputeHoldingsBetween([]model.Trade{}, start, end)
		if err != nil {
			t.Fatal(err)
		}
		if len(holdings) != 0 {
			t.Errorf("expected 0 holdings, got %d", len(holdings))
		}
	})

	t.Run("multiple symbols", func(t *testing.T) {
		trades := []model.Trade{
			{Symbol: "ABC", BuyDate: "2024-03-01", Quantity: 10, Action: constants.Buy},
			{Symbol: "XYZ", BuyDate: "2024-02-01", Quantity: 5, Action: constants.Buy},
			{Symbol: "XYZ", BuyDate: "2024-06-01", Quantity: 10, Action: constants.Buy},
		}

		holdings, err := ComputeHoldingsBetween(trades, start, end)
		if err != nil {
			t.Fatal(err)
		}

		if len(holdings) != 2 {
			t.Errorf("expected 2 holdings, got %d", len(holdings))
		}

		tests := []struct {
			symbol        string
			wantStart     float64
			wantEnd       float64
			wantNumTrades int
		}{
			{"ABC", 10, 10, 0},
			{"XYZ", 5, 15, 1},
		}

		holdingMap := make(map[string]*model.HoldingQuantity)
		for _, h := range holdings {
			holdingMap[h.Symbol] = h
		}

		for _, tt := range tests {
			t.Run(tt.symbol, func(t *testing.T) {
				h, ok := holdingMap[tt.symbol]
				if !ok {
					t.Fatalf("missing holding for symbol %s", tt.symbol)
				}
				if h.QuantityStart != tt.wantStart {
					t.Errorf("start qty for %s: got %.2f, want %.2f", tt.symbol, h.QuantityStart, tt.wantStart)
				}
				if h.QuantityEnd != tt.wantEnd {
					t.Errorf("end qty for %s: got %.2f, want %.2f", tt.symbol, h.QuantityEnd, tt.wantEnd)
				}
				if h.NumberOfTrades != tt.wantNumTrades {
					t.Errorf("got number of trades; want %d, got %d", tt.wantNumTrades, h.NumberOfTrades)
				}
			})
		}
	})

	t.Run("multiple symbols where one has 0 trades and no start quantity in period", func(t *testing.T) {
		trades := []model.Trade{
			{Symbol: "ABC", BuyDate: "2024-03-01", Quantity: 10, Action: constants.Buy},
			{Symbol: "ABC", BuyDate: "2024-03-02", Quantity: 10, Action: constants.Sell},
			{Symbol: "XYZ", BuyDate: "2024-02-01", Quantity: 5, Action: constants.Buy},
			{Symbol: "XYZ", BuyDate: "2024-06-01", Quantity: 10, Action: constants.Buy},
		}

		holdings, err := ComputeHoldingsBetween(trades, start, end)
		if err != nil {
			t.Fatal(err)
		}
		want := 1
		if len(holdings) != want {
			t.Errorf("got %d holdings, want %d", len(holdings), want)
		}
	})

	t.Run("multiple symbols where one has has trades in period but sold all", func(t *testing.T) {
		trades := []model.Trade{
			{Symbol: "ABC", BuyDate: "2024-05-01", Quantity: 10, Action: constants.Buy},
			{Symbol: "ABC", BuyDate: "2024-05-02", Quantity: 10, Action: constants.Sell},
			{Symbol: "XYZ", BuyDate: "2024-02-01", Quantity: 5, Action: constants.Buy},
			{Symbol: "XYZ", BuyDate: "2024-06-01", Quantity: 10, Action: constants.Buy},
		}

		holdings, err := ComputeHoldingsBetween(trades, start, end)
		if err != nil {
			t.Fatal(err)
		}
		want := 2
		if len(holdings) != want {
			t.Errorf("got %d holdings, want %d", len(holdings), want)
		}
	})
}
