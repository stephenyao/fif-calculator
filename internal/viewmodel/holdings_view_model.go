package viewmodel

type Currency struct {
	Code string
	Name string
}

func (c Currency) DisplayValue() string {
	return c.Code + " - " + c.Name
}

type HoldingFormViewModel struct {
	Currencies      []Currency
	DefaultCurrency Currency
}

func NewHoldingFormViewModel() *HoldingFormViewModel {
	defaultCurrency := Currency{Code: "USD", Name: "US Dollar"}

	return &HoldingFormViewModel{
		Currencies: []Currency{
			defaultCurrency,
			{Code: "AUD", Name: "Australian Dollar"},
			{Code: "EUR", Name: "Euro"},
			{Code: "GBP", Name: "British Pound"},
			{Code: "JPY", Name: "Japanese Yen"},
			{Code: "CAD", Name: "Canadian Dollar"},
			{Code: "CHF", Name: "Swiss Franc"},
			{Code: "SGD", Name: "Singapore Dollar"},
			{Code: "HKD", Name: "Hong Kong Dollar"},
			{Code: "KRW", Name: "South Korean Won"},
			{Code: "INR", Name: "Indian Rupee"},
		},
		DefaultCurrency: defaultCurrency,
	}
}

type HoldingViewModel struct {
	ID       int
	Name     string
	Symbol   string
	Currency string
}
