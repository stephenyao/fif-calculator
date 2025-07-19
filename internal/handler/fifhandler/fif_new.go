package fifhandler

import (
	"fif-calculator/internal/model"
	"fif-calculator/internal/service/costbasisservice"
	"fif-calculator/internal/service/fifservice"
	"fif-calculator/internal/utils"
	. "fif-calculator/internal/viewmodel"
	"fif-calculator/views/fif"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (h *FIFHandler) New(w http.ResponseWriter, r *http.Request) {
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
	startDate, endDate := fifservice.StartEndDates(year)

	userID, err := utils.GetUID(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	trades, err := h.TradeRepository.GetAllByAscendingDate(userID)
	if err != nil {
		http.Error(w, "Failed to fetch all trades", http.StatusInternalServerError)
	}

	if !costbasisservice.IsEligibleForFIF(trades, startDate, endDate) {
		err = fif.NoFIFApplicable().Render(r.Context(), w)
		return
	}

	holdings, err := fifservice.ComputeHoldingsBetween(trades, startDate, endDate)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = fif.RenderFIFHoldingQuantities(holdings, year).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Failed to render fif", http.StatusInternalServerError)
	}
}

func (h *FIFHandler) FIFFormSubmit(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	switch r.FormValue("action") {
	case "calculate":
		h.calculateFIF(w, r)
	case "save":
		h.saveFIFCalculation(w, r)
	}
}

func (h *FIFHandler) calculateFIF(w http.ResponseWriter, r *http.Request) {
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

func (h *FIFHandler) saveFIFCalculation(w http.ResponseWriter, r *http.Request) {
	year, err := strconv.Atoi(r.FormValue("financialYear"))
	if err != nil {
		http.Error(w, "invalid financial year", http.StatusBadRequest)
		return
	}

	// Step 1: Gather all symbols
	var symbols []string
	for key := range r.Form {
		if strings.HasPrefix(key, "symbol_") {
			symbol := r.FormValue(key)
			symbols = append(symbols, symbol)
		}
	}

	// Step 2: Extract holdings
	var holdings []*model.FIFHolding
	for _, symbol := range symbols {
		get := func(field string) float64 {
			v := r.FormValue(field + "_" + symbol)
			f, _ := strconv.ParseFloat(v, 64)
			return f
		}

		holdings = append(holdings, &model.FIFHolding{
			Symbol:            symbol,
			QuantityStart:     get("quantity_start"),
			QuantityEnd:       get("quantity_end"),
			PriceStart:        get("price_start"),
			PriceEnd:          get("price_end"),
			ProceedsFromSales: get("proceeds_from_sales"),
			Dividends:         get("dividends"),
			TaxCredits:        get("tax_credits"),
			OtherGains:        get("other_gains"),
			CostOfPurchases:   get("cost_of_purchases"),
			ForeignIncomeTax:  get("foreign_income_tax"),
			OtherCosts:        get("other_cost"),
		})
	}

	calc := &model.FIFCalculation{
		UserID:        1, // Or extract from session/context
		FinancialYear: year,
		CalculatedAt:  time.Now(),
	}

	if err := h.FIFRepository.CreateOrUpdateCalculation(calc, holdings); err != nil {
		http.Error(w, "failed to save FIF calculation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Redirect", "/fif")
	w.WriteHeader(http.StatusOK)
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
