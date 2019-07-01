package util

import "container/list"

import "fmt"

type InterfaceQueue struct {
	data *list.List
}

func NewInterfaceQueue() *InterfaceQueue {
	q := InterfaceQueue{}
	q.data = list.New()
	return &q
}

func (q *InterfaceQueue) Empty() bool {
	return q.data.Len() == 0
}
func (q *InterfaceQueue) Clear() {
	q.data.Init()
}

func (q *InterfaceQueue) Len() int {
	return q.data.Len()
}
func (q *InterfaceQueue) Available() bool {
	return q.data.Len() > 0
}
func (q *InterfaceQueue) Get() interface{} {
	e := q.data.Front()
	q.data.Remove(e)
	return e.Value
}
func (q *InterfaceQueue) Put(b interface{}) {
	q.data.PushBack(b)
}
func (q *InterfaceQueue) String() string {
	return fmt.Sprintf("%v", q.data.Front())
}
