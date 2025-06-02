package fifservice

import (
	"fif-calculator/internal/constants"
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
