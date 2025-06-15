package viewmodel

import "fif-calculator/internal/model"

type FIFCalculationViewModel struct {
	Year       int
	FDRResults []model.FRDResult
	TotalFDR   float64
	CVResults  []model.CVResult
	TotalCV    float64
}
