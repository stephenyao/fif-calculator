package fifhandler

import (
	"fif-calculator/internal/model"
	"fif-calculator/internal/service/costbasisservice"
	"fif-calculator/internal/service/fifservice"
	. "fif-calculator/internal/viewmodel"
	"fif-calculator/views/fif"
	"net/http"
	"strconv"
	"time"
)

func (h *FIFHandler) Start(w http.ResponseWriter, r *http.Request) {
	err := fif.Start(r.URL.Path).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Could not render start fif page", http.StatusInternalServerError)
	}
}

func (h *FIFHandler) HoldingsInfo(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		http.Error(w, "Could not parse start fif form", http.StatusInternalServerError)
	}

	yearStr := r.FormValue("financialYear")
	year, _ := strconv.Atoi(yearStr)

	startDate := time.Date(year-1, 4, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, 3, 31, 0, 0, 0, 0, time.UTC)

	trades, err := h.Repo.GetAll()
	if err != nil {
		http.Error(w, "Failed to fetch all trades", http.StatusInternalServerError)
	}

	maxCostBasis := costbasisservice.MaxCostBasisDuringYear(trades, startDate, endDate)
	if maxCostBasis < 50000 {
		fif.NoFIFApplicable().Render(r.Context(), w)
		return
	}

	// You’ll compute this from the Trades table
	holdings, err := fifservice.ComputeHoldingsBetween(trades, startDate, endDate)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = fif.RenderFIFHoldingQuantities(holdings, year).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Failed to render fif", http.StatusInternalServerError)
	}
}

func (h *FIFHandler) CalculateFIF(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	yearStr := r.FormValue("financialYear")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "Invalid financial year", http.StatusBadRequest)
		return
	}

	startDate := time.Date(year-1, 4, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, 3, 31, 0, 0, 0, 0, time.UTC)

	trades, err := h.Repo.GetAllByAscendingDate()
	if err != nil {
		http.Error(w, "Failed to load trades", http.StatusInternalServerError)
		return
	}

	holdings, err := fifservice.ComputeHoldingsBetween(trades, startDate, endDate)

	for _, hq := range holdings {
		priceStartStr := r.FormValue("price_start_" + hq.Symbol)
		priceEndStr := r.FormValue("price_end_" + hq.Symbol)

		priceStart, _ := strconv.ParseFloat(priceStartStr, 64)
		priceEnd, _ := strconv.ParseFloat(priceEndStr, 64)

		hq.OpeningPrice = priceStart
		hq.ClosingPrice = priceEnd

		gainLossParams, _ := getGainLossParams(hq.Symbol, r)

		hq.GainLoss = gainLossParams
	}

	frdResults, err := fifservice.ComputeFRDIncome(trades, holdings, startDate, endDate)
	var totalFDR float64 = 0

	for _, result := range frdResults {
		totalFDR += result.TotalFDRIncome()
	}

	var totalCV float64 = 0
	cvResults, err := fifservice.ComputeCVIncome(trades, holdings, startDate, endDate)

	for _, result := range cvResults {
		totalCV += result.TotalIncome()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	vm := FIFCalculationViewModel{
		Year:       year,
		FDRResults: frdResults,
		TotalFDR:   totalFDR,
		CVResults:  cvResults,
		TotalCV:    totalCV,
	}

	err = fif.RenderFIFResult(vm).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Failed to render fif", http.StatusInternalServerError)
	}
}

func getGainLossParams(symbol string, r *http.Request) (model.GainLossParams, error) {
	err := r.ParseForm()

	if err != nil {
		return model.GainLossParams{}, err
	}

	dividends, err := strconv.ParseFloat(r.FormValue("dividends_"+symbol), 64)
	taxCredits, err := strconv.ParseFloat(r.FormValue("tax_credits"+symbol), 64)
	otherGains, err := strconv.ParseFloat(r.FormValue("other_gains"+symbol), 64)
	foreignIncomeTax, err := strconv.ParseFloat(r.FormValue("foreign_income_tax"+symbol), 64)
	otherCosts, err := strconv.ParseFloat(r.FormValue("other_costs"+symbol), 64)

	if err != nil {
		return model.GainLossParams{}, err
	}

	return model.GainLossParams{
		Dividends:        dividends,
		TaxCredits:       taxCredits,
		OtherGains:       otherGains,
		ForeignIncomeTax: foreignIncomeTax,
		OtherCosts:       otherCosts,
	}, nil
}
