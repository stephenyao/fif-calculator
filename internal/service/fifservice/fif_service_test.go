package fifservice

import (
	"fif-calculator/internal/model"
	"testing"
	"time"
)

func TestComputeFRDIncome(t *testing.T) {
	t.Run("No quick sale adjustment", func(t *testing.T) {
		trades := []model.Trade{
			{0, "XYZ", "2021-01-01", 1000, 100, "USD", "buy"},
			{0, "XYZ", "2021-02-01", 1000, 100, "USD", "buy"},
			{0, "GOOG", "2021-02-01", 1500, 100, "USD", "buy"},
			{0, "GOOG", "2021-02-01", 1500, 100, "USD", "buy"},
		}

		start := time.Date(2022, time.April, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2023, time.March, 31, 0, 0, 0, 0, time.UTC)

		got, err := ComputeFRDIncome(trades, start, end)
		var want float64 = 2000*0.05 + 3000*0.05

		if err != nil {
			t.Errorf("ComputeFRDIncome() error = %v", err)
		}

		if got != want {
			t.Errorf("got %f, want %f", got, want)
		}
	})
}
