package viewmodel

import "github.com/a-h/templ"

type TradeFormViewModel struct {
	Title     string
	ActionURL templ.SafeURL
}

func CreateTradeFormViewModel(holding string, actionURL string) TradeFormViewModel {
	title := "New Trade for " + holding
	return TradeFormViewModel{
		Title:     title,
		ActionURL: templ.SafeURL(actionURL),
	}
}
