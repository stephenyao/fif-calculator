package viewmodel

import "fif-calculator/internal/model"

type FIFCalculationViewModel struct {
	Year     int
	Results  []model.FRDResult
	TotalFDR float64
}
