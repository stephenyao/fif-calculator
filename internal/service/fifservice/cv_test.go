package fifservice

import (
	"fif-calculator/internal/model"
	"testing"
	"time"
)

func TestComputeCVIncome(t *testing.T) {
	trades := []model.Trade{
		{0, "XYZ", "2020-01-01", 10000, 20, "USD", "buy"},
		{0, "XYZ", "2021-10-01", 5000, 22, "USD", "buy"},
		{0, "XYZ", "2021-12-01", 4000, 10, "USD", "sell"},
		{0, "XYZ", "2021-12-23", 2000, 22, "USD", "buy"},
	}

	holdings := []*model.HoldingInfo{
		{Symbol: "XYZ", QuantityStart: 10000, QuantityEnd: 13000, OpeningPrice: 10, ClosingPrice: 25},
	}

	start := time.Date(2021, time.April, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2022, time.March, 31, 0, 0, 0, 0, time.UTC)
	gainLoss := model.GainLossParams{
		Dividends:        5000,
		TaxCredits:       1000,
		OtherGains:       0,
		ForeignIncomeTax: 2000,
		OtherCosts:       0,
	}

	holdings[0].GainLoss = gainLoss

	results, err := ComputeCVIncome(trades, holdings, start, end)

	if err != nil {
		t.Errorf("ComputeCVIncome failed: %v", err)
	}

	got := results[0]
	want := model.CVResult{
		Symbol:       "XYZ",
		OpeningValue: 100000,
		ClosingValue: 325000,
		Gains:        46000,
		Costs:        156000,
	}

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	var expectedIncome float64 = 115000
	gotIncome := got.TotalIncome()

	if gotIncome != expectedIncome {
		t.Errorf("got %v, want %v", gotIncome, expectedIncome)
	}
}
