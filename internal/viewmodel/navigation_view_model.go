package viewmodel

type NavigationItem struct {
	Title string
	URL   string
}

type NavigationViewModel struct {
	Items       []NavigationItem
	CurrentPath string
}

func NewNavigationViewModel() *NavigationViewModel {
	return &NavigationViewModel{
		Items: []NavigationItem{
			NavigationItem{"Manage Trades", "/trades"},
			NavigationItem{"Cost Basis", "/cost-basis"},
			NavigationItem{"FIF calculation", "/fif-calculation"},
		},
		CurrentPath: "/trades",
	}
}
