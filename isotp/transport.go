package isotp

import "github.com/andrewarrow/go-uds/util"

//import "fmt"

const FLOW = "FLOW"
const WAIT = "WAIT"
const CONTINUE = "CONTINUE"
const OVERFLOW = "OVERFLOW"
const IDLE = "IDLE"
const TRANSMIT = "TRANSMIT"
const SINGLE = "SINGLE"
const FIRST = "FIRST"
const CONSECUTIVE = "CONSECUTIVE"

type Transport struct {
	rxfn                         func() (Message, bool)
	txfn                         func(m Message)
	rx_queue                     *util.Queue
	tx_queue                     *util.Queue
	address                      Address
	rx_state                     string
	tx_state                     string
	last_seqnum                  int
	rx_block_counter             int
	rx_frame_length              int
	tx_frame_length              int
	rx_buffer                    []byte
	pending_flow_control_tx      bool
	last_flow_control_frame      *PDU
	Stmin                        int
	Blocksize                    int
	squash_stmin_requirement     bool
	rx_flowcontrol_timeout       int
	rx_consecutive_frame_timeout int
	wftmax                       int
	data_length                  int
	timer_rx_cf                  *Timer
	timer_rx_fc                  *Timer
	timer_tx_stmin               *Timer
	tx_padding                   int
	tx_buffer                    []byte
	remote_blocksize             int
	tx_block_counter             int
	tx_seqnum                    int
	wft_counter                  int
}

func NewTransport(a Address, rxfn func() (Message, bool), txfn func(m Message)) *Transport {
	t := Transport{}
	t.rxfn = rxfn
	t.txfn = txfn
	t.rx_queue = util.NewQueue()
	t.tx_queue = util.NewQueue()
	t.address = a
	t.rx_state = IDLE
	t.tx_state = IDLE
	t.rx_buffer = []byte{}
	t.rx_frame_length = 0
	t.Stmin = 1
	t.Blocksize = 8
	t.data_length = 8
	t.rx_flowcontrol_timeout = 1000
	t.rx_consecutive_frame_timeout = 1000
	t.makeTimers()
	return &t
}

func (t *Transport) makeTimers() {
	t.timer_rx_fc = NewTimer(float32(t.rx_flowcontrol_timeout) / 1000.0)
	t.timer_rx_cf = NewTimer(float32(t.rx_consecutive_frame_timeout) / 1000.0)
}

func (t *Transport) Process() {
	i := 0
	for {
		msg, ok := t.rxfn()
		if ok == false {
			break
		}
		t.process_rx(msg)
		i++
		if i > 50 {
			break
		}
	}

	for {
		msg, ok := t.process_tx()
		if ok == false {
			break
		}
		t.txfn(msg)
	}
}

func (t *Transport) start_reception_after_first_frame(frame PDU) {
	t.last_seqnum = 0
	t.rx_block_counter = 0
	t.rx_frame_length = frame.length
	t.rx_state = WAIT
	t.rx_buffer = append([]byte{}, frame.payload...)
	t.pending_flow_control_tx = true
	//fmt.Println("start_reception_after_first_frame", t.timer_rx_cf.startedAt)
	t.timer_rx_cf.start()
	//fmt.Println("start_reception_after_first_frame", t.timer_rx_cf.startedAt)
}

func (t *Transport) Send(data []byte) {
	t.tx_queue.Put(data)
	//self.tx_queue.put( {'data':data, 'target_address_type':target_address_type})    # frame is always an IsoTPFrame here
}
func (t *Transport) Recv() []byte {
	if t.rx_queue.Len() == 0 {
		return []byte{}
	}
	return t.rx_queue.Get()
}

func (t *Transport) stop_receiving() {
	t.rx_state = IDLE
	t.rx_buffer = []byte{}
	t.pending_flow_control_tx = false
	t.last_flow_control_frame = nil
	t.timer_rx_cf.stop()
}
func (t *Transport) stop_sending() {
	t.tx_state = IDLE
	t.tx_buffer = []byte{}
	t.tx_frame_length = 0
	t.timer_rx_fc.stop()
	t.timer_tx_stmin.stop()
	t.remote_blocksize = 0
	t.tx_block_counter = 0
	t.tx_seqnum = 0
	t.wft_counter = 0
}

func (t *Transport) available() bool {
	return t.rx_queue.Len() > 0
}
