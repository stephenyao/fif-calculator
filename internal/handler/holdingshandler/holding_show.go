package holdingshandler

import (
	"database/sql"
	"errors"
	"fif-calculator/internal/model"
	"fif-calculator/internal/service/costbasisservice"
	"fif-calculator/internal/utils"
	"fif-calculator/internal/viewmodel"
	"fif-calculator/views/holdings"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

const pageLimit = 10

func (h *HoldingsHandler) Show(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(param)

	if err != nil {
		http.Error(w, "id should be a number", http.StatusBadRequest)
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))

	if err != nil {
		page = 0
	}

	userId, err := utils.GetUID(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	holding, err := h.HoldingsRepository.GetHolding(id, userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Holding not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	uid, err := utils.GetUID(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	trades, err := h.TradeRepository.GetByHoldingID(id, uid)
	paginatedTrades, _, err := h.TradeRepository.GetByHoldingIDPaginated(id, pageLimit, page, uid)
	costBasis := costbasisservice.CostBasisBySymbol(trades, nil)[holding.Symbol]

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := (len(trades) + pageLimit - 1) / pageLimit
	startPage, endPage := calculateStartEndPageIndices(page, totalPages)

	holdingVm := viewmodel.HoldingViewModel{
		ID:              id,
		Name:            holding.Name,
		Symbol:          holding.Symbol,
		Currency:        holding.Currency,
		TotalTrades:     fmt.Sprintf("%d", costBasis.TotalTrades),
		CurrentQuantity: strconv.FormatFloat(costBasis.CurrentQuantity, 'f', -1, 64),
		CostBasis:       fmt.Sprintf("$%.2f", costBasis.CostBasisNZD),
	}

	tradesVm := viewmodel.TradesViewModel{
		HoldingID: id,
		Trades:    convertTradesToViewModel(param, paginatedTrades),
		PageInfo: viewmodel.PageInfo{
			TotalPages:   totalPages,
			CurrentPage:  page,
			StartPage:    startPage,
			EndPage:      endPage,
			PreviousPage: max(page-1, 0),
			NextPage:     min(page+1, totalPages-1),
		},
	}

	err = holdings.ViewHolding(r.URL.Path, holdingVm, tradesVm).Render(r.Context(), w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *HoldingsHandler) GetHoldingTrades(w http.ResponseWriter, r *http.Request) {
	// Parse params
	uid, err := utils.GetUID(r.Context())
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	page, err := strconv.Atoi(r.URL.Query().Get("page"))

	count, err := h.TradeRepository.CountTrades(id, uid)
	paginatedTrades, _, err := h.TradeRepository.GetByHoldingIDPaginated(id, pageLimit, page, uid)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	totalPages := (count + pageLimit - 1) / pageLimit
	startPage, endPage := calculateStartEndPageIndices(page, totalPages)

	vm := viewmodel.TradesViewModel{
		HoldingID: id,
		Trades:    convertTradesToViewModel(idStr, paginatedTrades),
		PageInfo: viewmodel.PageInfo{
			TotalPages:   totalPages,
			CurrentPage:  page,
			StartPage:    startPage,
			EndPage:      endPage,
			PreviousPage: max(page-1, 0),
			NextPage:     min(page+1, totalPages-1),
		},
	}

	// If this is an htmx request, return ONLY the panel
	if r.Header.Get("HX-Request") == "true" {
		if err := holdings.TradesPanel(vm).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	//
	// Otherwise render the full page
	// Non-HTMX: redirect to the holding details page
	http.Redirect(w, r, "/holdings/"+idStr, http.StatusSeeOther)
}

func calculateStartEndPageIndices(currentPage int, totalPages int) (int, int) {
	offset := 2
	start := currentPage - offset
	end := currentPage + offset

	if start < 0 {
		start = 0
	}

	if end > totalPages-1 {
		end = totalPages - 1
	}

	return start, end
}

func convertTradesToViewModel(holdingID string, trades []model.Trade) []viewmodel.TradeViewModel {
	viewTrades := make([]viewmodel.TradeViewModel, len(trades))

	for i, trade := range trades {
		viewTrades[i] = viewmodel.TradeViewModel{
			TransactionDate: trade.BuyDate,
			Quantity:        trade.Quantity,
			Price:           trade.Price,
			Currency:        trade.Currency,
			Action:          trade.Action,
			URL:             "/holdings/" + holdingID + "/trades/" + strconv.Itoa(trade.ID),
			BackURL:         "/holdings/" + holdingID,
		}
	}

	return viewTrades
}
