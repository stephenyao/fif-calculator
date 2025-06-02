package fifservice

import "time"
import . "fif-calculator/internal/model"

func ComputeHoldingsBetween(start, end time.Time) []HoldingQuantity {
	// Fetch trades up to and including `start`, and separately up to and including `end`
	// Sum buy/sell quantities for each symbol to get net holding at each point

	return []HoldingQuantity{
		{Symbol: "AAPL", QuantityStart: 20, QuantityEnd: 30},
		{Symbol: "TSLA", QuantityStart: 5, QuantityEnd: 0},
	}
}
