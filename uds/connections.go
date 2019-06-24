package uds

import "container/list"
import "time"

//import "fmt"
import "sync"

type AnyConn interface {
	empty_rxqueue()
	empty_txqueue()
	send(payload []byte)
	wait_frame() []byte
}

type QueueConnection struct {
	name      string
	mtu       int
	fromuser  *list.List
	fromuserm sync.Mutex
	touser    *list.List
	touserm   sync.Mutex
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
	q.fromuserm.Lock()
	defer q.fromuserm.Unlock()
	q.fromuser.Init()
}
func (q *QueueConnection) empty_txqueue() {
	q.touserm.Lock()
	defer q.touserm.Unlock()
	q.touser.Init()
}
func (q *QueueConnection) send(payload []byte) {
	//fmt.Printf("                      sending to touser %v\n", payload)
	q.touserm.Lock()
	defer q.touserm.Unlock()
	q.touser.PushBack(payload)
}
func (q *QueueConnection) touser_frame() []byte {
	for {
		q.touserm.Lock()
		if q.touser.Len() > 0 {
			e := q.touser.Front()
			q.touser.Remove(e)
			q.touserm.Unlock()
			val := e.Value.([]byte)
			//fmt.Printf("                     reading from touser %v\n", val)
			return val
		}
		q.touserm.Unlock()
		time.Sleep(1 * time.Millisecond)
	}
	return []byte{}
}
func (q *QueueConnection) wait_frame() []byte {
	for {
		q.fromuserm.Lock()
		if q.fromuser.Len() > 0 {
			e := q.fromuser.Front()
			q.fromuser.Remove(e)
			q.fromuserm.Unlock()
			val := e.Value.([]byte)
			//fmt.Printf("                     reading from fromuser %v\n", val)
			return val
		}
		q.fromuserm.Unlock()
		time.Sleep(1 * time.Millisecond)
	}
	return []byte{}
}
