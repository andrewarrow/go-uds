package uds

import "time"

//import "fmt"
import "github.com/andrewarrow/go-uds/util"
import "sync"

type QueueConnection struct {
	name      string
	mtu       int
	fromuser  *util.Queue
	fromuserm sync.Mutex
	touser    *util.Queue
	touserm   sync.Mutex
}

func NewQueueConnection(name string, mtu int) *QueueConnection {
	q := QueueConnection{}
	q.name = name
	q.mtu = mtu
	q.fromuser = util.NewQueue()
	q.touser = util.NewQueue()
	return &q
}

func (q *QueueConnection) Send_and_wait_for_reply(payload []byte) []byte {
	return []byte{}
}
func (q *QueueConnection) Send_and_grant_flow_request(payload []byte, length int) []byte {
	return []byte{}
}
func (q *QueueConnection) Empty_rxqueue() {
	q.fromuserm.Lock()
	defer q.fromuserm.Unlock()
	q.fromuser.Clear()
}
func (q *QueueConnection) Empty_txqueue() {
	q.touserm.Lock()
	defer q.touserm.Unlock()
	q.touser.Clear()
}
func (q *QueueConnection) Send(payload []byte) {
	//fmt.Printf("                      sending to touser %v\n", payload)
	q.touserm.Lock()
	defer q.touserm.Unlock()
	q.touser.Put(payload)
}
func (q *QueueConnection) touser_frame() []byte {
	for {
		q.touserm.Lock()
		if q.touser.Len() > 0 {
			val := q.touser.Get()
			q.touserm.Unlock()
			//fmt.Printf("                     reading from touser %v\n", val)
			return val
		}
		q.touserm.Unlock()
		time.Sleep(1 * time.Millisecond)
	}
	return []byte{}
}
func (q *QueueConnection) Wait_frame() []byte {
	for {
		q.fromuserm.Lock()
		if q.fromuser.Len() > 0 {
			val := q.fromuser.Get()
			q.fromuserm.Unlock()
			//fmt.Printf("                     reading from fromuser %v\n", val)
			return val
		}
		q.fromuserm.Unlock()
		time.Sleep(1 * time.Millisecond)
	}
	return []byte{}
}
