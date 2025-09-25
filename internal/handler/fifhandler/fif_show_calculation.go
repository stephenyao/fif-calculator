package fifhandler

import (
	"fif-calculator/internal/repository"
	"fif-calculator/internal/service/fifservice"
	"fif-calculator/internal/utils"
	"fif-calculator/views/fif"
	"net/http"
	"strconv"
	"time"
)

func (h *FIFHandler) ShowCalculation(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUID(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	trades, err := h.TradeRepository.GetAllByAscendingDate(userID)
	if err != nil {
		http.Error(w, "Failed to load trades", http.StatusInternalServerError)
		return
	}

	holdings, err := fifservice.ComputeHoldingsBetween(trades, year)
	var fdrInputs []fifservice.FDRHoldingInput

	for _, hq := range holdings {
		priceStartStr := r.FormValue("price_start_" + hq.Symbol)
		priceStartExStr := r.FormValue("price_start_exchange_" + hq.Symbol)
		priceEndStr := r.FormValue("price_end_" + hq.Symbol)
		priceEndExStr := r.FormValue("price_end_exchange_" + hq.Symbol)

		priceStart, _ := strconv.ParseFloat(priceStartStr, 64)
		priceEnd, _ := strconv.ParseFloat(priceEndStr, 64)
		priceStartEx, _ := strconv.ParseFloat(priceStartExStr, 64)
		pricenEndEx, _ := strconv.ParseFloat(priceEndExStr, 64)

		hq.OpeningPrice = priceStart * priceStartEx
		hq.ClosingPrice = priceEnd * pricenEndEx

		gainLossParams, _ := getGainLossParams(hq.Symbol, r)

		hq.GainLoss = gainLossParams

		fdrInputs = append(fdrInputs, fifservice.FDRHoldingInput{
			OpeningPrice:      priceStart,
			ExchangeRateToNZD: priceStartEx,
			HoldingID:         repository.HoldingID(hq.HoldingId),
		})
	}

	fdrInput := fifservice.FDRInput{Holdings: fdrInputs}
	results := h.Service.FDRIncome(fdrInput, startDate, endDate)

	err = fif.ShowCalculation(r.URL.Path, results).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
