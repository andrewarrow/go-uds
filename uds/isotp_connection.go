package uds

import "time"
import "github.com/andrewarrow/go-isotp/isotp"

type IsotpConnection struct {
	name  string
	mtu   int
	stack *isotp.Transport
	rxfn  func() (isotp.Message, bool)
	txfn  func(msg isotp.Message)
}

func NewIsotpConnection(rx, tx int, rxfn func() (isotp.Message, bool),
	txfn func(msg isotp.Message)) *IsotpConnection {
	ic := IsotpConnection{}
	a := isotp.NewAddress(rx, tx)
	ic.stack = isotp.NewTransport(a, rxfn, txfn)
	ic.rxfn = rxfn
	ic.txfn = txfn
	return &ic
}

func (ic *IsotpConnection) empty_rxqueue() {
}
func (ic *IsotpConnection) empty_txqueue() {
}
func (ic *IsotpConnection) send(payload []byte) {
	//todo
	msg := isotp.Message{}
	ic.txfn(msg)
}
func (ic *IsotpConnection) wait_frame() []byte {
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
