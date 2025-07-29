package repository

import (
	"fif-calculator/internal/model"
	"github.com/jmoiron/sqlx"
	"strings"
)

type HoldingsRepository interface {
	CreateHolding(record *model.HoldingRecord) error
	GetHolding(id int, userID string) (*model.HoldingRecord, error)
	AllHoldings(userID string) ([]*model.HoldingRecord, error)
	DeleteByID(id int, userID string) error
	Update(record *model.HoldingRecord) error
}

type SQLHoldingsRepository struct {
	DB *sqlx.DB
}

func NewHoldingsRepository(db *sqlx.DB) *SQLHoldingsRepository {
	return &SQLHoldingsRepository{
		DB: db,
	}
}

func (r *SQLHoldingsRepository) CreateHolding(record *model.HoldingRecord) error {
	// Uppercase the symbol as the ticker is not typically lower case
	record.Symbol = strings.ToUpper(record.Symbol)

	result, err := r.DB.Exec(`
		INSERT INTO holdings (user_id, name, symbol, currency)
		VALUES (?, ?, ?, ?)`,
		record.UserID, record.Name, record.Symbol, record.Currency,
	)

	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	record.ID = int(id)
	return nil
}

func (r *SQLHoldingsRepository) GetHolding(id int, userID string) (*model.HoldingRecord, error) {
	var record model.HoldingRecord
	err := r.DB.Get(&record, `SELECT * FROM holdings WHERE id=? AND user_id=?`, id, userID)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *SQLHoldingsRepository) AllHoldings(userID string) ([]*model.HoldingRecord, error) {
	var records []*model.HoldingRecord
	err := r.DB.Select(&records, `SELECT * FROM holdings WHERE user_id=? ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r *SQLHoldingsRepository) DeleteByID(id int, userID string) error {
	_, err := r.DB.Exec("DELETE FROM holdings WHERE id = ? AND user_id = ?", id, userID)
	return err
}

func (r *SQLHoldingsRepository) Update(record *model.HoldingRecord) error {
	// Uppercase the symbol as the ticker is not typically lower case
	_, err := r.DB.Exec(`
		UPDATE holdings
		SET name = ?, symbol = ?, currency = ?
		WHERE id = ? AND user_id = ?`,
		record.Name, record.Symbol, record.Currency, record.ID, record.UserID,
	)
	return err
}
