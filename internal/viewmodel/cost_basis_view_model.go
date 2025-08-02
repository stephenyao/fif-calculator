package viewmodel

import "time"

type SymbolCostBasis struct {
	CostBasisFX     float64
	CostBasisNZD    float64
	TotalBought     float64
	TotalSold       float64
	Oversold        bool
	CurrentQuantity float64
	TotalTrades     int
	UntilDate       time.Time
}

type CostBasisViewModel struct {
	CostBasisBySymbol map[string]SymbolCostBasis
	TotalCostBasis    float64
	IsValidForFIF     bool
}
