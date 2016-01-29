package queue

import (
	"testing"
)

func TestQueue(t *testing.T) {
	q := New()
	q.Push("Hello test")

	msg := q.Pop()
	if msg != "Hello test" {
		t.Errorf("Expecting 'Hello test', but got '%s' instead", msg)
	}
}
