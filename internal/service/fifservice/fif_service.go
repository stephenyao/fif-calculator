package fifservice

import (
	"fif-calculator/internal/constants"
	"fif-calculator/internal/datastructures"
	"time"
)
import . "fif-calculator/internal/model"

func ComputeHoldingsBetween(trades []Trade, startDate, endDate time.Time) ([]HoldingQuantity, error) {
	quantityStart := make(map[string]float64)
	quantityEnd := make(map[string]float64)

	for _, trade := range trades {
		tradeDate, err := time.Parse("2006-01-02", trade.BuyDate)

		if err != nil {
			return nil, err
		}

		var delta float64
		switch trade.Action {
		case constants.Buy:
			delta = trade.Quantity
		case constants.Sell:
			delta = -trade.Quantity
		}

		if !tradeDate.After(startDate) {
			quantityStart[trade.Symbol] += delta
		}
		if !tradeDate.After(endDate) {
			quantityEnd[trade.Symbol] += delta
		}
	}

	symbolSet := map[string]struct{}{}
	for _, trade := range trades {
		symbolSet[trade.Symbol] = struct{}{}
	}

	var result []HoldingQuantity
	for symbol := range symbolSet {
		quantityStart := quantityStart[symbol]
		quantityEnd := quantityEnd[symbol]

		if quantityStart < 0 || quantityEnd < 0 {
			continue
		}

		result = append(result, HoldingQuantity{
			Symbol:        symbol,
			QuantityStart: quantityStart,
			QuantityEnd:   quantityEnd,
		})
	}

	return result, nil
}

func ComputeFRDIncome(trades []Trade, holdings []HoldingQuantity, startDate, endDate time.Time) (float64, error) {
	var income float64

	tradesBySymbol := tradesBySymbol(trades)

	for _, holding := range holdings {
		income += holding.QuantityStart * holding.OpeningPrice * 0.05

		peakDifferential, err := peakDifferentialForSymbol(holding, tradesBySymbol[holding.Symbol], startDate, endDate)
		actualGain, err := calculateRealGainForSymbol(tradesBySymbol[holding.Symbol], startDate, endDate)

		if err != nil {
			return 0, err
		}

		quickSaleAdjustment := max(0, min(peakDifferential, actualGain))
		income += quickSaleAdjustment
	}

	return income, nil
}

func peakDifferentialForSymbol(holding HoldingQuantity, trades []Trade, startDate, endDate time.Time) (float64, error) {
	var trackedQuantity float64
	var totalCost float64
	var totalBuyQuantity float64
	var peakQuantity float64

	isValid := false
	for _, trade := range trades {
		// find the trade with the closest start date
		tradeDate, err := time.Parse("2006-01-02", trade.BuyDate)

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
			totalCost += trade.Quantity * trade.Price
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
		tradeDate, err := time.Parse("2006-01-02", trade.BuyDate)

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
				realGain += quantityUsed*trade.Price - quantityUsed*lastTrade.Price
				quantityRemaining -= quantityToDrain
			}
		}
	}
	return realGain, nil
}

func tradesBySymbol(trades []Trade) map[string][]Trade {
	tradesBySymbol := map[string][]Trade{}

	for _, trade := range trades {
		tradesBySymbol[trade.Symbol] = append(tradesBySymbol[trade.Symbol], trade)
	}

	return tradesBySymbol
}
