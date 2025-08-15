package viewmodel

import (
	"fmt"
	"time"
)

type FIFStartViewModel struct {
	FinancialYearOptions SelectOptions
}

func CreateFIFStartViewModel() FIFStartViewModel {
	endYear := time.Now().Year() + 1
	startYear := endYear - 9
	options := []Option{}

	for year := startYear; year <= endYear; year++ {
		options = append(options, Option{
			Value:   fmt.Sprintf("%d", year),
			Display: fmt.Sprintf("1 Apr %d – 31 Mar %d", year-1, year),
		})
	}

	return FIFStartViewModel{
		FinancialYearOptions: SelectOptions{Options: options},
	}
}
