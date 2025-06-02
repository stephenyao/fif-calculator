package handler

import (
	. "fif-calculator/internal/model"
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

func (h *FIFHandler) StartFIF(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	yearStr := r.FormValue("financialYear")
	year, _ := strconv.Atoi(yearStr)

	startDate := time.Date(year-1, 4, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, 3, 31, 0, 0, 0, 0, time.UTC)

	trades, err := h.Repo.GetAll()

	if err != nil {
		http.Error(w, "Failed to fetch all trades", http.StatusInternalServerError)
	}

	// You’ll compute this from the Trades table
	holdings, err := fifservice.ComputeHoldingsBetween(trades, startDate, endDate)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fif.RenderFIFHoldingQuantities(holdings, year).Render(r.Context(), w)
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

	trades, err := h.Repo.GetAll()
	if err != nil {
		http.Error(w, "Failed to load trades", http.StatusInternalServerError)
		return
	}

	holdings, err := fifservice.ComputeHoldingsBetween(trades, startDate, endDate)

	var results []FIFResult
	totalFDR := 0.0

	for _, hq := range holdings {
		priceStartStr := r.FormValue("price_start_" + hq.Symbol)
		priceEndStr := r.FormValue("price_end_" + hq.Symbol)

		priceStart, _ := strconv.ParseFloat(priceStartStr, 64)
		priceEnd, _ := strconv.ParseFloat(priceEndStr, 64)

		startVal := hq.QuantityStart * priceStart
		endVal := hq.QuantityEnd * priceEnd
		fdr := startVal * 0.05

		totalFDR += fdr

		results = append(results, FIFResult{
			Symbol:     hq.Symbol,
			StartValue: startVal,
			EndValue:   endVal,
			FDRAmount:  fdr,
		})
	}

	vm := FIFCalculationViewModel{
		Year:     year,
		Results:  results,
		TotalFDR: totalFDR,
	}

	fif.RenderFIFResult(vm).Render(r.Context(), w)
}
