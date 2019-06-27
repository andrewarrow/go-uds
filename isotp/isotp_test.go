package isotp

import "testing"
import "os"
import "time"
import "github.com/andrewarrow/go-uds/util"
import "fmt"

func TestMain(m *testing.M) {
	RXID = 0x456
	TXID = 0x123
	test_rx_queue = util.NewInterfaceQueue()
	test_tx_queue = util.NewInterfaceQueue()
	a := NewAddress(RXID, TXID)
	test_stack = NewTransport(a, stack_rxfn, stack_txfn)
	os.Exit(m.Run())
}
func TestSingleFrame(t *testing.T) {
	a := []byte{0x05, 0x11, 0x22, 0x33, 0x44, 0x55}
	msg := NewMessage(RXID, a)
	test_rx_queue.Put(msg)
	test_stack.Process()
	compareStrings(t, test_stack.Recv(), a[1:], "")
}

func TestMultiSingleFrame(t *testing.T) {
	test_stack.Process()
	test_stack.Process()
	a := []byte{0x05, 0x11, 0x22, 0x33, 0x44, 0x55}
	msg := NewMessage(RXID, a)
	test_rx_queue.Put(msg)
	test_stack.Process()
	compareStrings(t, test_stack.Recv(), a[1:], "")

	if len(test_stack.Recv()) != 0 {
		t.Fail()
	}
	test_stack.Process()
	if len(test_stack.Recv()) != 0 {
		t.Fail()
	}

	b := []byte{0x05, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE}
	msg = NewMessage(RXID, b)
	test_rx_queue.Put(msg)
	test_stack.Process()
	compareStrings(t, test_stack.Recv(), b[1:], "")
	if len(test_stack.Recv()) != 0 {
		t.Fail()
	}
	test_stack.Process()
	if len(test_stack.Recv()) != 0 {
		t.Fail()
	}
}

func TestMultipleSingleProcess(t *testing.T) {
	a := []byte{0x05, 0x11, 0x22, 0x33, 0x44, 0x55}
	b := []byte{0x05, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE}
	test_rx_queue.Put(NewMessage(RXID, a))
	test_rx_queue.Put(NewMessage(RXID, b))
	test_stack.Process()
	compareStrings(t, test_stack.Recv(), a[1:], "")
	compareStrings(t, test_stack.Recv(), b[1:], "")
	if len(test_stack.Recv()) != 0 {
		t.Fail()
	}
}

func TestMultiFrame(t *testing.T) {
	size := 10
	payload := make_payload(size, 0)
	simulate_rx(append([]byte{0x10, byte(size)}, payload[0:6]...))
	simulate_rx(append([]byte{0x21}, payload[6:10]...))

	test_stack.Process()
	compareStrings(t, test_stack.Recv(), payload, "")
	ensureEmpty(t, test_stack.Recv())
}
func TestTwoMultiFrame(t *testing.T) {
	size := 10
	payload := make_payload(size, 0)
	simulate_rx(append([]byte{0x10, byte(size)}, payload[0:6]...))
	simulate_rx(append([]byte{0x21}, payload[6:10]...))
	simulate_rx(append([]byte{0x10, byte(size)}, payload[0:6]...))
	simulate_rx(append([]byte{0x21}, payload[6:10]...))
	test_stack.Process()
	compareStrings(t, test_stack.Recv(), payload, "")
	compareStrings(t, test_stack.Recv(), payload, "")
	ensureEmpty(t, test_stack.Recv())
}
func TestMultiFrameFlowControl(t *testing.T) {
	test_stack.Stmin = 0x02
	test_stack.Blocksize = 0x05
	size := 10
	payload := make_payload(size, 0)
	simulate_rx(append([]byte{0x10, byte(size)}, payload[0:6]...))
	test_stack.Process()
	assert_sent_flow_control(t, 2, 5, 0)
	ensureEmpty(t, test_stack.Recv())
	simulate_rx(append([]byte{0x21}, payload[6:10]...))
	test_stack.Process()
	compareStrings(t, test_stack.Recv(), payload, "")
	ensureEmpty(t, test_stack.Recv())
}
func TestMultiFrameFlowControlPadding(t *testing.T) {
	test_stack.Stmin = 0x02
	test_stack.Blocksize = 0x05
	test_stack.tx_padding = 0x22
	size := 10
	payload := make_payload(size, 0)
	simulate_rx(append([]byte{0x10, byte(size)}, payload[0:6]...))
	test_stack.Process()
	assert_sent_flow_control(t, 2, 5, test_stack.tx_padding)
	ensureEmpty(t, test_stack.Recv())
	simulate_rx(append([]byte{0x21}, payload[6:10]...))
	test_stack.Process()
	compareStrings(t, test_stack.Recv(), payload, "")
	ensureEmpty(t, test_stack.Recv())
}
func TestLongMultiframeFlowControl(t *testing.T) {
	size := 30
	payload := make_payload(size, 0)
	test_stack.Stmin = 0x05
	test_stack.Blocksize = 0x3
	test_stack.tx_padding = 0
	simulate_rx(append([]byte{0x10, byte(size)}, payload[0:6]...))
	test_stack.Process()
	assert_sent_flow_control(t, 5, 3, 0)
	ensureEmpty(t, test_stack.Recv())
	simulate_rx(append([]byte{0x21}, payload[6:13]...))
	test_stack.Process()
	ensureEmpty(t, test_stack.Recv())
	simulate_rx(append([]byte{0x22}, payload[13:20]...))
	test_stack.Process()
	ensureEmpty(t, test_stack.Recv())
	simulate_rx(append([]byte{0x23}, payload[20:27]...))
	test_stack.Process()
	assert_sent_flow_control(t, 5, 3, 0)
	ensureEmpty(t, test_stack.Recv())
	simulate_rx(append([]byte{0x24}, payload[27:30]...))
	test_stack.Process()
	compareStrings(t, test_stack.Recv(), payload, "TestLongMultiframeFlowControl")
	ensureEmpty(t, test_stack.Recv())
}
func TestMultiFrameBadSeqNum(t *testing.T) {
	test_stack.Stmin = 0x02
	test_stack.Blocksize = 0x05
	size := 10
	payload := make_payload(size, 0)
	simulate_rx(append([]byte{0x10, byte(size)}, payload[0:6]...))
	simulate_rx(append([]byte{0x22}, payload[6:10]...))
	test_stack.Process()
	ensureEmpty(t, test_stack.Recv())
	_, ok := get_tx_can_msg()
	if ok {
		fmt.Println("get_tx_can_msg has msg")
		t.Fail()
	}
}

func TestTimeoutFrameAfterFirst(t *testing.T) {
	test_stack.rx_consecutive_frame_timeout = 200
	test_stack.makeTimers()
	size := 10
	payload := make_payload(size, 0)
	simulate_rx(append([]byte{0x10, byte(size)}, payload[0:6]...))
	test_stack.Process()
	time.Sleep(200 * time.Millisecond)
	simulate_rx(append([]byte{0x21}, payload[6:10]...))
	test_stack.Process()
	ensureEmpty(t, test_stack.Recv())
}
func TestRecoverTimeoutFrameAfterFirst(t *testing.T) {
	test_stack.rx_consecutive_frame_timeout = 200
	test_stack.makeTimers()
	size := 10
	payload1 := make_payload(size, 0)
	payload2 := make_payload(size, 1)
	compareNotEqStrings(t, payload1, payload2, "TestRecoverTimeoutFrameAfterFirst")
	simulate_rx(append([]byte{0x10, byte(size)}, payload1[0:6]...))
	test_stack.Process()
	time.Sleep(200 * time.Millisecond)
	simulate_rx(append([]byte{0x21}, payload1[6:10]...))
	test_stack.Process()
	ensureEmpty(t, test_stack.Recv())
	simulate_rx(append([]byte{0x10, byte(size)}, payload2[0:6]...))
	simulate_rx(append([]byte{0x21}, payload2[6:10]...))
	test_stack.Process()
	compareStrings(t, test_stack.Recv(), payload2, "TestRecoverTimeoutFrameAfterFirst")
}

func TestReceive_multiframe_interrupting_another(t *testing.T)                {}
func TestReceive_single_frame_interrupt_multiframe_then_recover(t *testing.T) {}
func TestReceive_4095_multiframe(t *testing.T)                                {}
func TestReceive_4095_multiframe_check_blocksize(t *testing.T)                {}
func TestReceive_data_length_12_bytes(t *testing.T)                           {}
func TestReceive_data_length_5_bytes(t *testing.T)                            {}
func TestReceive_data_length_12_but_set_8_bytes(t *testing.T)                 {}
func TestSend_single_frame(t *testing.T)                                      {}
func TestPadding_single_frame(t *testing.T)                                   {}
func TestPadding_single_frame_dl_12_bytes(t *testing.T)                       {}
func TestSend_multiple_single_frame_one_process(t *testing.T)                 {}
func TestSend_small_multiframe(t *testing.T)                                  {}
func TestPadding_multi_frame(t *testing.T)                                    {}
func TestPadding_multi_frame_dl_12_bytes(t *testing.T)                        {}
func TestSend_2_small_multiframe(t *testing.T)                                {}
func TestSend_multiframe_flow_control_timeout(t *testing.T)                   {}
func TestSend_multiframe_flow_control_timeout_recover(t *testing.T)           {}
func TestSend_unexpected_flow_control(t *testing.T)                           {}
func TestSend_respect_wait_frame(t *testing.T)                                {}
func TestSend_respect_wait_frame_but_timeout(t *testing.T)                    {}
func TestSend_wait_frame_after_first_frame_wftmax_0(t *testing.T)             {}
func TestSend_wait_frame_after_consecutive_frame_wftmax_0(t *testing.T)       {}
func TestSend_wait_frame_after_first_frame_reach_max(t *testing.T)            {}
func TestSend_wait_frame_after_conscutive_frame_reach_max(t *testing.T)       {}
func TestSend_4095_multiframe_zero_stmin(t *testing.T)                        {}
func TestSend_128_multiframe_variable_blocksize(t *testing.T)                 {}
func TestSquash_timing_requirement(t *testing.T)                              {}
func TestStmin_requirement(t *testing.T)                                      {}
func TestSend_nothing_with_empty_payload(t *testing.T)                        {}
func TestSend_single_frame_after_empty_payload(t *testing.T)                  {}
func TestSend_blocksize_zero(t *testing.T)                                    {}
func TestTransmit_data_length_12_bytes(t *testing.T)                          {}
func TestTransmit_data_length_5_bytes(t *testing.T)                           {}
