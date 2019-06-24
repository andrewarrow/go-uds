package uds

import "time"
import "github.com/andrewarrow/go-isotp/isotp"

type IsotpConnection struct {
	name  string
	mtu   int
	stack *isotp.Transport
}

func NewIsotpConnection(rx, tx int, rxfn func() (isotp.Message, bool),
	txfn func(msg isotp.Message)) *IsotpConnection {
	ic := IsotpConnection{}
	a := isotp.NewAddress(rx, tx)
	ic.stack = isotp.NewTransport(a, rxfn, txfn)
	return &ic
}

func (ic *IsotpConnection) empty_rxqueue() {
}
func (ic *IsotpConnection) empty_txqueue() {
}
func (ic *IsotpConnection) send(payload []byte) {
	//C.hardware_send()
}
func (ic *IsotpConnection) touser_frame() []byte {
	for {
		time.Sleep(1 * time.Millisecond)
	}
	return []byte{}
}
func (ic *IsotpConnection) wait_frame() []byte {
	for {
		time.Sleep(1 * time.Millisecond)
	}
	return []byte{}
}
