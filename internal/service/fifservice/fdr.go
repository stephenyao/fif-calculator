package fifservice

import (
	"fif-calculator/internal/constants"
	"fif-calculator/internal/datastructures"
	"log"
	"time"
)
import . "fif-calculator/internal/model"

func ComputeFRDIncome(trades []Trade, holdings []*HoldingInfo, startDate, endDate time.Time) ([]FRDResult, error) {
	var result []FRDResult

	tradesBySymbol := tradesBySymbol(trades)

	for _, holding := range holdings {

		peakDifferential, err := peakDifferentialForSymbol(*holding, tradesBySymbol[holding.Symbol], startDate, endDate)
		log.Printf("peak differential for symbol %s: %v", holding.Symbol, peakDifferential)
		actualGain, err := calculateRealGainForSymbol(tradesBySymbol[holding.Symbol], startDate, endDate)

		if err != nil {
			return nil, err
		}

		quickSaleAdjustment := max(0, min(peakDifferential, actualGain))
		r := FRDResult{
			Symbol:              holding.Symbol,
			StartValue:          holding.QuantityStart * holding.OpeningPrice,
			QuickSaleAdjustment: quickSaleAdjustment,
		}

		result = append(result, r)
	}

	return result, nil
}

func peakDifferentialForSymbol(holding HoldingInfo, trades []Trade, startDate, endDate time.Time) (float64, error) {
	var trackedQuantity float64
	var totalCost float64
	var totalBuyQuantity float64
	var peakQuantity float64

	isValid := false
	for _, trade := range trades {
		// find the trade with the closest start date
		tradeDate, err := time.Parse(time.DateOnly, trade.BuyDate)

		if err != nil {
			return 0, err
		}

		if tradeDate.Before(startDate) || tradeDate.After(endDate) {
			continue
		}

		isValid = true

		switch trade.Action {
		case constants.Buy:
			trackedQuantity += trade.Quantity
			totalCost += trade.Quantity * trade.Price * trade.ExchangeRate
			totalBuyQuantity += trade.Quantity
		case constants.Sell:
			trackedQuantity -= trade.Quantity
			trackedQuantity = max(trackedQuantity, 0)
		}

		peakQuantity = max(peakQuantity, trackedQuantity)
	}

	if !isValid {
		return 0, nil
	}

	averageCost := totalCost / totalBuyQuantity
	peakQuantity = holding.QuantityStart + peakQuantity

	peakToStart := peakQuantity - holding.QuantityStart
	peakToEnd := peakQuantity - holding.QuantityEnd

	peakDifferential := min(peakToStart, peakToEnd)

	return peakDifferential * averageCost * 0.05, nil
}

func calculateRealGainForSymbol(trades []Trade, startDate, endDate time.Time) (float64, error) {
	var realGain float64
	queue := &datastructures.Stack{}
	for _, trade := range trades {

		// find the trade with the closest start date
		tradeDate, err := time.Parse(time.DateOnly, trade.BuyDate)

		if err != nil {
			return 0, err
		}

		if tradeDate.Before(startDate) || tradeDate.After(endDate) {
			continue
		}

		switch trade.Action {
		case constants.Buy:
			queue.Push(trade)
		case constants.Sell:
			quantityRemaining := trade.Quantity

			for quantityRemaining > 0 {
				lastTrade, success := queue.Pop()

				if !success {
					break
				}

				quantityToDrain := lastTrade.Quantity

				isDrained := (quantityToDrain - quantityRemaining) <= 0

				if !isDrained {
					lastTrade.Quantity -= quantityRemaining
					queue.Push(lastTrade)
				}

				quantityUsed := min(quantityRemaining, quantityToDrain)

				realGain += quantityUsed*trade.Price*trade.ExchangeRate - quantityUsed*lastTrade.Price*lastTrade.ExchangeRate
				quantityRemaining -= quantityToDrain
			}
		}
	}
	return realGain, nil
}
