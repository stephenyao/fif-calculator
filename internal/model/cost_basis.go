package model

import (
	"fif-clacultor/internal/constants"
	"math"
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
				// If buy queue is empty then move onto the next trade
				if len(queue) == 0 {
					continue
				}

				// If the first buy order quantity can absorb the sell order, reduce that buy order quantity and skip to next order
				if queue[0].Quantity > trade.Quantity {
					queue[0].Quantity -= trade.Quantity
					continue
				}

				remainder := trade.Quantity - queue[0].Quantity
				// While there is more of the sell order left, keep draining the queue
				for remainder >= 0 {
					queue = queue[1:]

					if len(queue) == 0 {
						break
					}

					if remainder == 0 {
						break
					}

					remainder -= queue[0].Quantity
				}

				// Mutate and subtract the quantity from the first Buy trade in the queue
				if remainder < 0 {
					queue[0].Quantity -= math.Abs(remainder)
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
