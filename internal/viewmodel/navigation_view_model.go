package viewmodel

import "strings"

type NavigationItem struct {
	Title    string
	URL      string
	IsActive bool
}

type NavigationViewModel struct {
	Items       []NavigationItem
	CurrentPath string
}

func NewNavigationViewModel(currentPath string) *NavigationViewModel {
	items := []NavigationItem{
		{"Holdings", "/holdings", false},
		{"Cost Basis", "/cost-basis", false},
		{"FIF calculation", "/fif", false},
		{"Login", "/login", false},
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
