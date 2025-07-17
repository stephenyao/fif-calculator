package viewmodel

type TradeFormViewModel struct {
	Title     string
	ActionURL string
}

func CreateTradeFormViewModel(title string, actionURL string) TradeFormViewModel {
	return TradeFormViewModel{
		Title:     title,
		ActionURL: actionURL,
	}
}
