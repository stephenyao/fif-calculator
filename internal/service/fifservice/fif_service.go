package fifservice

import (
	"fif-calculator/internal/constants"
	. "fif-calculator/internal/repository"
	"time"
)

type FIFService interface {
	FDRIncome(input FDRInput, startDate time.Time, endDate time.Time) FDRResult
	PeakHoldingDifferential(
		holdingInfo FIFHoldingQuantity,
		trades []FIFTradeActivity) PeakDifferentialResult
	RealGain(trades []FIFTradeActivity) RealGainResult
	QuickSaleAdjustment(peakDifferential PeakDifferentialResult, realGain RealGainResult) QuickSaleAdjustmentResult
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
	QuickSaleAdjustment QuickSaleAdjustmentResult
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

const fdrRate = 0.05

func NewFIFService(repo FIFRepository) FIFService {
	return FIFCalculationService{
		repository: repo,
	}
}

func (s FIFCalculationService) FDRIncome(input FDRInput, startDate time.Time, endDate time.Time) FDRResult {
	// 1. Get all the holding IDs that need to be computed
	holdingIDs := make([]HoldingID, 0, len(input.Holdings))
	for _, holding := range input.Holdings {
		holdingIDs = append(holdingIDs, holding.HoldingID)
	}

	// 2. Get the quantity for each holding at the start date

	holdings := s.repository.GetHoldingQuantities(holdingIDs, startDate)
	activities := s.repository.GetTrades(holdingIDs, startDate, endDate)

	// 3. For each holding calculate the FDR (5% * opening market value)

	result := FDRResult{}
	for _, holding := range input.Holdings {
		info, ok := holdings[holding.HoldingID]

		if !ok {
			// Potentially return an error here
			continue
		}

		openingValue := info.Quantity * holding.OpeningPrice * holding.ExchangeRateToNZD
		peak := s.PeakHoldingDifferential(holdings[holding.HoldingID], activities[holding.HoldingID])
		realGain := s.RealGain(activities[holding.HoldingID])
		quickSales := s.QuickSaleAdjustment(peak, realGain)

		r := FDRHoldingResult{
			Name:                info.Name,
			Symbol:              info.Symbol,
			OpeningValue:        openingValue,
			QuickSaleAdjustment: quickSales,
			Income:              openingValue*fdrRate + quickSales.Result,
		}
		result.Holdings = append(result.Holdings, r)
	}

	return result
}

type QuickSaleAdjustmentResult struct {
	PeakDifferentialResult PeakDifferentialResult
	RealGainResult         RealGainResult
	Result                 float64
}

func (s FIFCalculationService) QuickSaleAdjustment(peakDifferential PeakDifferentialResult, realGain RealGainResult) QuickSaleAdjustmentResult {
	result := min(peakDifferential.Result, realGain.Result)

	return QuickSaleAdjustmentResult{
		PeakDifferentialResult: peakDifferential,
		RealGainResult:         realGain,
		Result:                 result,
	}
}

type PeakDifferentialResult struct {
	PeakQuantity     float64
	PeakDifferential float64
	QuantityStart    float64
	QuantityEnd      float64
	AverageCost      float64
	Result           float64
	TimeStamp        time.Time
}

func (s FIFCalculationService) PeakHoldingDifferential(
	holdingInfo FIFHoldingQuantity,
	trades []FIFTradeActivity) PeakDifferentialResult {

	peakQuantity := holdingInfo.Quantity
	currentQuantity := holdingInfo.Quantity
	var totalBuyQuantity float64
	var totalBuyAmount float64
	var totalSellQuantity float64
	var peakDate time.Time

	for _, trade := range trades {
		switch trade.Action {
		case constants.Buy:
			currentQuantity += trade.Quantity
			totalBuyQuantity += trade.Quantity
			totalBuyAmount += trade.Quantity * trade.Price * trade.ExchangeRate
		case constants.Sell:
			currentQuantity -= trade.Quantity
			totalSellQuantity += trade.Quantity
		}
		if currentQuantity > peakQuantity {
			peakQuantity = currentQuantity
			peakDate = trade.Date
		}
	}

	if totalBuyQuantity == 0 || totalSellQuantity == 0 {
		return PeakDifferentialResult{}
	}

	peakToStart := peakQuantity - holdingInfo.Quantity
	peakToEnd := peakQuantity - currentQuantity
	peakDifferential := min(peakToStart, peakToEnd)

	averageCost := totalBuyAmount / totalBuyQuantity

	return PeakDifferentialResult{
		PeakQuantity:     peakQuantity,
		PeakDifferential: peakDifferential,
		QuantityStart:    holdingInfo.Quantity,
		QuantityEnd:      currentQuantity,
		AverageCost:      averageCost,
		Result:           averageCost * fdrRate * peakDifferential,
		TimeStamp:        peakDate,
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

func (s FIFCalculationService) RealGain(trades []FIFTradeActivity) RealGainResult {
	type BuyActivity struct {
		quantity     float64
		price        float64
		exchangeRate float64
	}

	var result RealGainResult = RealGainResult{}
	var queue []*BuyActivity

	var totalGain float64
	for _, trade := range trades {
		switch trade.Action {
		case constants.Buy:
			activity := BuyActivity{
				quantity:     trade.Quantity,
				price:        trade.Price,
				exchangeRate: trade.ExchangeRate,
			}
			queue = append(queue, &activity)
		case constants.Sell:
			if len(queue) == 0 {
				break
			}

			sellQuantity := trade.Quantity
			var costOfAcquisition float64
			for sellQuantity > 0 {
				if len(queue) == 0 {
					break
				}

				buyActivity := queue[0]

				// if the buy order is bigger than the sell order, drain the quanitity. Otherwise, pop the activity.
				if buyActivity.quantity > sellQuantity {
					buyActivity.quantity -= sellQuantity
					costOfAcquisition += buyActivity.price * buyActivity.exchangeRate * sellQuantity
					sellQuantity = 0
				} else {
					queue = queue[1:]
					sellQuantity -= buyActivity.quantity
					costOfAcquisition += buyActivity.price * buyActivity.exchangeRate * buyActivity.quantity
				}
			}

			gainOnSale := (trade.Quantity-sellQuantity)*trade.Price*trade.ExchangeRate - costOfAcquisition

			if gainOnSale <= 0 {
				continue
			}

			matchedQuantity := trade.Quantity - sellQuantity
			totalGain += gainOnSale
			sale := GainOnSale{
				Quantity:          matchedQuantity,
				Gain:              gainOnSale,
				CostOfAcquisition: costOfAcquisition,
			}
			result.Sales = append(result.Sales, sale)
		}
	}

	result.Result = totalGain

	return result
}
