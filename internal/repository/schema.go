package repository

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func InitSchema(db *sqlx.DB) {
	schema := `
	CREATE TABLE IF NOT EXISTS holdings (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  user_id TEXT NOT NULL,
	  name TEXT,
	  symbol TEXT NOT NULL,
	  currency TEXT NOT NULL,
	  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS trades (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		holding_id INTEGER NOT NULL,
		buy_date TEXT,
		quantity REAL,
		price REAL,
		exchange_rate REAL,
		action TEXT,
		FOREIGN KEY (holding_id) REFERENCES holdings(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS fif_calculations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		financial_year INTEGER NOT NULL,
		calculated_at DATETIME NOT NULL DEFAULT (DATETIME('now'))
	);

	-- FIF holdings table
	CREATE TABLE IF NOT EXISTS fif_holdings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		fif_calculation_id INTEGER NOT NULL,
		symbol TEXT NOT NULL,
		quantity_start REAL NOT NULL,
		quantity_end REAL NOT NULL,
		price_start REAL,
		price_end REAL,
		proceeds_from_sales REAL,
		dividends REAL,
		tax_credits REAL,
		other_gains REAL,
		cost_of_purchases REAL,
		foreign_income_tax REAL,
		other_costs REAL,
		FOREIGN KEY(fif_calculation_id) REFERENCES fif_calculations(id) ON DELETE CASCADE
	);

	-- Index for quick lookups
	CREATE INDEX IF NOT EXISTS idx_fif_holdings_calculation ON fif_holdings(fif_calculation_id);
	`
	_, err := db.Exec(schema)
	if err != nil {
		log.Fatalf("failed to create schema: %v", err)
	}
}
