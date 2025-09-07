package repository

import (
	"database/sql"
	"fif-calculator/internal/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"testing"
	"time"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	db, err := sqlx.Open("sqlite3", "file::memory:?cache=shared&_foreign_keys=on")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	InitSchema(db)

	t.Cleanup(func() { _ = db.Close() })
	return db
}

const userID = "1"

func TestHoldingsRepository(t *testing.T) {
	t.Run("CreateHolding stores record and assigns ID", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewHoldingsRepository(db)
		record := &model.HoldingRecord{
			UserID:   userID,
			Name:     "Apple",
			Symbol:   "APPL",
			Currency: "USD",
		}

		err := repo.CreateHolding(record)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if record.ID == 0 {
			t.Fatal("expected ID to be set after insert")
		}

		actual, err := repo.GetHolding(record.ID, userID)
		if err != nil {
			t.Fatalf("failed to fetch inserted record: %v", err)
		}

		assertEqual(t, "UserID", actual.UserID, record.UserID)
		assertEqual(t, "Name", actual.Name, record.Name)
		assertEqual(t, "Symbol", actual.Symbol, record.Symbol)
		assertEqual(t, "Currency", actual.Currency, record.Currency)
		assertTimeClose(t, "CreatedAt", actual.CreatedAt, time.Now(), 2*time.Second)
	})

	t.Run("GetHolding returns error if not found", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewHoldingsRepository(db)
		_, err := repo.GetHolding(9999, userID)
		if err == nil {
			t.Fatal("expected error for missing holding, got nil")
		}
		if err != sql.ErrNoRows {
			t.Errorf("expected sql.ErrNoRows, got %v", err)
		}
	})

	t.Run("AllHoldings returns list of all holdings in desc order", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewHoldingsRepository(db)
		record1 := &model.HoldingRecord{
			UserID:   userID,
			Name:     "Microsoft",
			Symbol:   "MSFT",
			Currency: "USD",
		}
		record2 := &model.HoldingRecord{
			UserID:   userID,
			Name:     "Apple",
			Symbol:   "APPL",
			Currency: "USD",
		}

		err := repo.CreateHolding(record1)
		if err != nil {
			t.Fatalf("failed to insert record1: %v", err)
		}
		err = repo.CreateHolding(record2)
		if err != nil {
			t.Fatalf("failed to insert record2: %v", err)
		}

		records, err := repo.AllHoldings(userID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(records) != 2 {
			t.Fatalf("expected 2 records, got %d", len(records))
		}

		assertEqual(t, "Symbol", records[0].Symbol, "MSFT")
		assertEqual(t, "Symbol", records[1].Symbol, "APPL")
	})
}
func assertEqual(t *testing.T, field string, got, want any) {
	if got != want {
		t.Errorf("expected %s to be %v, got %v", field, want, got)
	}
}

func assertTimeClose(t *testing.T, field string, got, want time.Time, margin time.Duration) {
	if got.Sub(want) > margin || want.Sub(got) > margin {
		t.Errorf("expected %s within %v of %v, got %v", field, margin, want, got)
	}
}
