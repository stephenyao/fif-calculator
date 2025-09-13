package fifservice

import (
	"reflect"
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

func TestQuickSaleAdjustment(t *testing.T) {
	t.Run("peak differential is smaller than real gain", func(t *testing.T) {
		peakDiff := PeakDifferentialResult{
			PeakQuantity:  10000,
			QuantityStart: 9000,
			QuantityEnd:   8000,
			AverageCost:   25,
			Result:        50000,
		}

		realGain := RealGainResult{
			Sales:  nil,
			Result: 20000,
		}

		service := NewFIFService(MockFIFRepository{})
		got := service.QuickSaleAdjustment(peakDiff, realGain)

		want := QuickSaleAdjustmentResult{
			PeakDifferentialResult: peakDiff,
			RealGainResult:         realGain,
			Result:                 20000,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("peak differential is bigger than real gain", func(t *testing.T) {
		peakDiff := PeakDifferentialResult{
			PeakQuantity:  10000,
			QuantityStart: 9000,
			QuantityEnd:   8000,
			AverageCost:   25,
			Result:        50000,
		}

		realGain := RealGainResult{
			Sales:  nil,
			Result: 80000,
		}

		service := NewFIFService(MockFIFRepository{})
		got := service.QuickSaleAdjustment(peakDiff, realGain)

		want := QuickSaleAdjustmentResult{
			PeakDifferentialResult: peakDiff,
			RealGainResult:         realGain,
			Result:                 50000,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
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

	testCases := []struct {
		name     string
		quantity FDRHoldingQuantity
		trades   []FDRTradeActivity
		want     PeakDifferentialResult
	}{
		{
			name: "buy and sell activity",
			quantity: FDRHoldingQuantity{
				Quantity: 10000,
				Name:     "Block",
				Symbol:   "XYZ",
			}, trades: buySellActivities(),
			want: PeakDifferentialResult{
				PeakQuantity:  2000,
				QuantityStart: 10000,
				QuantityEnd:   13000,
				AverageCost:   22,
				Result:        2200,
			},
		}, {
			name: "no trade activity throughout period",
			quantity: FDRHoldingQuantity{
				Quantity: 100,
				Name:     "Block",
				Symbol:   "XYZ",
			},
			trades: []FDRTradeActivity{},
			want:   PeakDifferentialResult{},
		},
		{
			name: "no buy activities",
			quantity: FDRHoldingQuantity{
				Quantity: 10000,
				Name:     "Block",
				Symbol:   "XYZ",
			},
			trades: noBuyActivities(),
			want:   PeakDifferentialResult{},
		},
		{
			name: "no sell activities",
			quantity: FDRHoldingQuantity{
				Quantity: 10000,
				Name:     "Block",
				Symbol:   "XYZ",
			},
			trades: noSellActivities(),
			want:   PeakDifferentialResult{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := service.PeakHoldingDifferential(tc.quantity, tc.trades)
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestRealGain(t *testing.T) {
	service := NewFIFService(MockFIFRepository{})

	testCases := []struct {
		name   string
		trades []FDRTradeActivity
		want   RealGainResult
	}{
		{
			name:   "buy activity bigger than sell activity",
			trades: buySellActivities(),
			want: RealGainResult{
				Sales: []GainOnSale{
					{
						Quantity:          4000,
						Gain:              12000,
						CostOfAcquisition: 88000,
					},
				},
				Result: 12000,
			},
		}, {
			name:   "buy activity not big enough for sell activity",
			trades: sellQuantityGreaterThanBuy(),
			want: RealGainResult{
				Sales: []GainOnSale{
					{
						Quantity:          10000,
						Gain:              -12000,
						CostOfAcquisition: 210000,
					},
				},
				Result: -12000,
			},
		}, {
			name:   "no sell activities",
			trades: noSellActivities(),
			want:   RealGainResult{},
		}, {
			name:   "no buy activities",
			trades: noBuyActivities(),
			want:   RealGainResult{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := service.RealGain(tc.trades)
			if !slices.Equal(got.Sales, tc.want.Sales) {
				t.Errorf("got %v, want %v", got.Sales, tc.want.Sales)
			}

			if got.Result != tc.want.Result {
				t.Errorf("got %v, want %v", got.Result, tc.want.Result)
			}
		})
	}
}

func buySellActivities() []FDRTradeActivity {
	return []FDRTradeActivity{
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
}

func sellQuantityGreaterThanBuy() []FDRTradeActivity {
	return []FDRTradeActivity{
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
			Action:       "buy",
			Quantity:     4000,
			Price:        25,
			ExchangeRate: 1,
			AmountInNZD:  100000,
		},
		{
			Date:         time.Date(2024, 12, 23, 0, 0, 0, 0, time.UTC),
			Action:       "sell",
			Quantity:     10000,
			Price:        22,
			ExchangeRate: 1,
			AmountInNZD:  220000,
		},
	}
}

func noBuyActivities() []FDRTradeActivity {
	return []FDRTradeActivity{
		{
			Date:         time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			Action:       "sell",
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
	}
}

func noSellActivities() []FDRTradeActivity {
	return []FDRTradeActivity{
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
			Action:       "buy",
			Quantity:     4000,
			Price:        25,
			ExchangeRate: 1,
			AmountInNZD:  100000,
		},
	}
}
