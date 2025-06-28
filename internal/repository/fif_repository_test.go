package repository_test

import (
	"fif-calculator/internal/model"
	"fif-calculator/internal/repository"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func TestCreateCalculation(t *testing.T) {
	db := sqlx.MustConnect("sqlite3", ":memory:")
	db.MustExec("PRAGMA foreign_keys = ON;")
	repository.InitSchema(db)

	repo := repository.NewFIFRepository(db)

	calc := &model.FIFCalculation{
		UserID:        1,
		FinancialYear: 2025,
		CalculatedAt:  time.Now(),
	}

	holdings := []*model.FIFHolding{
		{
			Symbol:            "AAPL",
			QuantityStart:     10,
			QuantityEnd:       12,
			PriceStart:        150,
			PriceEnd:          155,
			ProceedsFromSales: 1000,
			Dividends:         30,
			TaxCredits:        5,
			OtherGains:        20,
			CostOfPurchases:   900,
			ForeignIncomeTax:  3,
			OtherCosts:        2,
		},
		{
			Symbol:            "MSFT",
			QuantityStart:     5,
			QuantityEnd:       6,
			PriceStart:        300,
			PriceEnd:          310,
			ProceedsFromSales: 800,
			Dividends:         25,
			TaxCredits:        4,
			OtherGains:        10,
			CostOfPurchases:   700,
			ForeignIncomeTax:  2,
			OtherCosts:        1,
		},
	}

	err := repo.CreateCalculation(calc, holdings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Run("stores calculation and holdings", func(t *testing.T) {
		var count int
		err = db.Get(&count, "SELECT COUNT(*) FROM fif_calculations")
		if err != nil || count != 1 {
			t.Fatalf("expected 1 calculation, got %d (err=%v)", count, err)
		}

		err = db.Get(&count, "SELECT COUNT(*) FROM fif_holdings")
		if err != nil || count != 2 {
			t.Fatalf("expected 2 holdings, got %d (err=%v)", count, err)
		}
	})

	t.Run("links holdings to calculation", func(t *testing.T) {
		var calcID int
		err = db.Get(&calcID, "SELECT id FROM fif_calculations LIMIT 1")
		if err != nil {
			t.Fatalf("failed to fetch calculation id: %v", err)
		}

		var holdingCalcIDs []int
		err = db.Select(&holdingCalcIDs, "SELECT fif_calculation_id FROM fif_holdings")
		if err != nil {
			t.Fatalf("failed to fetch holdings: %v", err)
		}
		for _, id := range holdingCalcIDs {
			if id != calcID {
				t.Errorf("expected holding to be linked to calculation ID %d, got %d", calcID, id)
			}
		}
	})

	t.Run("round-trip values match inserted values", func(t *testing.T) {
		var actualCalc model.FIFCalculation
		err := db.Get(&actualCalc, "SELECT * FROM fif_calculations LIMIT 1")
		if err != nil {
			t.Fatalf("failed to fetch calculation: %v", err)
		}

		if actualCalc.UserID != calc.UserID ||
			actualCalc.FinancialYear != calc.FinancialYear {
			t.Errorf("calculation mismatch: got %+v, expected %+v", actualCalc, calc)
		}

		if !actualCalc.CalculatedAt.Equal(calc.CalculatedAt) {
			t.Errorf("calculatedAt mismatch: got %v, expected %v", actualCalc.CalculatedAt, calc.CalculatedAt)
		}

		var actualHoldings []model.FIFHolding
		err = db.Select(&actualHoldings, "SELECT * FROM fif_holdings ORDER BY symbol")
		if err != nil {
			t.Fatalf("failed to fetch holdings: %v", err)
		}

		if len(actualHoldings) != len(holdings) {
			t.Fatalf("expected %d holdings, got %d", len(holdings), len(actualHoldings))
		}

		for i, h := range holdings {
			a := actualHoldings[i]
			if a.Symbol != h.Symbol ||
				a.QuantityStart != h.QuantityStart ||
				a.QuantityEnd != h.QuantityEnd ||
				a.PriceStart != h.PriceStart ||
				a.PriceEnd != h.PriceEnd ||
				a.ProceedsFromSales != h.ProceedsFromSales ||
				a.Dividends != h.Dividends ||
				a.TaxCredits != h.TaxCredits ||
				a.OtherGains != h.OtherGains ||
				a.CostOfPurchases != h.CostOfPurchases ||
				a.ForeignIncomeTax != h.ForeignIncomeTax ||
				a.OtherCosts != h.OtherCosts {
				t.Errorf("holding mismatch:\ngot:  %+v\nwant: %+v", a, h)
			}
		}
	})
}
