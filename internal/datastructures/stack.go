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

type GenericStack[T any] struct {
	data []T
}

func (s *GenericStack[T]) Push(t T) {
	s.data = append(s.data, t)
}

func (s *GenericStack[T]) Pop() (T, bool) {
	var zero T
	if len(s.data) == 0 {
		return zero, false
	}

	last := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return last, true
}

func (s *GenericStack[T]) Peek() (T, bool) {
	var zero T
	if len(s.data) == 0 {
		return zero, false
	}
	return s.data[len(s.data)-1], true
}

func (s *GenericStack[T]) Len() int {
	return len(s.data)
}
