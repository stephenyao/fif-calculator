package datastructures

import . "fif-calculator/internal/model"

type Stack struct {
	data []Trade
}

func (s *Stack) Push(trade Trade) {
	s.data = append(s.data, trade)
}

func (s *Stack) Pop() (Trade, bool) {
	if len(s.data) == 0 {
		var zero Trade
		return zero, false
	}
	last := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return last, true
}

func (s *Stack) Peek() (Trade, bool) {
	if len(s.data) == 0 {
		var zero Trade
		return zero, false
	}
	return s.data[len(s.data)-1], true
}

func (s *Stack) Len() int {
	return len(s.data)
}
