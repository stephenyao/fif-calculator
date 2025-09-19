package repository

import (
	"testing"

	"github.com/jmoiron/sqlx"
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

func seedDB(t *testing.T, db *sqlx.DB) {
	t.Helper()
	insertHoldings(t, db)
	insertTradesSingleHolding(t, db)
}

func insertHoldings(t *testing.T, db *sqlx.DB) {
	queryInsertHoldings := `
		INSERT INTO holdings (user_id, name, symbol, currency) 
		VALUES 
		    (?, ?, ?, ?),
		    (?, ?, ?, ?)
	`
	_, err := db.Exec(queryInsertHoldings,
		"1", "Google", "GOOG", "USD",
		"2", "Apple", "APPL", "USD",
	)

	if err != nil {
		t.Fatalf("insert holdings error: %v", err)
	}
}

func insertTradesSingleHolding(t *testing.T, db *sqlx.DB) {
	query := `
		INSERT INTO trades (holding_id, buy_date, quantity, price, exchange_rate, action)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(query, "1", "2024-08-08", 100, 100, 1.6, "buy")
	_, err = db.Exec(query, "1", "2024-09-08", 50, 200, 1.6, "sell")
	_, err = db.Exec(query, "1", "2024-09-09", 50, 200, 1.6, "buy")
	_, err = db.Exec(query, "1", "2024-09-10", 50, 200, 1.6, "sell")
	_, err = db.Exec(query, "1", "2024-09-11", 50, 200, 1.6, "buy")
	_, err = db.Exec(query, "1", "2024-09-12", 50, 200, 1.6, "buy")

	if err != nil {
		t.Fatalf("insert trades error: %v", err)
	}
}

func insertTradesMultipleHoldings(t *testing.T, db *sqlx.DB) {
	query := `
		INSERT INTO trades (holding_id, buy_date, quantity, price, exchange_rate, action)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(query, "1", "2024-08-08", 100, 100, 1.6, "buy")
	_, err = db.Exec(query, "1", "2024-09-08", 30, 200, 1.6, "sell")
	_, err = db.Exec(query, "1", "2024-09-09", 50, 200, 1.6, "buy")
	_, err = db.Exec(query, "1", "2024-09-10", 20, 200, 1.6, "sell")
	_, err = db.Exec(query, "1", "2024-09-11", 50, 200, 1.6, "buy")
	_, err = db.Exec(query, "1", "2024-09-12", 50, 200, 1.6, "buy")

	_, err = db.Exec(query, "2", "2024-08-08", 100, 100, 1.6, "buy")
	_, err = db.Exec(query, "2", "2024-09-08", 50, 200, 1.6, "sell")

	if err != nil {
		t.Fatalf("insert trades error: %v", err)
	}
}

func insertTradesNegativeQuantity(t *testing.T, db *sqlx.DB) {
	query := `
		INSERT INTO trades (holding_id, buy_date, quantity, price, exchange_rate, action)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := db.Exec(query, "2", "2024-08-08", 100, 100, 1.6, "buy")
	_, err = db.Exec(query, "2", "2024-09-08", 10000, 200, 1.6, "sell")

	if err != nil {
		t.Fatalf("insert trades error: %v", err)
	}
}

func deleteTrades() {

}
