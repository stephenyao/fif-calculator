package fifservice

import (
	"fif-calculator/internal/repository"
	"time"
)

type FIFService interface {
	ComputeFRDIncome(startDate time.Time, endDate time.Time) FairDividendRateResult
}

type FIFCalculationService struct {
	repository repository.TradeRepository
}

type FairDividendRateResult struct {
	Holdings []FairDividendRateHolding
}

type FairDividendRateHolding struct {
	Name                string
	Symbol              string
	OpeningValue        float64
	QuickSaleAdjustment float64
	Income              float64
}

func NewFIFService(trade repository.TradeRepository) FIFService {
	return FIFCalculationService{
		repository: trade,
	}
}

func (s FIFCalculationService) ComputeFRDIncome(startDate time.Time, endDate time.Time) FairDividendRateResult {
	return FairDividendRateResult{}
}
