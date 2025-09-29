package fifservice

import (
	. "fif-calculator/internal/repository"
	"reflect"
	"slices"
	"testing"
	"time"
)

type MockFIFRepository struct {
	holdingQuantities map[HoldingID]FIFHoldingQuantity
	tradeActivities   map[HoldingID][]FIFTradeActivity
}

func (r MockFIFRepository) GetHoldingQuantities(holdingsIDs []HoldingID, upUntil time.Time) (map[HoldingID]FIFHoldingQuantity, error) {
	return r.holdingQuantities, nil
}

func (r MockFIFRepository) GetTrades(holdingsIDs []HoldingID, start, end time.Time) (map[HoldingID][]FIFTradeActivity, error) {
	return r.tradeActivities, nil
}

func TestQuickSaleAdjustment(t *testing.T) {
	service := NewFIFService(MockFIFRepository{})

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

	t.Run("no peak differential or real gain", func(t *testing.T) {
		peakDiff := PeakDifferentialResult{}
		realGain := RealGainResult{}

		got := service.QuickSaleAdjustment(peakDiff, realGain)
		want := QuickSaleAdjustmentResult{}

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
		service := NewFIFService(MockFIFRepository{
			holdingQuantities: holdingQuantities(),
		})
		got := service.FDRIncome(input, start, end)
		want := FDRResult{
			Holdings: []FDRHoldingResult{
				FDRHoldingResult{
					Name:                "Google",
					Symbol:              "GOOG",
					OpeningValue:        30000,
					QuickSaleAdjustment: QuickSaleAdjustmentResult{},
					Income:              1500,
				},
				{
					Name:                "Block",
					Symbol:              "XYZ",
					OpeningValue:        60000,
					QuickSaleAdjustment: QuickSaleAdjustmentResult{},
					Income:              3000,
				},
			},
			TotalIncome: 4500,
			Year:        2021,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("fif repository returns no holdings", func(t *testing.T) {
		service := NewFIFService(MockFIFRepository{})

		got := service.FDRIncome(input, start, end)
		want := FDRResult{
			Year: 2021,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("with quick sales adjustment", func(t *testing.T) {
		service := NewFIFService(MockFIFRepository{
			holdingQuantities: holdingWithQuickSales(),
			tradeActivities:   tradeActivitiesWithQuickSales(),
		})

		input := FDRInput{
			Holdings: []FDRHoldingInput{
				{
					OpeningPrice:      20,
					ExchangeRateToNZD: 1,
					HoldingID:         0,
				},
			},
		}

		got := service.FDRIncome(input, start, end)
		want := FDRResult{
			Holdings: []FDRHoldingResult{
				{
					Name:         "Google",
					Symbol:       "GOOG",
					OpeningValue: 200000,
					QuickSaleAdjustment: QuickSaleAdjustmentResult{
						PeakDifferentialResult: PeakDifferentialResult{
							PeakQuantity:     15000,
							PeakDifferential: 2000,
							QuantityStart:    10000,
							QuantityEnd:      13000,
							AverageCost:      22,
							Result:           2200,
							TimeStamp:        time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
						},
						RealGainResult: RealGainResult{
							Sales: []GainOnSale{
								{
									Quantity:          4000,
									Gain:              12000,
									CostOfAcquisition: 88000,
								},
							},
							Result: 12000,
						},
						Result: 2200,
					},
					Income: 12200,
				},
			},
			TotalIncome: 12200,
			Year:        2021,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestPeakDifferential(t *testing.T) {
	service := NewFIFService(MockFIFRepository{})

	testCases := []struct {
		name     string
		quantity FIFHoldingQuantity
		trades   []FIFTradeActivity
		want     PeakDifferentialResult
	}{
		{
			name: "buy and sell activity",
			quantity: FIFHoldingQuantity{
				Quantity: 10000,
				Name:     "Block",
				Symbol:   "XYZ",
			}, trades: buySellActivities(),
			want: PeakDifferentialResult{
				PeakQuantity:     15000,
				PeakDifferential: 2000,
				QuantityStart:    10000,
				QuantityEnd:      13000,
				AverageCost:      22,
				Result:           2200,
				TimeStamp:        time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			},
		}, {
			name: "no trade activity throughout period",
			quantity: FIFHoldingQuantity{
				Quantity: 100,
				Name:     "Block",
				Symbol:   "XYZ",
			},
			trades: []FIFTradeActivity{},
			want:   PeakDifferentialResult{},
		},
		{
			name: "no buy activities",
			quantity: FIFHoldingQuantity{
				Quantity: 10000,
				Name:     "Block",
				Symbol:   "XYZ",
			},
			trades: noBuyActivities(),
			want:   PeakDifferentialResult{},
		},
		{
			name: "no sell activities",
			quantity: FIFHoldingQuantity{
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
		trades []FIFTradeActivity
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
			want:   RealGainResult{},
		}, {
			name:   "multiple buy and sells should drain first buy in the queue",
			trades: buyQuantityGreaterThanSell(),
			want: RealGainResult{
				Sales: []GainOnSale{
					{
						Quantity:          3000,
						Gain:              9000,
						CostOfAcquisition: 66000,
					},
					{
						Quantity:          3000,
						Gain:              6000,
						CostOfAcquisition: 69000,
					},
				},
				Result: 15000,
			},
		}, {
			name:   "no sell activities",
			trades: noSellActivities(),
			want:   RealGainResult{},
		}, {
			name:   "no buy activities",
			trades: noBuyActivities(),
			want:   RealGainResult{},
		}, {
			name:   "sell at loss",
			trades: sellAtLossActivities(),
			want:   RealGainResult{},
		},
		{
			name:   "sell at loss then profit",
			trades: sellAtLossThenProfit(),
			want: RealGainResult{
				Sales: []GainOnSale{
					{
						Quantity:          4000,
						Gain:              112000,
						CostOfAcquisition: 88000,
					},
				},
				Result: 112000,
			},
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

func buySellActivities() []FIFTradeActivity {
	return []FIFTradeActivity{
		{
			Date:         time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			Action:       "buy",
			Quantity:     5000,
			Price:        22,
			ExchangeRate: 1,
		},
		{
			Date:         time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
			Action:       "sell",
			Quantity:     4000,
			Price:        25,
			ExchangeRate: 1,
		},
		{
			Date:         time.Date(2024, 12, 23, 0, 0, 0, 0, time.UTC),
			Action:       "buy",
			Quantity:     2000,
			Price:        22,
			ExchangeRate: 1,
		},
	}
}

func sellQuantityGreaterThanBuy() []FIFTradeActivity {
	return []FIFTradeActivity{
		{
			Date:         time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			Action:       "buy",
			Quantity:     5000,
			Price:        22,
			ExchangeRate: 1,
		},
		{
			Date:         time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
			Action:       "buy",
			Quantity:     4000,
			Price:        25,
			ExchangeRate: 1,
		},
		{
			Date:         time.Date(2024, 12, 23, 0, 0, 0, 0, time.UTC),
			Action:       "sell",
			Quantity:     10000,
			Price:        22,
			ExchangeRate: 1,
		},
	}
}

func buyQuantityGreaterThanSell() []FIFTradeActivity {
	return []FIFTradeActivity{
		{
			Date:         time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			Action:       "buy",
			Quantity:     5000,
			Price:        22,
			ExchangeRate: 1,
		},
		{
			Date:         time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
			Action:       "buy",
			Quantity:     4000,
			Price:        25,
			ExchangeRate: 1,
		},
		{
			Date:         time.Date(2024, 12, 23, 0, 0, 0, 0, time.UTC),
			Action:       "sell",
			Quantity:     3000,
			Price:        25,
			ExchangeRate: 1,
		},
		{
			Date:         time.Date(2024, 12, 23, 0, 0, 0, 0, time.UTC),
			Action:       "sell",
			Quantity:     3000,
			Price:        25,
			ExchangeRate: 1,
		},
	}
}

func noBuyActivities() []FIFTradeActivity {
	return []FIFTradeActivity{
		{
			Date:         time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			Action:       "sell",
			Quantity:     5000,
			Price:        22,
			ExchangeRate: 1,
		},
		{
			Date:         time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
			Action:       "sell",
			Quantity:     4000,
			Price:        25,
			ExchangeRate: 1,
		},
	}
}

func sellAtLossActivities() []FIFTradeActivity {
	return []FIFTradeActivity{
		{
			Date:         time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			Action:       "buy",
			Quantity:     5000,
			Price:        22,
			ExchangeRate: 1,
		},
		{
			Date:         time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
			Action:       "sell",
			Quantity:     4000,
			Price:        10,
			ExchangeRate: 1,
		},
		{
			Date:         time.Date(2024, 12, 10, 0, 0, 0, 0, time.UTC),
			Action:       "buy",
			Quantity:     5000,
			Price:        22,
			ExchangeRate: 1,
		},
		{
			Date:         time.Date(2024, 12, 22, 0, 0, 0, 0, time.UTC),
			Action:       "sell",
			Quantity:     4000,
			Price:        10,
			ExchangeRate: 1,
		},
	}
}

func sellAtLossThenProfit() []FIFTradeActivity {
	return []FIFTradeActivity{
		{
			Date:         time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			Action:       "buy",
			Quantity:     5000,
			Price:        22,
			ExchangeRate: 1,
		},
		{
			Date:         time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
			Action:       "sell",
			Quantity:     4000,
			Price:        10,
			ExchangeRate: 1,
		},
		{
			Date:         time.Date(2024, 12, 10, 0, 0, 0, 0, time.UTC),
			Action:       "buy",
			Quantity:     5000,
			Price:        22,
			ExchangeRate: 1,
		},
		{
			Date:         time.Date(2024, 12, 22, 0, 0, 0, 0, time.UTC),
			Action:       "sell",
			Quantity:     4000,
			Price:        50,
			ExchangeRate: 1,
		},
	}
}

func noSellActivities() []FIFTradeActivity {
	return []FIFTradeActivity{
		{
			Date:         time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			Action:       "buy",
			Quantity:     5000,
			Price:        22,
			ExchangeRate: 1,
		},
		{
			Date:         time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
			Action:       "buy",
			Quantity:     4000,
			Price:        25,
			ExchangeRate: 1,
		},
	}
}

func holdingQuantities() map[HoldingID]FIFHoldingQuantity {
	return map[HoldingID]FIFHoldingQuantity{
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

func holdingWithQuickSales() map[HoldingID]FIFHoldingQuantity {
	return map[HoldingID]FIFHoldingQuantity{
		0: {
			Quantity: 10000,
			Name:     "Google",
			Symbol:   "GOOG",
		},
	}
}

func tradeActivitiesWithQuickSales() map[HoldingID][]FIFTradeActivity {
	return map[HoldingID][]FIFTradeActivity{
		0: buySellActivities(),
	}
}
