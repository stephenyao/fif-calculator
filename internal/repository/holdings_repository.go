package repository

import (
	"fif-calculator/internal/model"
	"github.com/jmoiron/sqlx"
)

type HoldingsRepository interface {
	CreateHolding(record *model.HoldingRecord) error
	GetHolding(id int) (*model.HoldingRecord, error)
	AllHoldings() ([]*model.HoldingRecord, error)
	DeleteByID(id int) error
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

func (r *SQLHoldingsRepository) GetHolding(id int) (*model.HoldingRecord, error) {
	var record model.HoldingRecord
	err := r.DB.Get(&record, `SELECT * FROM holdings WHERE id=?`, id)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *SQLHoldingsRepository) AllHoldings() ([]*model.HoldingRecord, error) {
	var records []*model.HoldingRecord
	err := r.DB.Select(&records, `SELECT * FROM holdings ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r *SQLHoldingsRepository) DeleteByID(id int) error {
	_, err := r.DB.Exec("DELETE FROM holdings WHERE id = ?", id)
	return err
}
