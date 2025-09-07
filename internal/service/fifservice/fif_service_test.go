package fifservice

import (
	"slices"
	"testing"
	"time"
)

type MockFIFRepository struct {
}

func (r MockFIFRepository) GetHoldingQuantities(holdingsIDs []HoldingID, upUntil time.Time) map[HoldingID]float64 {
	return map[HoldingID]float64{
		0: 200,
	}
}

func TestFDRIncome(t *testing.T) {
	service := NewFIFService(MockFIFRepository{})

	start, _ := time.Parse(time.DateOnly, "2021-04-01")
	end, _ := time.Parse(time.DateOnly, "2022-03-31")
	input := FDRInput{
		Holdings: []FDRHoldingInput{
			FDRHoldingInput{
				OpeningPrice:      100,
				ExchangeRateToNZD: 1.5,
				HoldingID:         0,
			},
		},
	}

	t.Run("test FDR income - no sales throughout period", func(t *testing.T) {
		got := service.FDRDIncome(input, start, end)
		want := FDRResult{
			Holdings: []FDRHoldingResult{
				FDRHoldingResult{
					Name:                "",
					Symbol:              "",
					OpeningValue:        30000,
					QuickSaleAdjustment: 0,
					Income:              1500,
				},
			},
		}

		if !slices.Equal(got.Holdings, want.Holdings) {
			t.Errorf("got %v, want %v", got.Holdings, want.Holdings)
		}
	})
}
