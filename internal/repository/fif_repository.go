package repository

import (
	"database/sql"
	"fif-calculator/internal/model"
	"github.com/jmoiron/sqlx"
	"strings"
)

type FIFRepository interface {
	CreateCalculation(calc *model.FIFCalculation, holdings []*model.FIFHolding) error
	GetCalculationsByUser(userID int) ([]*model.FIFCalculation, error)
	GetCalculationWithHoldings(id int) (model.FIFCalculation, []*model.FIFHolding, error)
	CreateOrUpdateCalculation(calc *model.FIFCalculation, holdings []*model.FIFHolding) error
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

func (r *SqlFIFRepository) GetCalculationWithHoldings(id int) (model.FIFCalculation, []*model.FIFHolding, error) {
	var calc model.FIFCalculation
	err := r.DB.Get(&calc, `
		SELECT id, user_id, financial_year, calculated_at
		FROM fif_calculations
		WHERE id = ?
	`, id)
	if err != nil {
		return model.FIFCalculation{}, nil, err
	}

	var holdings []*model.FIFHolding
	err = r.DB.Select(&holdings, `
		SELECT id, fif_calculation_id, symbol, quantity_start, quantity_end,
		       price_start, price_end, proceeds_from_sales, dividends,
		       tax_credits, other_gains, cost_of_purchases,
		       foreign_income_tax, other_costs
		FROM fif_holdings
		WHERE fif_calculation_id = ?
	`, id)
	if err != nil {
		return model.FIFCalculation{}, nil, err
	}

	return calc, holdings, nil
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

func (r *SqlFIFRepository) CreateOrUpdateCalculation(calc *model.FIFCalculation, holdings []*model.FIFHolding) error {
	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	var calculationID int
	err = tx.Get(&calculationID, `
		SELECT id FROM fif_calculations
		WHERE user_id = ? AND financial_year = ?
	`, calc.UserID, calc.FinancialYear)

	if err != nil {
		if err == sql.ErrNoRows {
			// No existing calculation: insert new
			result, err := tx.Exec(`
				INSERT INTO fif_calculations (user_id, financial_year, calculated_at)
				VALUES (?, ?, ?)
			`, calc.UserID, calc.FinancialYear, calc.CalculatedAt)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			calculationID64, err := result.LastInsertId()
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			calculationID = int(calculationID64)
		} else {
			_ = tx.Rollback()
			return err // Unexpected error
		}
	} else {
		// Existing calculation found: delete old holdings
		_, err := tx.Exec(`
			DELETE FROM fif_holdings
			WHERE fif_calculation_id = ?
		`, calculationID)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	// Re-insert holdings
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
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
