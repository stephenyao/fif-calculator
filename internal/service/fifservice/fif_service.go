package fifservice

import (
	"time"
)

type FIFService interface {
	FDRDIncome(input FDRInput, startDate time.Time, endDate time.Time) FDRResult
}

type HoldingID int

type FIFRepository interface {
	GetHoldingQuantities(holdingsIDs []HoldingID, upUntil time.Time) map[HoldingID]float64
}

type FIFCalculationService struct {
	repository FIFRepository
}

type FDRResult struct {
	Holdings []FDRHoldingResult
}

type FDRHoldingResult struct {
	Name                string
	Symbol              string
	OpeningValue        float64
	QuickSaleAdjustment float64
	Income              float64
}

type FDRInput struct {
	Holdings []FDRHoldingInput
}

type FDRHoldingInput struct {
	OpeningPrice      float64
	ExchangeRateToNZD float64
	HoldingID         HoldingID
}

func NewFIFService(repo FIFRepository) FIFService {
	return FIFCalculationService{
		repository: repo,
	}
}

func (s FIFCalculationService) FDRDIncome(input FDRInput, startDate time.Time, endDate time.Time) FDRResult {
	// 1. Get all the holding IDs that need to be computed
	holdingIDs := make([]HoldingID, len(input.Holdings))
	for _, holding := range input.Holdings {
		holdingIDs = append(holdingIDs, holding.HoldingID)
	}

	// 2. Get the quantity for each holding at the start date

	holdingQuantities := s.repository.GetHoldingQuantities(holdingIDs, startDate)

	// 3. For each holding calculate the FDR (5% * opening market value)

	result := FDRResult{}
	for _, holding := range input.Holdings {
		quantity := holdingQuantities[holding.HoldingID]
		openingValue := quantity * holding.OpeningPrice * holding.ExchangeRateToNZD
		fdrRate := 0.05

		r := FDRHoldingResult{
			Name:                "",
			Symbol:              "",
			OpeningValue:        openingValue,
			QuickSaleAdjustment: 0,
			Income:              openingValue * fdrRate,
		}
		result.Holdings = append(result.Holdings, r)
	}

	return result
}
