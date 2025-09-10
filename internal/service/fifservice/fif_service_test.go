package fifservice

import (
	"slices"
	"testing"
	"time"
)

type MockFIFRepository struct {
	returnsEmpty bool
}

func (r MockFIFRepository) GetHoldingQuantities(holdingsIDs []HoldingID, upUntil time.Time) map[HoldingID]FDRHoldingQuantity {
	if r.returnsEmpty {
		return make(map[HoldingID]FDRHoldingQuantity)
	}

	return map[HoldingID]FDRHoldingQuantity{
		0: {
			Quantity: 200,
			Name:     "Google",
			Symbol:   "GOOG",
		},
		1: {
			Quantity: 200,
			Name:     "Block",
			Symbol:   "XYZ",
		},
	}
}

func (r MockFIFRepository) GetTrades(holdingsIDs []HoldingID, start, end time.Time) map[HoldingID][]FDRTradeActivity {
	if r.returnsEmpty {
		return make(map[HoldingID][]FDRTradeActivity)
	}
	return make(map[HoldingID][]FDRTradeActivity)
}

func TestFDRIncome(t *testing.T) {
	start, _ := time.Parse(time.DateOnly, "2021-04-01")
	end, _ := time.Parse(time.DateOnly, "2022-03-31")
	input := FDRInput{
		Holdings: []FDRHoldingInput{
			{
				OpeningPrice:      100,
				ExchangeRateToNZD: 1.5,
				HoldingID:         0,
			},
			{
				OpeningPrice:      200,
				ExchangeRateToNZD: 1.5,
				HoldingID:         1,
			},
		},
	}

	t.Run("test FDR income - no sales throughout period", func(t *testing.T) {
		service := NewFIFService(MockFIFRepository{})
		got := service.FDRIncome(input, start, end)
		want := FDRResult{
			Holdings: []FDRHoldingResult{
				FDRHoldingResult{
					Name:                "Google",
					Symbol:              "GOOG",
					OpeningValue:        30000,
					QuickSaleAdjustment: 0,
					Income:              1500,
				},
				{
					Name:                "Block",
					Symbol:              "XYZ",
					OpeningValue:        60000,
					QuickSaleAdjustment: 0,
					Income:              3000,
				},
			},
		}

		if !slices.Equal(got.Holdings, want.Holdings) {
			t.Errorf("got %v, want %v", got.Holdings, want.Holdings)
		}
	})

	t.Run("fif repository returns no holdings", func(t *testing.T) {
		service := NewFIFService(MockFIFRepository{
			returnsEmpty: true,
		})

		got := service.FDRIncome(input, start, end)
		want := FDRResult{}

		if !slices.Equal(got.Holdings, want.Holdings) {
			t.Errorf("got %v, want %v", got.Holdings, want.Holdings)
		}
	})
}

func TestPeakDifferential(t *testing.T) {
	service := NewFIFService(MockFIFRepository{})

	t.Run("no trade activity throughout period", func(t *testing.T) {
		got := service.PeakHoldingDifferential(FDRHoldingQuantity{
			Quantity: 100,
			Name:     "Block",
			Symbol:   "XYZ",
		}, []FDRTradeActivity{})
		want := PeakDifferentialResult{}
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("there was a buy activity in the period", func(t *testing.T) {
		holdingInfo := FDRHoldingQuantity{
			Quantity: 10000,
			Name:     "Block",
			Symbol:   "XYZ",
		}
		trades := []FDRTradeActivity{
			{
				Date:         time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
				Action:       "buy",
				Quantity:     5000,
				Price:        22,
				ExchangeRate: 1,
				AmountInNZD:  110000,
			},
			{
				Date:         time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
				Action:       "sell",
				Quantity:     4000,
				Price:        25,
				ExchangeRate: 1,
				AmountInNZD:  100000,
			},
			{
				Date:         time.Date(2024, 12, 23, 0, 0, 0, 0, time.UTC),
				Action:       "buy",
				Quantity:     2000,
				Price:        22,
				ExchangeRate: 1,
				AmountInNZD:  44000,
			},
		}

		got := service.PeakHoldingDifferential(holdingInfo, trades)
		want := PeakDifferentialResult{
			PeakQuantity:  2000,
			QuantityStart: 10000,
			QuantityEnd:   13000,
			AverageCost:   22,
			Result:        2200,
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestRealGain(t *testing.T) {
	service := NewFIFService(MockFIFRepository{})
	trades := []FDRTradeActivity{
		{
			Date:         time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			Action:       "buy",
			Quantity:     5000,
			Price:        22,
			ExchangeRate: 1,
			AmountInNZD:  110000,
		},
		{
			Date:         time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
			Action:       "sell",
			Quantity:     4000,
			Price:        25,
			ExchangeRate: 1,
			AmountInNZD:  100000,
		},
		{
			Date:         time.Date(2024, 12, 23, 0, 0, 0, 0, time.UTC),
			Action:       "buy",
			Quantity:     2000,
			Price:        22,
			ExchangeRate: 1,
			AmountInNZD:  44000,
		},
	}

	got := service.RealGain(trades)
	want := RealGainResult{
		Sales: []GainOnSale{
			{
				Quantity:          4000,
				Gain:              12000,
				CostOfAcquisition: 88000,
			},
		},
		Result: 12000,
	}

	if !slices.Equal(got.Sales, want.Sales) {
		t.Errorf("got %v, want %v", got, want)
	}

	if got.Result != want.Result {
		t.Errorf("got %v, want %v", got.Result, want.Result)
	}
}
