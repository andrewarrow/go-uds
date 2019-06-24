package isotp

import "time"

type AnyConn interface {
	Empty_rxqueue()
	Empty_txqueue()
	Send(payload []byte)
	Wait_frame() []byte
}

type IsotpConnection struct {
	name  string
	mtu   int
	stack *Transport
	rxfn  func() (Message, bool)
	txfn  func(msg Message)
}

func NewIsotpConnection(rx, tx int, rxfn func() (Message, bool),
	txfn func(msg Message)) *IsotpConnection {
	ic := IsotpConnection{}
	a := NewAddress(rx, tx)
	ic.stack = NewTransport(a, rxfn, txfn)
	ic.rxfn = rxfn
	ic.txfn = txfn
	return &ic
}

func (ic *IsotpConnection) Empty_rxqueue() {
}
func (ic *IsotpConnection) Empty_txqueue() {
}
func (ic *IsotpConnection) Send(payload []byte) {
	//todo
	msg := NewMessage(0, payload)
	ic.txfn(msg)
}
func (ic *IsotpConnection) Wait_frame() []byte {
	count := 0
	for {
		msg, ok := ic.rxfn()
		if ok {
			return msg.GetData()
		}
		time.Sleep(1 * time.Millisecond)
		count++
		if count > 200 {
			break
		}
	}
	return []byte{}
}
