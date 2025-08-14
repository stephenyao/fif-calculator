package viewmodel

type SelectOptions struct {
	Options []Option
}

type Option struct {
	Value   string
	Display string
}
