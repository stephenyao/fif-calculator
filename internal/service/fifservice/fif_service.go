package fifservice

import (
	"fif-calculator/internal/constants"
	"time"
)

type FIFService interface {
	FDRIncome(input FDRInput, startDate time.Time, endDate time.Time) FDRResult
	PeakHoldingDifferential(
		holdingInfo FDRHoldingQuantity,
		trades []FDRTradeActivity) PeakDifferentialResult
}

type HoldingID int

type FIFRepository interface {
	GetHoldingQuantities(holdingsIDs []HoldingID, upUntil time.Time) map[HoldingID]FDRHoldingQuantity
	GetTrades(holdingsIDs []HoldingID, startDate time.Time, endDate time.Time) map[HoldingID][]FDRTradeActivity
}

type FDRHoldingQuantity struct {
	Quantity float64
	Name     string
	Symbol   string
}

type FDRTradeActivity struct {
	Date         time.Time
	Action       string
	Quantity     float64
	Price        float64
	ExchangeRate float64
	AmountInNZD  float64
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

const fdrRate = 0.05

func (s FIFCalculationService) FDRIncome(input FDRInput, startDate time.Time, endDate time.Time) FDRResult {
	// 1. Get all the holding IDs that need to be computed
	holdingIDs := make([]HoldingID, 0, len(input.Holdings))
	for _, holding := range input.Holdings {
		holdingIDs = append(holdingIDs, holding.HoldingID)
	}

	// 2. Get the quantity for each holding at the start date

	holdings := s.repository.GetHoldingQuantities(holdingIDs, startDate)

	// 3. For each holding calculate the FDR (5% * opening market value)

	result := FDRResult{}
	for _, holding := range input.Holdings {
		info, ok := holdings[holding.HoldingID]

		if !ok {
			// Potentially return an error here
			continue
		}

		openingValue := info.Quantity * holding.OpeningPrice * holding.ExchangeRateToNZD

		r := FDRHoldingResult{
			Name:                info.Name,
			Symbol:              info.Symbol,
			OpeningValue:        openingValue,
			QuickSaleAdjustment: 0,
			Income:              openingValue * fdrRate,
		}
		result.Holdings = append(result.Holdings, r)
	}

	return result
}

type PeakDifferentialResult struct {
	PeakQuantity  float64
	QuantityStart float64
	QuantityEnd   float64
	AverageCost   float64
	Result        float64
}

func (s FIFCalculationService) PeakHoldingDifferential(
	holdingInfo FDRHoldingQuantity,
	trades []FDRTradeActivity) PeakDifferentialResult {

	peakQuantity := holdingInfo.Quantity
	currentQuantity := holdingInfo.Quantity
	var totalBuyQuantity float64
	var totalBuyAmount float64

	for _, trade := range trades {
		switch trade.Action {
		case constants.Buy:
			currentQuantity += trade.Quantity
			totalBuyQuantity += trade.Quantity
			totalBuyAmount += trade.Quantity * trade.Price * trade.ExchangeRate
		case constants.Sell:
			currentQuantity -= trade.Quantity
		}
		peakQuantity = max(peakQuantity, currentQuantity)
	}

	if totalBuyQuantity == 0 && totalBuyAmount == 0 {
		return PeakDifferentialResult{}
	}

	peakToStart := peakQuantity - holdingInfo.Quantity
	peakToEnd := peakQuantity - currentQuantity
	peakDifferential := min(peakToStart, peakToEnd)

	averageCost := totalBuyAmount / totalBuyQuantity

	return PeakDifferentialResult{
		PeakQuantity:  peakDifferential,
		QuantityStart: holdingInfo.Quantity,
		QuantityEnd:   currentQuantity,
		AverageCost:   averageCost,
		Result:        averageCost * fdrRate * peakDifferential,
	}
}
