package costbasisservice

import (
	"fif-clacultor/internal/constants"
	. "fif-clacultor/internal/model"
	. "fif-clacultor/internal/viewmodel"
)

func CostBasisBySymbol(trades []Trade) map[string]SymbolCostBasis {
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
		costBasisBySymbol[symbol] = SymbolCostBasis{costBasisForSymbol, totalBought, totalSold, overSold}
	}

	return costBasisBySymbol
}
