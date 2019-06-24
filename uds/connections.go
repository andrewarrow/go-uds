package uds

import "container/list"
import "time"
import "fmt"
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
	fmt.Println("                      empty_fromuser")
	q.fromuserm.Lock()
	defer q.fromuserm.Unlock()
	q.fromuser.Init()
}
func (q *QueueConnection) empty_txqueue() {
	fmt.Println("                      empty_touser")
	q.touserm.Lock()
	defer q.touserm.Unlock()
	q.touser.Init()
}
func (q *QueueConnection) send(payload []byte) {
	fmt.Println("                      sending to touser")
	q.touserm.Lock()
	defer q.touserm.Unlock()
	q.touser.PushBack(payload)
}
func (q *QueueConnection) other() []byte {
	for {
		q.touserm.Lock()
		if q.touser.Len() > 0 {
			e := q.touser.Front()
			q.touser.Remove(e)
			fmt.Println("                     reading from touser")
			q.touserm.Unlock()
			return e.Value.([]byte)
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
			fmt.Println("                     reading from fromuser")
			q.fromuserm.Unlock()
			return e.Value.([]byte)
		}
		q.fromuserm.Unlock()
		time.Sleep(1 * time.Millisecond)
	}
	return []byte{}
}
