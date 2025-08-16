package fifhandler

import (
	"fif-calculator/internal/service/costbasisservice"
	"fif-calculator/internal/service/fifservice"
	"fif-calculator/internal/utils"
	"fif-calculator/views/fif"
	"log"
	"net/http"
	"strconv"
)

func (h *FIFHandler) HoldingsInfo(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Could not parse start fif form", http.StatusInternalServerError)
	}
	yearStr := r.FormValue("financialYear")
	year, _ := strconv.Atoi(yearStr)
	startDate, endDate := fifservice.StartEndDates(year)

	userID, err := utils.GetUID(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	trades, err := h.TradeRepository.GetAllByAscendingDate(userID)

	if err != nil {
		http.Error(w, "Failed to fetch all trades", http.StatusInternalServerError)
		return
	}

	if !costbasisservice.IsEligibleForFIF(trades, startDate, endDate) {
		err = fif.NoFIFApplicable().Render(r.Context(), w)
		return
	}
	holdings, err := fifservice.ComputeHoldingsBetween(trades, startDate, endDate)
	if err != nil {
		log.Printf("%v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = fif.RenderFIFHoldingQuantities(holdings, year).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render fif", http.StatusInternalServerError)
		return
	}
}
