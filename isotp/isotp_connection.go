package isotp

import "time"

import "fmt"
import "github.com/andrewarrow/go-uds/util"

type AnyConn interface {
	Empty_rxqueue()
	Empty_txqueue()
	Send(payload []byte)
	Wait_frame() []byte
	Send_and_grant_flow_request(payload []byte, length int) []byte
	Send_and_wait_for_reply(payload []byte) []byte
}

type IsotpConnection struct {
	name           string
	mtu            int
	fromIsoTPQueue *util.Queue
	toIsoTPQueue   *util.Queue
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
	ic.fromIsoTPQueue = util.NewQueue()
	ic.toIsoTPQueue = util.NewQueue()
	return &ic
}

func (ic *IsotpConnection) asSingleFrameOrMulti(payload []byte) []Message {
	ms := []Message{}
	if len(payload) <= 8 {
		fmt.Println("less than 8")
		msg_data := append([]byte{byte(0x0 | len(payload))}, payload...)
		msg := ic.Stack.make_tx_msg(ic.Stack.address.txid, msg_data)
		//ic.Stack.txfn(msg)
		ms = append(ms, msg)
	} else {
		fmt.Println("more than 8")
		chunkSize := 6
		seq := 0
		for i := 0; i < len(payload); i += chunkSize {
			end := i + chunkSize

			if end > len(payload) {
				end = len(payload)
			}

			chunk := payload[i:end]
			tx_frame_length := len(payload)
			msg_data := []byte{}
			if seq == 0 {
				msg_data = append([]byte{0x10 | byte((tx_frame_length>>8)&0xF), byte(tx_frame_length & 0xFF)}, chunk...)
			} else {
				msg_data = append([]byte{0x20 | byte(seq)}, chunk...)
			}
			msg := ic.Stack.make_tx_msg(ic.Stack.address.txid, msg_data)
			fmt.Println("sending ", msg)
			//	ic.Stack.txfn(msg)
			ms = append(ms, msg)
			seq += 1
		}
	}
	return ms
}

func (ic *IsotpConnection) Send_and_wait_for_reply(payload []byte) []byte {

	msgs := ic.asSingleFrameOrMulti(payload)
	flow := []byte{}

	t1 := time.Now().Unix()
	for _, msg := range msgs {
		flow = []byte{}
		ic.Stack.txfn(msg)
		for {
			if time.Now().Unix()-t1 > 15 {
				fmt.Println("timeout1")
				break
			}
			msg, _ := ic.Stack.rxfn()
			if ic.Stack.address.is_for_me(msg) == false {
				continue
			}
			fmt.Println("got reply", msg.Payload)
			flow = append(flow, msg.Payload...)
			if true {
				break
			}
		}
	}
	return flow
}
func (ic *IsotpConnection) Send_and_grant_flow_request(payload []byte, length int) []byte {
	msg_data := append([]byte{byte(0x0 | len(payload))}, payload...)
	msg := ic.Stack.make_tx_msg(ic.Stack.address.txid, msg_data)
	ic.Stack.txfn(msg)
	flow := []byte{}
	// wait for flow request

	t1 := time.Now().Unix()
	for {
		if time.Now().Unix()-t1 > 5 {
			fmt.Println("timeout1")
			break
		}
		msg, _ := ic.Stack.rxfn()
		if ic.Stack.address.is_for_me(msg) == false {
			continue
		}
		flow = append(flow, msg.Payload...)
		if length < 7 {
			flow = append([]byte{0}, flow...)
			return flow
		}
		if true {
			break
		}
	}
	msg = ic.Stack.make_flow_control(CONTINUE)
	ic.Stack.txfn(msg)
	// read flow
	t1 = time.Now().Unix()
	for {
		if time.Now().Unix()-t1 > 5 {
			fmt.Println("timeout2")
			break
		}
		msg, _ := ic.Stack.rxfn()
		if ic.Stack.address.is_for_me(msg) == false {
			continue
		}
		flow = append(flow, msg.Payload[1:]...)
		if len(flow) > length {
			break
		}
	}

	return flow
}

func (ic *IsotpConnection) Open() {
	go func() {
		for {
			//fmt.Println("  [ml] toIsoTP")
			for {
				if ic.toIsoTPQueue.Len() == 0 {
					break
				}
				payload := ic.toIsoTPQueue.Get()
				ic.Stack.Send(payload)
			}

			ic.Stack.Process()

			//fmt.Println("  [ml] fromIsoTP")
			for {
				if ic.Stack.available() == false {
					break
				}
				stuff := ic.Stack.Recv()
				ic.fromIsoTPQueue.Put(stuff)
			}
			//fmt.Println("  [ml] sleep")
			time.Sleep(time.Millisecond * 10)
		}
	}()
}

func (ic *IsotpConnection) Empty_rxqueue() {
	ic.fromIsoTPQueue.Clear()
}
func (ic *IsotpConnection) Empty_txqueue() {
	ic.toIsoTPQueue.Clear()
}
func (ic *IsotpConnection) Send(payload []byte) {
	//msg := NewMessage(ic.Stack.address.rxid, payload)
	ic.toIsoTPQueue.Put(payload)
}
func (ic *IsotpConnection) Wait_frame() []byte {

	count := 0
	for {
		if ic.fromIsoTPQueue.Len() > 0 {
			m := ic.fromIsoTPQueue.Get()
			return m
		}
		time.Sleep(500 * time.Millisecond)
		count++
		if count > 30 {
			break
		}
	}
	return []byte{}
}
