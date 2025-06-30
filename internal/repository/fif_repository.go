package repository

import (
	"fif-calculator/internal/model"
	"github.com/jmoiron/sqlx"
	"strings"
)

type FIFRepository interface {
	CreateCalculation(calc *model.FIFCalculation, holdings []*model.FIFHolding) error
	GetCalculationsByUser(userID int) ([]*model.FIFCalculation, error)
}

type SqlFIFRepository struct {
	DB *sqlx.DB
}

func NewFIFRepository(db *sqlx.DB) FIFRepository {
	return &SqlFIFRepository{
		DB: db,
	}
}

func (r *SqlFIFRepository) GetCalculationsByUser(userID int) ([]*model.FIFCalculation, error) {
	var calculations []*model.FIFCalculation

	err := r.DB.Select(&calculations, `
		SELECT id, user_id, financial_year, calculated_at
		FROM fif_calculations
		WHERE user_id = ?
		ORDER BY calculated_at DESC
	`, userID)

	if err != nil {
		return nil, err
	}

	return calculations, nil
}

func (r *SqlFIFRepository) CreateCalculation(calc *model.FIFCalculation, holdings []*model.FIFHolding) error {
	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	// Insert into fif_calculations
	result, err := tx.Exec(`
		INSERT INTO fif_calculations (user_id, financial_year, calculated_at)
		VALUES (?, ?, ?)
	`, calc.UserID, calc.FinancialYear, calc.CalculatedAt)
	if err != nil {
		err = tx.Rollback()
		return err
	}

	calculationID, err := result.LastInsertId()
	if err != nil {
		err = tx.Rollback()
		return err
	}

	// Build batched INSERT query
	baseQuery := `
		INSERT INTO fif_holdings (
			fif_calculation_id, symbol, quantity_start, quantity_end,
			price_start, price_end, proceeds_from_sales, dividends,
			tax_credits, other_gains, cost_of_purchases,
			foreign_income_tax, other_costs
		) VALUES `

	placeholders := []string{}
	args := []interface{}{}

	for _, h := range holdings {
		placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		args = append(args,
			calculationID, h.Symbol, h.QuantityStart, h.QuantityEnd,
			h.PriceStart, h.PriceEnd, h.ProceedsFromSales, h.Dividends,
			h.TaxCredits, h.OtherGains, h.CostOfPurchases,
			h.ForeignIncomeTax, h.OtherCosts,
		)
	}

	fullQuery := baseQuery + strings.Join(placeholders, ",")
	_, err = tx.Exec(fullQuery, args...)
	if err != nil {
		err = tx.Rollback()
		return err
	}

	return tx.Commit()
}
