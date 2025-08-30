package fifhandler

import (
	"fif-calculator/internal/service/costbasisservice"
	"fif-calculator/internal/service/fifservice"
	"fif-calculator/internal/utils"
	"fif-calculator/internal/viewmodel"
	"fif-calculator/views/fif"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

func (h *FIFHandler) GetHolding(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	year, err := strconv.Atoi(chi.URLParam(r, "year"))

	if err != nil {
		log.Println(err)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	userID, err := utils.GetUID(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	trades, err := h.TradeRepository.GetByHoldingID(id, userID)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	holdings, err := fifservice.ComputeHoldingsBetween(trades, year)
	if err != nil {
		log.Printf("%v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = fif.RenderHoldingDetail(holdings[0], year).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render fif", http.StatusInternalServerError)
		return
	}
}

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

	isEligible := costbasisservice.IsEligibleForFIF(trades, startDate, endDate)

	if !isEligible {
		_ = fif.NoFIFApplicable().Render(r.Context(), w)
		return
	}

	holdings, err := fifservice.ComputeHoldingsBetween(trades, year)
	if err != nil {
		log.Printf("%v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vm := viewmodel.CreateFIFHoldingsViewModel(holdings, year)
	err = fif.RenderFIFHoldingQuantities(vm).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render fif", http.StatusInternalServerError)
		return
	}
}
