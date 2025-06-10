package model

type FIFResult struct {
	Symbol     string
	StartValue float64
	EndValue   float64
	FDRAmount  float64
}

type FRDResult struct {
	Symbol              string
	StartValue          float64
	QuickSaleAdjustment float64
}

func (f FRDResult) TotalFRDIncome() float64 {
	return f.QuickSaleAdjustment + f.StartValue*0.05
}
