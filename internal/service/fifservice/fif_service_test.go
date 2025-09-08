package fifservice

import (
	"slices"
	"testing"
	"time"
)

type MockFIFRepository struct {
	returnsEmpty bool
}

func (r MockFIFRepository) GetHoldingQuantities(holdingsIDs []HoldingID, upUntil time.Time) map[HoldingID]HoldingFDRInfo {
	if r.returnsEmpty {
		return make(map[HoldingID]HoldingFDRInfo)
	}

	return map[HoldingID]HoldingFDRInfo{
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
