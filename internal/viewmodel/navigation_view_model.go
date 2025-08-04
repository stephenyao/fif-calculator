package viewmodel

import (
	"context"
	"strings"
)

type NavigationItem struct {
	Title    string
	URL      string
	Icon     string // SVG path for icon
	IsActive bool
}

type NavigationViewModel struct {
	Items       []NavigationItem
	CurrentPath string
}

func NewNavigationViewModel(currentPath string, ctx context.Context) *NavigationViewModel {
	// Example icons (SVG path `d` values)
	const (
		holdingsIcon = "M3 3h18v18H3V3z"       // simple square
		fifIcon      = "M12 4v16m8-8H4"        // plus shape
		loginIcon    = "M15 12H3m6-6l-6 6 6 6" // arrow into box
		accountIcon  = "M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"
	)

	items := []NavigationItem{
		{Title: "Holdings", URL: "/holdings", Icon: holdingsIcon},
		{Title: "FIF calculation", URL: "/fif", Icon: fifIcon},
	}

	_, ok := ctx.Value("uid").(string)
	if !ok {
		items = append(items, NavigationItem{Title: "Login", URL: "/login", Icon: loginIcon})
	} else {
		items = append(items, NavigationItem{Title: "Account", URL: "/account", Icon: accountIcon})
	}

	// Set active item
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
