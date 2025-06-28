package model

import "time"

type FIFCalculation struct {
	ID            int       `db:"id"`
	UserID        int       `db:"user_id"`
	FinancialYear int       `db:"financial_year"`
	CalculatedAt  time.Time `db:"calculated_at"`
}

type FIFHolding struct {
	ID                int     `db:"id"`
	CalculationID     int     `db:"fif_calculation_id"`
	Symbol            string  `db:"symbol"`
	QuantityStart     float64 `db:"quantity_start"`
	QuantityEnd       float64 `db:"quantity_end"`
	PriceStart        float64 `db:"price_start"`
	PriceEnd          float64 `db:"price_end"`
	ProceedsFromSales float64 `db:"proceeds_from_sales"`
	Dividends         float64 `db:"dividends"`
	TaxCredits        float64 `db:"tax_credits"`
	OtherGains        float64 `db:"other_gains"`
	CostOfPurchases   float64 `db:"cost_of_purchases"`
	ForeignIncomeTax  float64 `db:"foreign_income_tax"`
	OtherCosts        float64 `db:"other_costs"`
}
