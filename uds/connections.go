package uds

import "container/list"

type AnyConn interface {
	empty_rxqueue()
	empty_txqueue()
	send(payload []byte)
	wait_frame() []byte
}

type QueueConnection struct {
	name     string
	mtu      int
	fromuser *list.List
	touser   *list.List
}

func NewQueueConnection(name string, mtu int) *QueueConnection {
	q := QueueConnection{}
	q.name = name
	q.mtu = mtu
	q.fromuser = list.New()
	q.touser = list.New()
	return &q
}

func (q *QueueConnection) empty_rxqueue() {
	q.fromuser.Init()
}
func (q *QueueConnection) empty_txqueue() {
	q.touser.Init()
}
func (q *QueueConnection) send(payload []byte) {
	q.touser.PushBack(payload)
}
func (q *QueueConnection) wait_frame() []byte {
	if q.fromuser.Len() > 0 {
		e := q.fromuser.Front()
		q.fromuser.Remove(e)
		return e.Value.([]byte)
	}
	return []byte{}
}
