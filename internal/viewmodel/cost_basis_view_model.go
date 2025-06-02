package viewmodel

type SymbolCostBasis struct {
	CostBasis   float64
	TotalBought float64
	TotalSold   float64
	Oversold    bool
}

type CostBasisViewModel struct {
	CostBasisBySymbol map[string]SymbolCostBasis
	TotalCostBasis    float64
	IsValidForFIF     bool
}
