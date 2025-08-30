package fifservice

import (
	"fif-calculator/internal/constants"
	. "fif-calculator/internal/model"
	"sort"
	"time"
)

func StartEndDates(year int) (time.Time, time.Time) {
	startDate := time.Date(year-1, 4, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, 3, 31, 0, 0, 0, 0, time.UTC)

	return startDate, endDate
}

func ComputeHoldingsBetween(trades []Trade, year int) ([]*HoldingInfo, error) {
	quantityStart := make(map[string]float64)
	quantityEnd := make(map[string]float64)
	numberOfTrades := make(map[string]int)
	proceedsFromSales := make(map[string]float64)
	costOfPurchases := make(map[string]float64)
	holdingIds := make(map[string]int)
	startDate, endDate := StartEndDates(year)

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

		holdingIds[trade.Symbol] = trade.HoldingID
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
		holdingID := holdingIds[symbol]

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
			Year:              year,
			HoldingId:         holdingID,
		})
	}

	// Sort the holdings alphbetically before returning
	sort.Slice(result, func(i, j int) bool {
		return result[i].Symbol < result[j].Symbol
	})

	return result, nil
}

func tradesBySymbol(trades []Trade) map[string][]Trade {
	tradesBySymbol := map[string][]Trade{}

	for _, trade := range trades {
		tradesBySymbol[trade.Symbol] = append(tradesBySymbol[trade.Symbol], trade)
	}

	return tradesBySymbol
}
