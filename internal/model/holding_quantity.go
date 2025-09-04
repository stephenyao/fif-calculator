package model

import "strconv"

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
	Year              int
	HoldingId         int
}

type GainLossParams struct {
	Dividends  float64
	OtherGains float64
	OtherCosts float64
}

func (h HoldingInfo) GetHoldingFIFInfoURL() string {
	return "/fif/holding/" + strconv.Itoa(h.HoldingId) + "/year/" + strconv.Itoa(h.Year)

}
