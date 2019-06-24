package isotp

import "fmt"

func (t *Transport) process_rx(msg Message) {
	if t.address.is_for_me(msg) == false {
		return
	}
	pdu := NewPDU(msg, t.address.rx_prefix_size, t.data_length)
	fmt.Println("calling process_rx ", pdu)
	if t.timer_rx_cf.is_timed_out() {
		fmt.Println("Reception of CONSECUTIVE_FRAME timed out.")
		t.stop_receiving()
	}
	if pdu.flavor == FLOW {
		t.last_flow_control_frame = &pdu
		if t.rx_state == WAIT {
			if pdu.flow == WAIT || pdu.flow == CONTINUE {
				t.timer_rx_cf.start()
			}
		}
		return
	}

	if t.rx_state == IDLE {
		t.rx_frame_length = 0
		t.timer_rx_cf.stop()
		if pdu.flavor == SINGLE && pdu.length > 0 {
			t.rx_queue.PushBack(pdu.payload)
		} else if pdu.flavor == FIRST {
			t.start_reception_after_first_frame(pdu)
		} else if pdu.flavor == CONSECUTIVE {
			fmt.Println("Received a ConsecutiveFrame while reception was idle.")
		}
	} else if t.rx_state == WAIT {
		if pdu.flavor == SINGLE && pdu.length > 0 {
			t.rx_queue.PushBack(pdu.payload)
			t.rx_state = IDLE
			fmt.Println("Reception of IsoTP frame interrupted with a new SingleFrame")
		} else if pdu.flavor == FIRST {
			t.start_reception_after_first_frame(pdu)
			fmt.Println("Reception of IsoTP frame interrupted with a new FirstFrame")
		} else if pdu.flavor == CONSECUTIVE {
			t.timer_rx_cf.start()
			expected_seqnum := (t.last_seqnum + 1) & 0xF
			if pdu.seqnum == expected_seqnum {
				t.last_seqnum = pdu.seqnum

				bytes_to_receive := (t.rx_frame_length - len(t.rx_buffer))

				if len(pdu.payload) > bytes_to_receive {
					t.rx_buffer = append(t.rx_buffer, pdu.payload[0:bytes_to_receive]...)
				} else {
					t.rx_buffer = append(t.rx_buffer, pdu.payload...)
				}
				if len(t.rx_buffer) >= t.rx_frame_length {
					t.rx_queue.PushBack(append([]byte{}, t.rx_buffer...))
					t.stop_receiving()
				} else {
					t.rx_block_counter += 1
					if t.rx_block_counter%t.blocksize == 0 {
						t.pending_flow_control_tx = true
						t.timer_rx_cf.stop()
					}
				}
			} else {
				t.stop_receiving()
				fmt.Println("Received a ConsecutiveFrame with wrong SequenceNumber")
			}
		}
	}
	//fmt.Println(pdu)
}
