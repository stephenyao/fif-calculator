package fifservice

import (
	"fif-calculator/internal/constants"
	"time"
)
import . "fif-calculator/internal/model"

func ComputeCVIncome(trades []Trade, holdings []*HoldingQuantity, parameters CVParameters, startDate, endDate time.Time) ([]CVResult, error) {
	var results []CVResult

	tradesBySymbol := tradesBySymbol(trades)

	for _, holding := range holdings {
		closingValue := holding.ClosingPrice * holding.QuantityEnd
		openingValue := holding.OpeningPrice * holding.QuantityStart
		var proceedsFromSales float64 = 0
		var costsOfPurchases float64 = 0

		for _, trade := range tradesBySymbol[holding.Symbol] {
			tradeDate, err := time.Parse("2006-01-02", trade.BuyDate)

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

		gains := proceedsFromSales + parameters.Dividends + parameters.TaxCredits + parameters.OtherGains
		costs := costsOfPurchases + parameters.ForeignIncomeTax + parameters.OtherCosts

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
