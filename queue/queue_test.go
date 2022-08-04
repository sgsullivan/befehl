package queue

import (
	"testing"
)

func TestNew(t *testing.T) {
	q := new(Queue).New(42)
	if q.count != 42 {
		t.Fatalf("count value of %d is unexpected", q.count)
	}
}

func TestDecrementCounter(t *testing.T) {
	q := new(Queue).New(42)
	newCount := q.DecrementCounter()
	if q.count != 41 {
		t.Fatalf("count value of %d is unexpected", q.count)
	}
	if newCount != 41 {
		t.Fatalf("newCount of %d is unexpected", newCount)
	}
}
