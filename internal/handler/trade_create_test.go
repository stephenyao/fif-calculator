package handler

import (
	"fif-calculator/internal/model"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestTradeHandler_Create(t *testing.T) {
	form := formValues()
	var received *model.Trade
	mockRepo := &mockTradeRepo{insert: func(trade *model.Trade) error {
		received = trade
		return nil
	}}
	handler := &TradeHandler{Repo: mockRepo}

	req := httptest.NewRequest(http.MethodPost, "/trade/create", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected status %d, got %d", http.StatusSeeOther, rec.Code)
	}

	if received == nil {
		t.Fatal("expected trade to be passed to repository, got nil")
	}
	if received.Symbol != "AAPL" || received.BuyDate != "2023-07-01" ||
		received.Quantity != 10 || received.Price != 150.50 ||
		received.Currency != "USD" || received.Action != "buy" {
		t.Errorf("unexpected trade received: %+v", received)
	}
}

func formValues() url.Values {
	form := url.Values{}
	form.Set("symbol", "AAPL")
	form.Set("buyDate", "2023-07-01")
	form.Set("quantity", "10")
	form.Set("price", "150.50")
	form.Set("currency", "USD")
	form.Set("action", "buy")
	return form
}

type mockTradeRepo struct {
	insert func(trade *model.Trade) error
}

func (m *mockTradeRepo) Insert(trade *model.Trade) error {
	return m.insert(trade)
}

// Stub methods for interface compliance
func (m *mockTradeRepo) Update(trade *model.Trade) error      { return nil }
func (m *mockTradeRepo) Delete(id int) error                  { return nil }
func (m *mockTradeRepo) GetByID(id int) (*model.Trade, error) { return nil, nil }
func (m *mockTradeRepo) GetAll() ([]*model.Trade, error)      { return nil, nil }
func (m *mockTradeRepo) DeleteByID(id int) error              { return nil }
func (m *mockTradeRepo) GetBySymbol(symbol string) ([]*model.Trade, error) {
	return nil, nil
}
