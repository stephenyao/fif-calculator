package fifservice

import (
	"fif-calculator/internal/constants"
	. "fif-calculator/internal/model"
	"time"
)

func ComputeHoldingsBetween(trades []Trade, startDate, endDate time.Time) ([]*HoldingInfo, error) {
	quantityStart := make(map[string]float64)
	quantityEnd := make(map[string]float64)
	numberOfTrades := make(map[string]int)
	proceedsFromSales := make(map[string]float64)
	costOfPurchases := make(map[string]float64)

	for _, trade := range trades {
		tradeDate, err := time.Parse(time.DateOnly, trade.BuyDate)

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
		if tradeDate.After(startDate) && tradeDate.Before(endDate) {
			numberOfTrades[trade.Symbol]++
			switch trade.Action {
			case constants.Sell:
				proceedsFromSales[trade.Symbol] += trade.Price * trade.Quantity
			case constants.Buy:
				costOfPurchases[trade.Symbol] += trade.Price * trade.Quantity
			}
		}
	}

	symbolSet := map[string]struct{}{}
	for _, trade := range trades {
		symbolSet[trade.Symbol] = struct{}{}
	}

	var result []*HoldingInfo
	for symbol := range symbolSet {
		quantityStart := quantityStart[symbol]
		quantityEnd := quantityEnd[symbol]
		numberOfTrades := numberOfTrades[symbol]
		proceedsFromSales := proceedsFromSales[symbol]
		costOfPurchases := costOfPurchases[symbol]

		if quantityStart < 0 || quantityEnd < 0 {
			continue
		}

		if quantityStart == 0 && quantityEnd == 0 && numberOfTrades == 0 {
			continue
		}

		result = append(result, &HoldingInfo{
			Symbol:            symbol,
			QuantityStart:     quantityStart,
			QuantityEnd:       quantityEnd,
			NumberOfTrades:    numberOfTrades,
			ProceedsFromSales: proceedsFromSales,
			CostOfPurchases:   costOfPurchases,
		})
	}

	return result, nil
}

func tradesBySymbol(trades []Trade) map[string][]Trade {
	tradesBySymbol := map[string][]Trade{}

	for _, trade := range trades {
		tradesBySymbol[trade.Symbol] = append(tradesBySymbol[trade.Symbol], trade)
	}

	return tradesBySymbol
}
