package repository

import (
	"fif-calculator/internal/model"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	t.Helper()
	db := sqlx.MustConnect("sqlite3", ":memory:")
	db.MustExec("PRAGMA foreign_keys = ON;")
	InitSchema(db)
	return db
}

func TestGetCalculationWithHoldings(t *testing.T) {
	db := setupTestDB(t)
	db.MustExec("PRAGMA foreign_keys = ON;")

	repo := NewFIFRepository(db)

	// Insert a calculation
	calc := model.FIFCalculation{
		UserID:        99,
		FinancialYear: 2026,
		CalculatedAt:  time.Now().UTC().Truncate(time.Second),
	}
	res, err := db.Exec(`
		INSERT INTO fif_calculations (user_id, financial_year, calculated_at)
		VALUES (?, ?, ?)
	`, calc.UserID, calc.FinancialYear, calc.CalculatedAt)
	if err != nil {
		t.Fatalf("failed to insert calc: %v", err)
	}
	calcID, _ := res.LastInsertId()

	// Insert holdings
	holdings := []model.FIFHolding{
		{
			CalculationID:     int(calcID),
			Symbol:            "AAPL",
			QuantityStart:     1,
			QuantityEnd:       2,
			PriceStart:        3,
			PriceEnd:          4,
			ProceedsFromSales: 5,
			Dividends:         6,
			TaxCredits:        7,
			OtherGains:        8,
			CostOfPurchases:   9,
			ForeignIncomeTax:  10,
			OtherCosts:        11,
		},
	}
	for _, h := range holdings {
		_, err := db.Exec(`
			INSERT INTO fif_holdings (
				fif_calculation_id, symbol, quantity_start, quantity_end,
				price_start, price_end, proceeds_from_sales, dividends,
				tax_credits, other_gains, cost_of_purchases,
				foreign_income_tax, other_costs
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, h.CalculationID, h.Symbol, h.QuantityStart, h.QuantityEnd,
			h.PriceStart, h.PriceEnd, h.ProceedsFromSales, h.Dividends,
			h.TaxCredits, h.OtherGains, h.CostOfPurchases,
			h.ForeignIncomeTax, h.OtherCosts)
		if err != nil {
			t.Fatalf("failed to insert holding: %v", err)
		}
	}

	// Retrieve and verify
	gotCalc, gotHoldings, err := repo.GetCalculationWithHoldings(int(calcID))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotCalc.UserID != calc.UserID ||
		gotCalc.FinancialYear != calc.FinancialYear ||
		!gotCalc.CalculatedAt.Equal(calc.CalculatedAt) {
		t.Errorf("calculation mismatch: got %+v, want %+v", gotCalc, calc)
	}

	if len(gotHoldings) != len(holdings) {
		t.Fatalf("expected %d holdings, got %d", len(holdings), len(gotHoldings))
	}

	want := holdings[0]
	got := gotHoldings[0]
	if got.Symbol != want.Symbol ||
		got.QuantityStart != want.QuantityStart ||
		got.QuantityEnd != want.QuantityEnd ||
		got.PriceStart != want.PriceStart ||
		got.PriceEnd != want.PriceEnd ||
		got.ProceedsFromSales != want.ProceedsFromSales ||
		got.Dividends != want.Dividends ||
		got.TaxCredits != want.TaxCredits ||
		got.OtherGains != want.OtherGains ||
		got.CostOfPurchases != want.CostOfPurchases ||
		got.ForeignIncomeTax != want.ForeignIncomeTax ||
		got.OtherCosts != want.OtherCosts {
		t.Errorf("holding mismatch: got %+v, want %+v", got, want)
	}
}

func TestCreateCalculation(t *testing.T) {
	db := setupTestDB(t)
	repo := NewFIFRepository(db)

	calc := &model.FIFCalculation{
		UserID:        1,
		FinancialYear: 2025,
		CalculatedAt:  time.Now().UTC().Truncate(time.Second),
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
		t.Fatalf("CreateCalculation failed: %v", err)
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

		if actualCalc.UserID != calc.UserID || actualCalc.FinancialYear != calc.FinancialYear {
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

func TestGetCalculationsByUser(t *testing.T) {
	db := setupTestDB(t)
	repo := NewFIFRepository(db)

	now := time.Now().UTC().Truncate(time.Second)

	_, err := db.Exec(`INSERT INTO fif_calculations (user_id, financial_year, calculated_at) VALUES (?, ?, ?)`, 1, 2024, now)
	if err != nil {
		t.Fatalf("insert calc1: %v", err)
	}
	_, err = db.Exec(`INSERT INTO fif_calculations (user_id, financial_year, calculated_at) VALUES (?, ?, ?)`, 1, 2025, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("insert calc2: %v", err)
	}

	calcs, err := repo.GetCalculationsByUser(1)
	if err != nil {
		t.Fatalf("GetCalculationsByUser failed: %v", err)
	}
	if len(calcs) != 2 {
		t.Fatalf("expected 2 results, got %d", len(calcs))
	}
	if calcs[0].FinancialYear != 2025 || calcs[1].FinancialYear != 2024 {
		t.Errorf("expected ordered years 2025 then 2024, got %d then %d", calcs[0].FinancialYear, calcs[1].FinancialYear)
	}
}

func TestCreateOrUpdateCalculation(t *testing.T) {
	db := setupTestDB(t)

	repo := NewFIFRepository(db)

	userID := 1
	finYear := 2025
	calcTime := time.Now()

	initialHoldings := []*model.FIFHolding{
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
	}

	updatedHoldings := []*model.FIFHolding{
		{
			Symbol:            "AAPL",
			QuantityStart:     20,
			QuantityEnd:       22,
			PriceStart:        250,
			PriceEnd:          255,
			ProceedsFromSales: 2000,
			Dividends:         60,
			TaxCredits:        10,
			OtherGains:        40,
			CostOfPurchases:   1800,
			ForeignIncomeTax:  6,
			OtherCosts:        4,
		},
		{
			Symbol:            "GOOG",
			QuantityStart:     5,
			QuantityEnd:       6,
			PriceStart:        3000,
			PriceEnd:          3100,
			ProceedsFromSales: 3000,
			Dividends:         0,
			TaxCredits:        0,
			OtherGains:        0,
			CostOfPurchases:   2500,
			ForeignIncomeTax:  0,
			OtherCosts:        0,
		},
	}

	t.Run("initial insert", func(t *testing.T) {
		err := repo.CreateOrUpdateCalculation(&model.FIFCalculation{
			UserID:        userID,
			FinancialYear: finYear,
			CalculatedAt:  calcTime,
		}, initialHoldings)
		if err != nil {
			t.Fatalf("failed initial insert: %v", err)
		}

		var count int
		_ = db.Get(&count, "SELECT COUNT(*) FROM fif_holdings")
		if count != 1 {
			t.Errorf("expected 1 holding, got %d", count)
		}
	})

	t.Run("update overwrite", func(t *testing.T) {
		err := repo.CreateOrUpdateCalculation(&model.FIFCalculation{
			UserID:        userID,
			FinancialYear: finYear,
			CalculatedAt:  calcTime.Add(1 * time.Hour),
		}, updatedHoldings)
		if err != nil {
			t.Fatalf("failed update insert: %v", err)
		}

		var count int
		_ = db.Get(&count, "SELECT COUNT(*) FROM fif_holdings")
		if count != 2 {
			t.Errorf("expected 2 holdings after update, got %d", count)
		}

		var symbol string
		err = db.Get(&symbol, "SELECT symbol FROM fif_holdings WHERE symbol = 'GOOG'")
		if err != nil {
			t.Errorf("GOOG holding not found after update")
		}
	})
}
