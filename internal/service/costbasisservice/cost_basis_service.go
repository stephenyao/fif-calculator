package costbasisservice

import (
	"fif-calculator/internal/constants"
	. "fif-calculator/internal/model"
	. "fif-calculator/internal/viewmodel"
	"sort"
	"time"
)

func CostBasisBySymbol(trades []Trade, untilDate time.Time) map[string]SymbolCostBasis {
	costBasisBySymbol := make(map[string]SymbolCostBasis)

	// Build up a map of trades by symbol so they can be iterated through

	tradesBySymbol := make(map[string][]Trade)

	for _, trade := range trades {
		tradesBySymbol[trade.Symbol] = append(tradesBySymbol[trade.Symbol], trade)
	}

	// Iterate through each trade by symbol and calculate the cost basis
	for symbol, trades := range tradesBySymbol {
		var queue []Trade
		var totalBought, totalSold float64

		for _, trade := range trades {
			tradeDate, _ := time.Parse(time.DateOnly, trade.BuyDate)

			if tradeDate.After(untilDate) {
				break
			}

			switch trade.Action {
			case constants.Buy:
				queue = append(queue, trade)
				totalBought += trade.Quantity
			case constants.Sell:
				totalSold += trade.Quantity
				if len(queue) == 0 {
					continue
				}

				sellQuantity := trade.Quantity

				// As long as there's still sell quantity, keep draining the queue
				for sellQuantity > 0 {
					buyQuantity := queue[0].Quantity

					if buyQuantity > sellQuantity {
						queue[0].Quantity -= sellQuantity
						sellQuantity = 0
					} else {
						queue = queue[1:]
						sellQuantity -= buyQuantity
					}

					if len(queue) == 0 {
						break
					}
				}
			}
		}

		var costBasisForSymbol float64
		for _, buyTrade := range queue {
			costBasisForSymbol += buyTrade.Price * buyTrade.Quantity
		}
		overSold := totalSold > totalBought
		costBasisBySymbol[symbol] =
			SymbolCostBasis{
				CostBasis:   costBasisForSymbol,
				TotalBought: totalBought,
				TotalSold:   totalSold,
				Oversold:    overSold,
				UntilDate:   untilDate,
			}
	}

	return costBasisBySymbol
}

func MaxCostBasisDuringYear(trades []Trade, startDate, endDate time.Time) float64 {
	tradesInRange := filterAndSortTrades(trades, startDate, endDate)

	var queue []Trade
	var currentCostBasis, maxCostBasis float64

	for _, trade := range tradesInRange {
		switch trade.Action {
		case constants.Buy:
			queue = append(queue, trade)
			currentCostBasis += trade.Price * trade.Quantity
		case constants.Sell:
			sellQty := trade.Quantity
			for sellQty > 0 && len(queue) > 0 {
				buy := &queue[0]
				if buy.Quantity > sellQty {
					currentCostBasis -= trade.Price * sellQty // conservative, or track exact match
					buy.Quantity -= sellQty
					sellQty = 0
				} else {
					currentCostBasis -= trade.Price * buy.Quantity
					sellQty -= buy.Quantity
					queue = queue[1:]
				}
			}
		}
		if currentCostBasis > maxCostBasis {
			maxCostBasis = currentCostBasis
		}
	}

	return maxCostBasis
}

func filterAndSortTrades(trades []Trade, startDate, endDate time.Time) []Trade {
	var result []Trade

	for _, trade := range trades {
		tradeDate, err := time.Parse("2006-01-02", trade.BuyDate)
		if err != nil {
			continue // skip invalid dates
		}
		if !tradeDate.Before(startDate) && !tradeDate.After(endDate) {
			result = append(result, trade)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		dateI, _ := time.Parse("2006-01-02", result[i].BuyDate)
		dateJ, _ := time.Parse("2006-01-02", result[j].BuyDate)
		return dateI.Before(dateJ)
	})

	return result
}
