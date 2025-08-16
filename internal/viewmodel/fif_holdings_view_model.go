package viewmodel

import (
	"fif-calculator/internal/model"
	"strconv"
)

type FIFHoldingsViewModel struct {
	Holdings []*model.HoldingInfo
	Title    string
	Year     int
}

func CreateFIFHoldingsViewModel(holdings []*model.HoldingInfo, financialYear int) FIFHoldingsViewModel {
	title := "Holdings for Financial Year " + strconv.Itoa(financialYear-1) + "-" + strconv.Itoa(financialYear)

	return FIFHoldingsViewModel{
		Holdings: holdings,
		Title:    title,
		Year:     financialYear,
	}
}
