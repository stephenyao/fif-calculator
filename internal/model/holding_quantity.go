package model

type HoldingInfo struct {
	Symbol            string
	QuantityStart     float64
	QuantityEnd       float64
	OpeningPrice      float64
	ClosingPrice      float64
	NumberOfTrades    int
	ProceedsFromSales float64
	CostOfPurchases   float64
	GainLoss          GainLossParams
}

type GainLossParams struct {
	Dividends        float64
	TaxCredits       float64
	OtherGains       float64
	ForeignIncomeTax float64
	OtherCosts       float64
}
