package queue

import (
	"sync"
)

type Queue struct {
	sync.Mutex
	data *[]string
}

func (q *Queue) Len() int {
	return len(*q.data)
}

func (q *Queue) Push(s string) {
	*q.data = append(*q.data, s)
}

func (q *Queue) Pop() string {
	n := (*q.data)[0]
	*q.data = (*q.data)[1:]
	return n
}

func New() *Queue {
	q := &Queue{
		data: new([]string),
	}
	return q
}
