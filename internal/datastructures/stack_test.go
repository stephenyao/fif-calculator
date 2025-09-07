package datastructures

import (
	"fif-calculator/internal/model"
	"testing"
)

func TestStack_Push(t *testing.T) {
	trade := model.Trade{
		ID: 0,
	}
	stack := &Stack{}

	stack.Push(trade)
	if stack.Len() != 1 {
		t.Errorf("Stack should have length 1")
	}
}

func TestStack_Pop(t *testing.T) {
	trade := model.Trade{
		ID: 0,
	}
	stack := &Stack{}

	stack.Push(trade)

	pop, success := stack.Pop()

	if !success {
		t.Errorf("Stack should have popped")
	}

	if pop != trade {
		t.Errorf("Stack should have popped %v", pop)
	}
}
