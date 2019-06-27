package util

import "container/list"

type Queue struct {
	data *list.List
}

func NewQueue() *Queue {
	q := Queue{}
	q.data = list.New()
	return &q
}

func (q *Queue) Empty() bool {
	return q.data.Len() == 0
}
func (q *Queue) Clear() {
	q.data.Init()
}

func (q *Queue) Len() int {
	return q.data.Len()
}
func (q *Queue) Available() bool {
	return q.data.Len() > 0
}
func (q *Queue) Get() []byte {
	e := q.data.Front()
	q.data.Remove(e)
	return e.Value.([]byte)
}
func (q *Queue) Put(b []byte) {
	q.data.PushBack(b)
}
