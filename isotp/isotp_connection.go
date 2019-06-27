package isotp

import "time"

//import "fmt"
import "container/list"

type AnyConn interface {
	Empty_rxqueue()
	Empty_txqueue()
	Send(payload []byte)
	Wait_frame() []byte
}

type IsotpConnection struct {
	name           string
	mtu            int
	fromIsoTPQueue *list.List
	toIsoTPQueue   *list.List
	Stack          *Transport
	rxfn           func() (Message, bool)
	txfn           func(msg Message)
}

func NewIsotpConnection(rx, tx int64, rxfn func() (Message, bool),
	txfn func(msg Message)) *IsotpConnection {
	ic := IsotpConnection{}
	a := NewAddress(rx, tx)
	ic.Stack = NewTransport(a, rxfn, txfn)
	ic.rxfn = rxfn
	ic.txfn = txfn
	ic.fromIsoTPQueue = list.New()
	ic.toIsoTPQueue = list.New()
	return &ic
}

func (ic *IsotpConnection) Open() {
	go func() {
		for {
			for {
				if ic.toIsoTPQueue.Len() == 0 {
					break
				}
				e := ic.toIsoTPQueue.Front()
				ic.toIsoTPQueue.Remove(e)
				payload := e.Value.([]byte)
				ic.Stack.Send(payload)
			}

			ic.Stack.Process()

			for {
				if ic.Stack.available() == false {
					break
				}
				ic.fromIsoTPQueue.PushBack(ic.Stack.Recv())
			}
		}
	}()
}

/*
   time.sleep(self.isotp_layer.sleep_time())

*/

func (ic *IsotpConnection) Empty_rxqueue() {
	ic.fromIsoTPQueue.Init()
}
func (ic *IsotpConnection) Empty_txqueue() {
	ic.toIsoTPQueue.Init()
}
func (ic *IsotpConnection) Send(payload []byte) {
	//msg := NewMessage(ic.Stack.address.rxid, payload)
	ic.toIsoTPQueue.PushBack(payload)
}
func (ic *IsotpConnection) Wait_frame() []byte {

	count := 0
	for {
		if ic.fromIsoTPQueue.Len() > 0 {
			e := ic.fromIsoTPQueue.Front()
			ic.fromIsoTPQueue.Remove(e)
			m := e.Value.([]byte)
			return m
		}
		time.Sleep(1 * time.Millisecond)
		count++
		if count > 20 {
			break
		}
	}
	return []byte{}
}
