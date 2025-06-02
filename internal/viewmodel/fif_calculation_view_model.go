package viewmodel

import "fif-calculator/internal/model"

type FIFCalculationViewModel struct {
	Year     int
	Results  []model.FIFResult
	TotalFDR float64
}
