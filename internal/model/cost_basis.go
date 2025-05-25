package model

import (
	"fif-clacultor/internal/constants"
)

func CostBasisBySymbol(trades []Trade) map[string]float64 {
	costBasisBySymbol := make(map[string]float64)

	// Build up a map of trades by symbol so they can be iterated through

	tradesBySymbol := make(map[string][]Trade)

	for _, trade := range trades {
		tradesBySymbol[trade.Symbol] = append(tradesBySymbol[trade.Symbol], trade)
	}

	// Iterate through each trade by symbol and calculate the cost basis
	for symbol, trades := range tradesBySymbol {
		var queue []Trade

		for _, trade := range trades {
			switch trade.Action {
			case constants.Buy:
				queue = append(queue, trade)
			case constants.Sell:
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
		costBasisBySymbol[symbol] = costBasisForSymbol
	}

	return costBasisBySymbol
}
