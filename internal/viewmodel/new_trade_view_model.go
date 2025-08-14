package viewmodel

import "fif-calculator/internal/constants"

type TradeFormViewModel struct {
	Title        string
	ActionURL    string
	Actions      SelectOptions
	SelectAction string
}

func CreateTradeFormViewModel(title, actionURL, selectedAction string) TradeFormViewModel {
	return TradeFormViewModel{
		Title:     title,
		ActionURL: actionURL,
		Actions: SelectOptions{
			Options: []Option{
				Option{
					Value:   constants.Buy,
					Display: "Buy",
				},
				Option{
					Value:   constants.Sell,
					Display: "Sell",
				},
			},
		},
		SelectAction: selectedAction,
	}
}
