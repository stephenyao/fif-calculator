package repository

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func InitSchema(db *sqlx.DB) {
	schema := `
	CREATE TABLE IF NOT EXISTS trades (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		symbol TEXT,
		buy_date TEXT,
		quantity REAL,
		price REAL,
		currency TEXT,
		action TEXT
	);`
	_, err := db.Exec(schema)
	if err != nil {
		log.Fatalf("failed to create schema: %v", err)
	}
}
