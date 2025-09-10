package fifservice

import (
	"fif-calculator/internal/constants"
	"fif-calculator/internal/datastructures"
	"time"
)

type FIFService interface {
	FDRIncome(input FDRInput, startDate time.Time, endDate time.Time) FDRResult
	PeakHoldingDifferential(
		holdingInfo FDRHoldingQuantity,
		trades []FDRTradeActivity) PeakDifferentialResult
	RealGain(trades []FDRTradeActivity) RealGainResult
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

	if totalBuyQuantity == 0 {
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

type RealGainResult struct {
	Sales  []GainOnSale
	Result float64
}

type GainOnSale struct {
	Quantity          float64
	Gain              float64
	CostOfAcquisition float64
}

func (s FIFCalculationService) RealGain(trades []FDRTradeActivity) RealGainResult {
	type BuyActivity struct {
		quantity     float64
		price        float64
		exchangeRate float64
	}

	var result RealGainResult = RealGainResult{}
	var stack datastructures.GenericStack[*BuyActivity]
	var totalGain float64

	for _, trade := range trades {
		switch trade.Action {
		case constants.Buy:
			activity := BuyActivity{
				quantity:     trade.Quantity,
				price:        trade.Price,
				exchangeRate: trade.ExchangeRate,
			}
			stack.Push(&activity)
		case constants.Sell:
			sellQuantity := trade.Quantity

			var costOfAcquisition float64
			for sellQuantity > 0 {
				buyActivity, ok := stack.Peek()
				if !ok {
					break
				}
				// if the buy order is bigger than the sell order, drain the quanitity. Otherwise, pop the activity.
				if buyActivity.quantity > sellQuantity {
					buyActivity.quantity -= sellQuantity
					costOfAcquisition += buyActivity.price * buyActivity.exchangeRate * sellQuantity
					sellQuantity = 0
				} else {
					stack.Pop()
					sellQuantity -= buyActivity.quantity
					costOfAcquisition += buyActivity.price * buyActivity.exchangeRate * buyActivity.quantity
				}
			}

			totalGain += trade.AmountInNZD - costOfAcquisition
			sale := GainOnSale{
				Quantity:          trade.Quantity,
				Gain:              totalGain,
				CostOfAcquisition: costOfAcquisition,
			}
			result.Sales = append(result.Sales, sale)
		}
	}

	result.Result = totalGain

	return result
}
