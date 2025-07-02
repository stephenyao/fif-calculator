package fifhandler

import (
	"fif-calculator/internal/model"
	"fif-calculator/views/fif"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *FIFHandler) View(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid calculation ID", http.StatusBadRequest)
		return
	}

	calc, holdings, err := h.FIFRepository.GetCalculationWithHoldings(id)
	if err != nil {
		http.Error(w, "Calculation not found", http.StatusNotFound)
		return
	}

	holdingsInfos := convertToHoldingInfos(holdings)

	err = fif.View(r.URL.Path, holdingsInfos, calc.FinancialYear).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Failed to render fif", http.StatusInternalServerError)
	}
}

func convertToHoldingInfos(fifHoldings []*model.FIFHolding) []*model.HoldingInfo {
	var infos []*model.HoldingInfo
	for _, h := range fifHoldings {
		infos = append(infos, &model.HoldingInfo{
			Symbol:            h.Symbol,
			QuantityStart:     h.QuantityStart,
			QuantityEnd:       h.QuantityEnd,
			OpeningPrice:      h.PriceStart,
			ClosingPrice:      h.PriceEnd,
			ProceedsFromSales: h.ProceedsFromSales,
			CostOfPurchases:   h.CostOfPurchases,
			GainLoss: model.GainLossParams{
				Dividends:        h.Dividends,
				TaxCredits:       h.TaxCredits,
				OtherGains:       h.OtherGains,
				ForeignIncomeTax: h.ForeignIncomeTax,
				OtherCosts:       h.OtherCosts,
			},
		})
	}
	return infos
}
