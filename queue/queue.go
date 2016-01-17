package queue

type Queue []string

func (q *Queue) Len() int {
	return len(*q)
}

func (q *Queue) Push(s string) {
	*q = append(*q, s)
}

func (q *Queue) Pop() string {
	n := (*q)[0]
	*q = (*q)[1:]
	return n
}

func New() *Queue {
	q := new(Queue)
	return q
}
