package handler

import (
	"fif-calculator/internal/service/fifservice"
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

	// You’ll compute this from the Trades table
	holdings := fifservice.ComputeHoldingsBetween(startDate, endDate)

	fif.RenderFIFHoldingQuantities(holdings, year).Render(r.Context(), w)
}
