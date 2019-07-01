package isotp

import "testing"
import "github.com/andrewarrow/go-uds/util"
import "fmt"

var test_rx_queue *util.InterfaceQueue
var test_tx_queue *util.InterfaceQueue
var test_stack *Transport
var RXID int64
var TXID int64

func eq(t *testing.T, a, b interface{}) {
	as := fmt.Sprintf("%v", a)
	bs := fmt.Sprintf("%v", b)
	if as != bs {
		fmt.Printf("%s: %s != %s\n", t.Name(), as, bs)
		t.Fail()
	}
}
func neq(t *testing.T, a, b interface{}) {
	as := fmt.Sprintf("%v", a)
	bs := fmt.Sprintf("%v", b)
	if as == bs {
		fmt.Printf("%s: %s == %s\n", t.Name(), as, bs)
		t.Fail()
	}
}
func compareStrings(t *testing.T, a, b interface{}, msg string) {
	as := fmt.Sprintf("%v", a)
	bs := fmt.Sprintf("%v", b)
	if as != bs {
		fmt.Printf("%s: %s != %s\n", msg, as, bs)
		t.Fail()
	}
}
func compareNotEqStrings(t *testing.T, a, b interface{}, msg string) {
	as := fmt.Sprintf("%v", a)
	bs := fmt.Sprintf("%v", b)
	if as == bs {
		fmt.Printf("%s: %s == %s\n", msg, as, bs)
		t.Fail()
	}
}

func make_payload(size, start_val int) []byte {
	a := []byte{}
	i := start_val
	for {
		a = append(a, byte(i%0x100))
		i++
		if i >= start_val+size {
			break
		}
	}
	return a
}

func simulate_rx(b []byte) {
	test_rx_queue.Put(NewMessage(RXID, b))
}
func ensureEmpty(t *testing.T, b []byte) {
	if len(b) != 0 {
		t.Logf("%v", b)
		t.Fail()
	}
}
func assert_sent_flow_control(t *testing.T, stmin, blocksize, tx_padding int) {
	msg, ok := get_tx_can_msg()
	if ok == false {
		fmt.Println("get_tx_can_msg has no msg")
		t.Fail()
	}
	data := []byte{}
	data = append(data, 0x30, byte(blocksize), byte(stmin))
	if tx_padding > 0 {
		for {
			data = append(data, byte(tx_padding))
			if len(data) == 8 {
				break
			}
		}
	}

	compareStrings(t, msg.Payload, data, "hi")
}

func stack_rxfn() (Message, bool) {
	if test_rx_queue.Len() > 0 {
		e := test_rx_queue.Get()
		return e.(Message), true
	}
	return Message{}, false
}
func get_tx_can_msg() (Message, bool) {
	if test_tx_queue.Len() > 0 {
		e := test_tx_queue.Get()
		return e.(Message), true
	}
	return Message{}, false
}
func stack_txfn(m Message) {
	test_tx_queue.Put(m)
}
