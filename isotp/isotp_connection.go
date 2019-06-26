package isotp

import "time"

import "fmt"
import "container/list"

type AnyConn interface {
	Empty_rxqueue()
	Empty_txqueue()
	Send(payload []byte)
	Wait_frame() []byte
}

type IsotpConnection struct {
	name     string
	mtu      int
	rx_queue *list.List
	Stack    *Transport
	rxfn     func() (Message, bool)
	txfn     func(msg Message)
}

func NewIsotpConnection(rx, tx int64, rxfn func() (Message, bool),
	txfn func(msg Message)) *IsotpConnection {
	ic := IsotpConnection{}
	a := NewAddress(rx, tx)
	ic.Stack = NewTransport(a, rxfn, txfn)
	ic.rxfn = rxfn
	ic.txfn = txfn
	ic.rx_queue = list.New()
	return &ic
}

func (ic *IsotpConnection) Open() {
	go func() {
		//rxthread
		for {
			msg, ok := ic.rxfn()
			if ok {
				fmt.Println(msg.ToBytes())
				ic.rx_queue.PushBack(msg.ToBytes())
			}
			time.Sleep(1 * time.Millisecond)
		}
	}()
}

func (ic *IsotpConnection) Empty_rxqueue() {
	ic.rx_queue.Init()
}
func (ic *IsotpConnection) Empty_txqueue() {
}
func (ic *IsotpConnection) Send(payload []byte) {
	msg := NewMessage(ic.Stack.address.rxid, payload)
	ic.txfn(msg)
}
func (ic *IsotpConnection) Wait_frame() []byte {
	count := 0
	for {
		if ic.rx_queue.Len() > 0 {
			e := ic.rx_queue.Front()
			ic.rx_queue.Remove(e)
			m := e.Value.([]byte)
			return m
		}
		time.Sleep(1 * time.Millisecond)
		count++
		if count > 200 {
			break
		}
	}
	return []byte{}
}
