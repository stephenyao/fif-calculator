package view_model

type CostBasisViewModel struct {
	CostBasisBySymbol map[string]float64
	TotalCostBasis    float64
	IsValidForFIF     bool
}
