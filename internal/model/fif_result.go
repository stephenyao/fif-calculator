package model

type FRDResult struct {
	Symbol              string
	StartValue          float64
	QuickSaleAdjustment float64
}

func (f FRDResult) TotalFDRIncome() float64 {
	return f.QuickSaleAdjustment + f.StartValue*0.05
}

type CVResult struct {
	Symbol       string
	OpeningValue float64
	ClosingValue float64
	Gains        float64
	Costs        float64
}

func (c CVResult) TotalIncome() float64 {
	return (c.ClosingValue + c.Gains) - (c.OpeningValue + c.Costs)
}
