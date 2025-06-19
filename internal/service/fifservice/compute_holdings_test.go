package fifservice

import (
	"fif-calculator/internal/constants"
	"fif-calculator/internal/model"
	"testing"
	"time"
)

func mustDate(s string) time.Time {
	d, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return d
}

func TestComputeHoldingsBetween(t *testing.T) {
	start := mustDate("2024-04-01")
	end := mustDate("2025-03-31")

	t.Run("compute holdings", func(t *testing.T) {
		cases := []struct {
			name   string
			trades []model.Trade
			want   map[string]model.HoldingQuantity
		}{
			{
				name: "only buy trades before start and end",
				trades: []model.Trade{
					{Symbol: "XYZ", BuyDate: "2024-03-01", Quantity: 10, Action: constants.Buy},
					{Symbol: "XYZ", BuyDate: "2024-03-15", Quantity: 5, Action: constants.Buy},
				},
				want: map[string]model.HoldingQuantity{
					"XYZ": {Symbol: "XYZ", QuantityStart: 15, QuantityEnd: 15, NumberOfTrades: 0},
				},
			},
			{
				name: "buy and sell before and after start",
				trades: []model.Trade{
					{Symbol: "XYZ", BuyDate: "2024-03-01", Quantity: 10, Action: constants.Buy},
					{Symbol: "XYZ", BuyDate: "2024-04-10", Quantity: 5, Action: constants.Buy},
					{Symbol: "XYZ", BuyDate: "2024-05-01", Quantity: 4, Action: constants.Sell},
				},
				want: map[string]model.HoldingQuantity{
					"XYZ": {Symbol: "XYZ", QuantityStart: 10, QuantityEnd: 11, NumberOfTrades: 2},
				},
			},
			{
				name: "negative quantity should be skipped",
				trades: []model.Trade{
					{Symbol: "XYZ", BuyDate: "2024-03-01", Quantity: 5, Action: constants.Sell},
				},
				want: map[string]model.HoldingQuantity{},
			},
			{
				name:   "empty trade list",
				trades: []model.Trade{},
				want:   map[string]model.HoldingQuantity{},
			},
			{
				name: "multiple symbols",
				trades: []model.Trade{
					{Symbol: "ABC", BuyDate: "2024-03-01", Quantity: 10, Action: constants.Buy},
					{Symbol: "XYZ", BuyDate: "2024-02-01", Quantity: 5, Action: constants.Buy},
					{Symbol: "XYZ", BuyDate: "2024-06-01", Quantity: 10, Action: constants.Buy},
				},
				want: map[string]model.HoldingQuantity{
					"ABC": {Symbol: "ABC", QuantityStart: 10, QuantityEnd: 10, NumberOfTrades: 0},
					"XYZ": {Symbol: "XYZ", QuantityStart: 5, QuantityEnd: 15, NumberOfTrades: 1},
				},
			},
			{
				name: "symbol with 0 net quantity and no trades in period is excluded",
				trades: []model.Trade{
					{Symbol: "ABC", BuyDate: "2024-03-01", Quantity: 10, Action: constants.Buy},
					{Symbol: "ABC", BuyDate: "2024-03-02", Quantity: 10, Action: constants.Sell},
					{Symbol: "XYZ", BuyDate: "2024-02-01", Quantity: 5, Action: constants.Buy},
					{Symbol: "XYZ", BuyDate: "2024-06-01", Quantity: 10, Action: constants.Buy},
				},
				want: map[string]model.HoldingQuantity{
					"XYZ": {Symbol: "XYZ", QuantityStart: 5, QuantityEnd: 15, NumberOfTrades: 1},
				},
			},
			{
				name: "trades in period but sold all",
				trades: []model.Trade{
					{Symbol: "ABC", BuyDate: "2024-05-01", Quantity: 10, Action: constants.Buy},
					{Symbol: "ABC", BuyDate: "2024-05-02", Quantity: 10, Action: constants.Sell},
					{Symbol: "XYZ", BuyDate: "2024-02-01", Quantity: 5, Action: constants.Buy},
					{Symbol: "XYZ", BuyDate: "2024-06-01", Quantity: 10, Action: constants.Buy},
				},
				want: map[string]model.HoldingQuantity{
					"ABC": {Symbol: "ABC", QuantityStart: 0, QuantityEnd: 0, NumberOfTrades: 2},
					"XYZ": {Symbol: "XYZ", QuantityStart: 5, QuantityEnd: 15, NumberOfTrades: 1},
				},
			},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				holdings, err := ComputeHoldingsBetween(tt.trades, start, end)
				if err != nil {
					t.Fatal(err)
				}

				if len(holdings) != len(tt.want) {
					t.Fatalf("expected %d holdings, got %d", len(tt.want), len(holdings))
				}

				gotMap := make(map[string]*model.HoldingQuantity)
				for _, h := range holdings {
					gotMap[h.Symbol] = h
				}

				for symbol, want := range tt.want {
					got, ok := gotMap[symbol]
					if !ok {
						t.Errorf("missing symbol %s", symbol)
						continue
					}
					if got.QuantityStart != want.QuantityStart {
						t.Errorf("%s QuantityStart: got %.2f, want %.2f", symbol, got.QuantityStart, want.QuantityStart)
					}
					if got.QuantityEnd != want.QuantityEnd {
						t.Errorf("%s QuantityEnd: got %.2f, want %.2f", symbol, got.QuantityEnd, want.QuantityEnd)
					}
					if got.NumberOfTrades != want.NumberOfTrades {
						t.Errorf("%s NumberOfTrades: got %d, want %d", symbol, got.NumberOfTrades, want.NumberOfTrades)
					}
				}
			})
		}
	})

	t.Run("compute holdings with sales and purchases", func(t *testing.T) {
		tests := []struct {
			name   string
			trades []model.Trade
			want   map[string]struct {
				QuantityStart     float64
				QuantityEnd       float64
				NumberOfTrades    int
				CostOfPurchases   float64
				ProceedsFromSales float64
			}
		}{
			{
				name: "simple buy and sell within period",
				trades: []model.Trade{
					{Symbol: "XYZ", BuyDate: "2024-04-10", Quantity: 10, Price: 100, Action: constants.Buy},
					{Symbol: "XYZ", BuyDate: "2024-06-01", Quantity: 5, Price: 120, Action: constants.Sell},
				},
				want: map[string]struct {
					QuantityStart     float64
					QuantityEnd       float64
					NumberOfTrades    int
					CostOfPurchases   float64
					ProceedsFromSales float64
				}{
					"XYZ": {
						QuantityStart:     0,
						QuantityEnd:       5,
						NumberOfTrades:    2,
						CostOfPurchases:   1000,
						ProceedsFromSales: 600,
					},
				},
			},
			{
				name: "buys before and inside range, one sell inside range",
				trades: []model.Trade{
					{Symbol: "ABC", BuyDate: "2024-03-15", Quantity: 5, Price: 100, Action: constants.Buy},  // before start
					{Symbol: "ABC", BuyDate: "2024-04-15", Quantity: 10, Price: 200, Action: constants.Buy}, // inside
					{Symbol: "ABC", BuyDate: "2024-05-01", Quantity: 8, Price: 150, Action: constants.Sell}, // inside
				},
				want: map[string]struct {
					QuantityStart     float64
					QuantityEnd       float64
					NumberOfTrades    int
					CostOfPurchases   float64
					ProceedsFromSales float64
				}{
					"ABC": {
						QuantityStart:     5,
						QuantityEnd:       7, // 5 (before) + 10 (buy) - 8 (sell)
						NumberOfTrades:    2,
						CostOfPurchases:   2000, // 10 × 200
						ProceedsFromSales: 1200, // 8 × 150
					},
				},
			},
			{
				name: "multiple symbols",
				trades: []model.Trade{
					{Symbol: "AAA", BuyDate: "2024-03-01", Quantity: 10, Price: 90, Action: constants.Buy},
					{Symbol: "XYZ", BuyDate: "2024-02-01", Quantity: 5, Price: 50, Action: constants.Buy},
					{Symbol: "XYZ", BuyDate: "2024-06-01", Quantity: 10, Price: 70, Action: constants.Buy},
					{Symbol: "XYZ", BuyDate: "2024-07-01", Quantity: 4, Price: 60, Action: constants.Sell},
				},
				want: map[string]struct {
					QuantityStart     float64
					QuantityEnd       float64
					NumberOfTrades    int
					CostOfPurchases   float64
					ProceedsFromSales float64
				}{
					"AAA": {
						QuantityStart:     10,
						QuantityEnd:       10,
						NumberOfTrades:    0,
						CostOfPurchases:   0,
						ProceedsFromSales: 0,
					},
					"XYZ": {
						QuantityStart:     5,
						QuantityEnd:       11,
						NumberOfTrades:    2,
						CostOfPurchases:   700, // 10 × 70
						ProceedsFromSales: 240, // 4 × 60
					},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				holdings, err := ComputeHoldingsBetween(tt.trades, start, end)
				if err != nil {
					t.Fatal(err)
				}

				if len(holdings) != len(tt.want) {
					t.Fatalf("expected %d holdings, got %d", len(tt.want), len(holdings))
				}

				for _, h := range holdings {
					expect := tt.want[h.Symbol]
					if h.QuantityStart != expect.QuantityStart {
						t.Errorf("%s start qty: got %.2f, want %.2f", h.Symbol, h.QuantityStart, expect.QuantityStart)
					}
					if h.QuantityEnd != expect.QuantityEnd {
						t.Errorf("%s end qty: got %.2f, want %.2f", h.Symbol, h.QuantityEnd, expect.QuantityEnd)
					}
					if h.NumberOfTrades != expect.NumberOfTrades {
						t.Errorf("%s trade count: got %d, want %d", h.Symbol, h.NumberOfTrades, expect.NumberOfTrades)
					}
					if h.CostOfPurchases != expect.CostOfPurchases {
						t.Errorf("%s cost of purchases: got %.2f, want %.2f", h.Symbol, h.CostOfPurchases, expect.CostOfPurchases)
					}
					if h.ProceedsFromSales != expect.ProceedsFromSales {
						t.Errorf("%s proceeds from sales: got %.2f, want %.2f", h.Symbol, h.ProceedsFromSales, expect.ProceedsFromSales)
					}
				}
			})
		}
	})

}
