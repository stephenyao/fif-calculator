package fifservice

import (
	"fif-calculator/internal/constants"
	"time"
)
import . "fif-calculator/internal/model"

func ComputeCVIncome(trades []Trade, holdings []*HoldingInfo, startDate, endDate time.Time) ([]CVResult, error) {
	var results []CVResult

	tradesBySymbol := tradesBySymbol(trades)

	for _, holding := range holdings {
		closingValue := holding.ClosingPrice * holding.QuantityEnd
		openingValue := holding.OpeningPrice * holding.QuantityStart
		var proceedsFromSales float64 = 0
		var costsOfPurchases float64 = 0

		for _, trade := range tradesBySymbol[holding.Symbol] {
			tradeDate, err := time.Parse(time.DateOnly, trade.BuyDate)

			if err != nil {
				return nil, err
			}

			if tradeDate.Before(startDate) || tradeDate.After(endDate) {
				continue
			}

			switch trade.Action {
			case constants.Sell:
				proceedsFromSales += trade.Price * trade.Quantity
			case constants.Buy:
				costsOfPurchases += trade.Price * trade.Quantity
			}
		}

		gains := proceedsFromSales + holding.GainLoss.Dividends + holding.GainLoss.OtherGains
		costs := costsOfPurchases + holding.GainLoss.OtherCosts

		result := CVResult{
			Symbol:       holding.Symbol,
			OpeningValue: openingValue,
			ClosingValue: closingValue,
			Gains:        gains,
			Costs:        costs,
		}

		results = append(results, result)
	}

	return results, nil
}
