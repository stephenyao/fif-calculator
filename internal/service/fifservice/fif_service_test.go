package fifservice

import (
	"fif-calculator/internal/model"
	"testing"
	"time"
)

func TestComputeFRDIncome(t *testing.T) {
	t.Run("No quick sale adjustment", func(t *testing.T) {
		trades := []model.Trade{
			{0, "XYZ", "2021-01-01", 1000, 100, "USD", "buy"},
			{0, "XYZ", "2021-02-01", 1000, 100, "USD", "buy"},
			{0, "GOOG", "2021-02-01", 1500, 100, "USD", "buy"},
			{0, "GOOG", "2021-02-01", 1500, 100, "USD", "buy"},
		}

		start := time.Date(2022, time.April, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2023, time.March, 31, 0, 0, 0, 0, time.UTC)

		holdings := []model.HoldingQuantity{
			{Symbol: "XYZ", QuantityStart: 2000, QuantityEnd: 2000, OpeningPrice: 100, ClosingPrice: 50},
			{Symbol: "GOOG", QuantityStart: 3000, QuantityEnd: 3000, OpeningPrice: 100, ClosingPrice: 50},
		}

		got, err := ComputeFRDIncome(trades, holdings, start, end)

		var want = 200000*0.05 + 300000*0.05

		if err != nil {
			t.Errorf("ComputeFRDIncome() error = %v", err)
		}

		if got != want {
			t.Errorf("got %f, want %f", got, want)
		}
	})

	t.Run("With quick sale adjustment", func(t *testing.T) {
		trades := []model.Trade{
			{0, "XYZ", "2020-01-01", 10000, 20, "USD", "buy"},
			{0, "XYZ", "2021-10-01", 5000, 22, "USD", "buy"},
			{0, "XYZ", "2021-12-01", 4000, 25, "USD", "sell"},
			{0, "XYZ", "2021-12-23", 2000, 22, "USD", "buy"},
		}

		holdings := []model.HoldingQuantity{
			{Symbol: "XYZ", QuantityStart: 10000, QuantityEnd: 13000, OpeningPrice: 20, ClosingPrice: 50},
		}

		start := time.Date(2021, time.April, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2022, time.March, 31, 0, 0, 0, 0, time.UTC)

		got, err := ComputeFRDIncome(trades, holdings, start, end)
		var want float64 = 12200

		if err != nil {
			t.Errorf("ComputeFRDIncome() error = %v", err)
		}

		if got != want {
			t.Errorf("ComputeFRDIncome() go = %v, want %v", got, want)
		}
	})
}

func TestTradesBySymbol(t *testing.T) {
	trades := []model.Trade{
		{0, "XYZ", "2021-01-01", 1000, 100, "USD", "buy"},
		{1, "XYZ", "2021-01-01", 1000, 100, "USD", "buy"},
		{2, "GOOG", "2021-01-01", 1000, 100, "USD", "sell"},
	}

	got := tradesBySymbol(trades)
	expectedXYZCount := 2
	expectedGOOGCount := 1

	if len(got["XYZ"]) != expectedXYZCount {
		t.Errorf("got len(got[XYZ]) = %v, want %v", len(got["XYZ"]), expectedXYZCount)
	}

	if len(got["GOOG"]) != expectedGOOGCount {
		t.Errorf("got len(got[GOOG]) = %v, want %v", len(got["GOOG"]), expectedGOOGCount)
	}
}

func TestCalculateRealGain(t *testing.T) {
	t.Run("testing real gains without draining the previous trade", func(t *testing.T) {
		trades := []model.Trade{
			{0, "XYZ", "2020-01-01", 10000, 20, "USD", "buy"},
			{0, "XYZ", "2021-10-01", 5000, 22, "USD", "buy"},
			{0, "XYZ", "2021-12-01", 4000, 25, "USD", "sell"},
			{0, "XYZ", "2021-12-23", 2000, 22, "USD", "buy"},
		}

		start := time.Date(2021, time.April, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2022, time.March, 31, 0, 0, 0, 0, time.UTC)

		got, err := calculateRealGainForSymbol(trades, start, end)
		var want float64 = 12000

		if err != nil {
			t.Errorf("calculateRealGain() error = %v", err)
		}

		if got != want {
			t.Errorf("got %f, want %f", got, want)
		}
	})

	t.Run("testing real gains with overruning the previous trade", func(t *testing.T) {
		trades := []model.Trade{
			{0, "XYZ", "2020-01-01", 10000, 20, "USD", "buy"},
			{0, "XYZ", "2021-10-01", 5000, 22, "USD", "buy"},
			{0, "XYZ", "2021-11-23", 2000, 23, "USD", "buy"},
			{0, "XYZ", "2021-12-01", 7000, 25, "USD", "sell"},
		}

		start := time.Date(2021, time.April, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2022, time.March, 31, 0, 0, 0, 0, time.UTC)

		got, err := calculateRealGainForSymbol(trades, start, end)
		var want float64 = 19000

		if err != nil {
			t.Errorf("calculateRealGain() error = %v", err)
		}

		if got != want {
			t.Errorf("got %f, want %f", got, want)
		}
	})
}

func TestPeakDifferential(t *testing.T) {
	t.Run("when peak differential is lesser of peak quantity - starting quantity", func(t *testing.T) {
		trades := []model.Trade{
			{0, "XYZ", "2020-01-01", 10000, 20, "USD", "buy"},
			{0, "XYZ", "2021-04-02", 5000, 22, "USD", "buy"},
			{0, "XYZ", "2021-04-23", 2000, 25, "USD", "buy"},
			{0, "XYZ", "2021-04-25", 7000, 25, "USD", "sell"},
			{0, "XYZ", "2021-05-01", 2000, 25, "USD", "sell"},
			{0, "XYZ", "2021-06-01", 1000, 25, "USD", "sell"},
			{0, "XYZ", "2021-07-01", 1000, 25, "USD", "buy"},
		}

		avgCost := 23.125
		start := time.Date(2021, time.April, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2022, time.March, 31, 0, 0, 0, 0, time.UTC)
		holdings, err := ComputeHoldingsBetween(trades, start, end)

		if err != nil {
			t.Errorf("ComputeHoldingsBetween() error = %v", err)
		}

		got, err := peakDifferentialForSymbol(holdings[0], trades, start, end)
		want := avgCost * 7000

		if err != nil {
			t.Errorf("peakDifferential() error = %v", err)
		}

		if got != want {
			t.Errorf("got %f, want %f", got, want)
		}
	})

	t.Run("when peak differential is lesser of peak quantity - ending quantity", func(t *testing.T) {
		trades := []model.Trade{
			{0, "XYZ", "2020-01-01", 10000, 20, "USD", "buy"},
			{0, "XYZ", "2021-04-02", 5000, 22, "USD", "buy"},
			{0, "XYZ", "2021-04-23", 2000, 25, "USD", "buy"},
			{0, "XYZ", "2021-04-25", 7000, 25, "USD", "sell"},
			{0, "XYZ", "2021-05-01", 2000, 25, "USD", "sell"},
			{0, "XYZ", "2021-06-01", 1000, 25, "USD", "sell"},
			{0, "XYZ", "2021-07-01", 5000, 25, "USD", "buy"},
		}

		avgCost := 23.75
		start := time.Date(2021, time.April, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2022, time.March, 31, 0, 0, 0, 0, time.UTC)
		holdings, err := ComputeHoldingsBetween(trades, start, end)

		got, err := peakDifferentialForSymbol(holdings[0], trades, start, end)
		want := avgCost * 5000

		if err != nil {
			t.Errorf("peakDifferential() error = %v", err)
		}

		if got != want {
			t.Errorf("got %f, want %f", got, want)
		}
	})
}
