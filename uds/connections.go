package uds

import "container/list"
import "time"
import "fmt"

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
	fmt.Println("                      empty_fromuser")
	q.fromuser.Init()
}
func (q *QueueConnection) empty_txqueue() {
	fmt.Println("                      empty_touser")
	q.touser.Init()
}
func (q *QueueConnection) send(payload []byte) {
	fmt.Println("                      sending to touser")
	q.touser.PushBack(payload)
}
func (q *QueueConnection) other() []byte {
	for {
		if q.touser.Len() > 0 {
			e := q.touser.Front()
			q.touser.Remove(e)
			fmt.Println("                     reading from touser")
			return e.Value.([]byte)
		}
		time.Sleep(1 * time.Millisecond)
	}
	return []byte{}
}
func (q *QueueConnection) wait_frame() []byte {
	for {
		if q.fromuser.Len() > 0 {
			e := q.fromuser.Front()
			q.fromuser.Remove(e)
			fmt.Println("                     reading from fromuser")
			return e.Value.([]byte)
		}
		time.Sleep(1 * time.Millisecond)
	}
	return []byte{}
}
