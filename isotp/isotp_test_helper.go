package isotp

import "testing"
import "container/list"
import "fmt"

var test_rx_queue *list.List
var test_tx_queue *list.List
var test_stack *Transport
var RXID int
var TXID int

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
		a = append(a, byte(i))
		i++
		if i >= start_val+size {
			break
		}
	}
	return a
}

func simulate_rx(b []byte) {
	test_rx_queue.PushBack(NewMessage(RXID, b))
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
		e := test_rx_queue.Front()
		test_rx_queue.Remove(e)
		return e.Value.(Message), true
	}
	return Message{}, false
}
func get_tx_can_msg() (Message, bool) {
	if test_tx_queue.Len() > 0 {
		e := test_tx_queue.Front()
		test_tx_queue.Remove(e)
		return e.Value.(Message), true
	}
	return Message{}, false
}
func stack_txfn(m Message) {
	test_tx_queue.PushBack(m)
}
