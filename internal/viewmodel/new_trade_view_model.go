package viewmodel

import "github.com/a-h/templ"

type TradeFormViewModel struct {
	Title     string
	ActionURL templ.SafeURL
}

func CreateTradeFormViewModel(title string, actionURL string) TradeFormViewModel {
	return TradeFormViewModel{
		Title:     title,
		ActionURL: templ.SafeURL(actionURL),
	}
}
