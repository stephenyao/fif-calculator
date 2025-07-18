package viewmodel

import (
	"context"
	"strings"
)

type NavigationItem struct {
	Title    string
	URL      string
	IsActive bool
}

type NavigationViewModel struct {
	Items       []NavigationItem
	CurrentPath string
}

func NewNavigationViewModel(currentPath string, ctx context.Context) *NavigationViewModel {
	items := []NavigationItem{
		{"Holdings", "/holdings", false},
		{"Cost Basis", "/cost-basis", false},
		{"FIF calculation", "/fif", false},
	}

	_, ok := ctx.Value("uid").(string)

	if !ok {
		items = append(items, NavigationItem{"Login", "/login", false})
	} else {
		items = append(items, NavigationItem{"Account", "/account", false})
	}

	for i := range items {
		if strings.HasPrefix(currentPath, items[i].URL) {
			items[i].IsActive = true
		}
	}

	return &NavigationViewModel{
		Items:       items,
		CurrentPath: currentPath,
	}
}
