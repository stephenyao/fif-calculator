package viewmodel

import (
	"fif-calculator/internal/model"
	"strconv"
)

type Currency struct {
	Code string
	Name string
}

func (c Currency) DisplayValue() string {
	return c.Code + " - " + c.Name
}

type HoldingFormViewModel struct {
	ID               string
	Name             string
	Ticker           string
	SelectedCurrency string
	Currencies       []Currency
	DefaultCurrency  Currency
}

var defaultCurrency = Currency{Code: "USD", Name: "US Dollar"}

var currencies = []Currency{
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
}

func NewHoldingFormViewModel() *HoldingFormViewModel {
	defaultCurrency := Currency{Code: "USD", Name: "US Dollar"}

	return &HoldingFormViewModel{
		Currencies:      currencies,
		DefaultCurrency: defaultCurrency,
	}
}

func NewHoldingFormViewModelFromRecord(record *model.HoldingRecord) *HoldingFormViewModel {
	defaultCurrency := Currency{Code: "USD", Name: "US Dollar"}

	return &HoldingFormViewModel{
		ID:               strconv.Itoa(record.ID),
		Name:             record.Name,
		Ticker:           record.Symbol,
		SelectedCurrency: record.Currency,
		Currencies:       currencies,
		DefaultCurrency:  defaultCurrency,
	}
}

type HoldingViewModel struct {
	ID              int
	Name            string
	Symbol          string
	Currency        string
	CostBasis       string
	CurrentQuantity string
	TotalTrades     string
	Trades          []TradeViewModel
	PageInfo        PageInfo
}

type TradeViewModel struct {
	TransactionDate string
	Quantity        float64
	Price           float64
	Currency        string
	Action          string
	URL             string
	BackURL         string
}

type PageInfo struct {
	TotalPages   int
	CurrentPage  int
	StartPage    int
	EndPage      int
	PreviousPage int
	NextPage     int
}
