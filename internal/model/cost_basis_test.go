package model

import (
	"fif-clacultor/internal/constants"
	"testing"
)

func TestCostBasisBySymbol(t *testing.T) {
	t.Run("buy and sell orders of exact same quantity", func(t *testing.T) {
		var trades []Trade = []Trade{
			{0, "XYZ", "2024/01/01", 15, 100, "USD", constants.Buy},
			{0, "XYZ", "2024/02/01", 10, 100, "USD", constants.Buy},
			{0, "XYZ", "2024/03/01", 30, 300, "USD", constants.Sell},
		}

		costBasisBySymbol := CostBasisBySymbol(trades)
		got, ok := costBasisBySymbol["XYZ"]

		if !ok {
			t.Errorf("not ok")
		}

		if got != 0 {
			t.Errorf("got %f, want 0", got)
		}
	})

	t.Run("buy orders is greater than sell order quantity", func(t *testing.T) {
		var trades []Trade = []Trade{
			{0, "XYZ", "2024/01/01", 10, 100, "USD", constants.Buy},
			{0, "XYZ", "2024/02/01", 20, 100, "USD", constants.Buy},
			{0, "XYZ", "2024/03/01", 10, 300, "USD", constants.Sell},
		}

		costBasisBySymbol := CostBasisBySymbol(trades)
		got, ok := costBasisBySymbol["XYZ"]
		var expected float64 = 2000

		if !ok {
			t.Errorf("not ok")
		}

		if got != expected {
			t.Errorf("got %f, want %f", got, expected)
		}
	})

	t.Run("buy orders is greater than sell order with overflow to next buy order", func(t *testing.T) {
		var trades []Trade = []Trade{
			{0, "XYZ", "2024/01/01", 15, 110, "USD", constants.Buy},
			{0, "XYZ", "2024/02/01", 10, 110, "USD", constants.Buy},
			{0, "XYZ", "2024/03/01", 5, 300, "USD", constants.Sell},
		}

		costBasisBySymbol := CostBasisBySymbol(trades)
		got, ok := costBasisBySymbol["XYZ"]
		var expected float64 = 2200

		if !ok {
			t.Errorf("not ok")
		}

		if got != expected {
			t.Errorf("got %f, want %f", got, expected)
		}
	})

	t.Run("buy orders is greater than sell order with overflow to second buy order", func(t *testing.T) {
		var trades []Trade = []Trade{
			{0, "XYZ", "2024/01/01", 10, 100, "USD", constants.Buy},
			{0, "XYZ", "2024/02/01", 10, 110, "USD", constants.Buy},
			{0, "XYZ", "2024/03/01", 15, 300, "USD", constants.Sell},
		}

		costBasisBySymbol := CostBasisBySymbol(trades)
		got, ok := costBasisBySymbol["XYZ"]
		var expected float64 = 550

		if !ok {
			t.Errorf("not ok")
		}

		if got != expected {
			t.Errorf("got %f, want %f", got, expected)
		}
	})

	t.Run("Multiple buy and sell orders should use FIFO", func(t *testing.T) {
		var trades []Trade = []Trade{
			{0, "XYZ", "2024/01/01", 10, 100, "USD", constants.Buy},
			{0, "XYZ", "2024/02/01", 20, 100, "USD", constants.Buy},
			{0, "XYZ", "2024/03/01", 5, 300, "USD", constants.Sell},
			{0, "XYZ", "2024/04/01", 10, 100, "USD", constants.Buy},
			{0, "XYZ", "2024/05/01", 20, 100, "USD", constants.Buy},
			{0, "XYZ", "2024/06/01", 10, 300, "USD", constants.Sell},
		}
		var expected float64 = 5*100 + 10*100 + 20*100
		costBasisBySymbol := CostBasisBySymbol(trades)
		got, ok := costBasisBySymbol["XYZ"]

		if !ok {
			t.Errorf("not ok")
		}

		if got != expected {
			t.Errorf("got %f, want %f", got, expected)
		}
	})

	t.Run("First order is sell order should not be counted", func(t *testing.T) {
		var trades []Trade = []Trade{
			{0, "XYZ", "2024/01/01", 10, 100, "USD", constants.Sell},
			{0, "XYZ", "2024/02/01", 20, 100, "USD", constants.Buy},
		}
		var expected float64 = 20 * 100
		costBasisBySymbol := CostBasisBySymbol(trades)
		got, ok := costBasisBySymbol["XYZ"]

		if !ok {
			t.Errorf("not ok")
		}

		if got != expected {
			t.Errorf("got %f, want %f", got, expected)
		}
	})

	t.Run("All sell orders", func(t *testing.T) {
		var trades []Trade = []Trade{
			{0, "XYZ", "2024/01/01", 10, 100, "USD", constants.Sell},
			{0, "XYZ", "2024/02/01", 20, 100, "USD", constants.Sell},
		}
		var expected float64 = 0
		costBasisBySymbol := CostBasisBySymbol(trades)
		got, ok := costBasisBySymbol["XYZ"]

		if !ok {
			t.Errorf("not ok")
		}

		if got != expected {
			t.Errorf("got %f, want %f", got, expected)
		}
	})

	t.Run("Two symbols", func(t *testing.T) {
		var trades []Trade = []Trade{
			{0, "XYZ", "2024/01/01", 10, 100, "USD", constants.Buy},
			{0, "XYZ", "2024/02/01", 20, 100, "USD", constants.Buy},
			{0, "XYZ", "2024/03/01", 5, 300, "USD", constants.Sell},
			{0, "GGG", "2024/01/01", 10, 100, "USD", constants.Buy},
			{0, "GGG", "2024/02/01", 20, 100, "USD", constants.Buy},
			{0, "GGG", "2024/03/01", 5, 300, "USD", constants.Sell},
		}

		costBasisBySymbol := CostBasisBySymbol(trades)

		testCases := []struct {
			got  float64
			want float64
		}{
			{got: costBasisBySymbol["XYZ"], want: 2500},
			{got: costBasisBySymbol["GGG"], want: 2500},
		}

		for _, tc := range testCases {
			if tc.got != tc.want {
				t.Errorf("got %f, want %f", tc.got, tc.want)
			}
		}

	})
}
