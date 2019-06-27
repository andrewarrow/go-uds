package isotp

//import "fmt"

func (t *Transport) process_tx() (Message, bool) {
	m := Message{}
	//output_msg := false

	//fmt.Println("calling process_tx ", t.pending_flow_control_tx)
	if t.pending_flow_control_tx {
		t.pending_flow_control_tx = false
		return t.make_flow_control(CONTINUE), true
	}
	if t.tx_state == IDLE {
		read_tx_queue := true
		for {
			read_tx_queue = false
			if t.tx_queue.Len() > 0 {
				payload := t.tx_queue.Get()
				if len(payload) == 0 {
					read_tx_queue = true
				} else {
					t.tx_buffer = payload
				}
			}
			if read_tx_queue == false {
				break
			}
		}
		//msg_data  = self.address.tx_payload_prefix + bytearray([0x0 | len(self.tx_buffer)]) + self.tx_buffer
		msg_data := append([]byte{byte(0x0 | len(t.tx_buffer))}, t.tx_buffer...)
		m = t.make_tx_msg(t.address.txid, msg_data)
		return m, true
	}

	/*
		if len(self.tx_buffer) <= self.params.ll_data_length-1-len(self.address.tx_payload_prefix):
		  msg_data  = self.address.tx_payload_prefix + bytearray([0x0 | len(self.tx_buffer)]) + self.tx_buffer
		  print("here1", msg_data)
		  arbitration_id  = self.address.get_tx_arbitraton_id(popped_object['target_address_type'])
		  output_msg              = self.make_tx_msg(arbitration_id, msg_data)

	*/
	return m, false
}

func (t *Transport) make_flow_control(status string) Message {
	payload := craft_flow_control_data(status, t.Blocksize, t.Stmin)
	return t.make_tx_msg(t.address._get_tx_arbitraton_id(Physical), append(t.address.tx_payload_prefix, payload...))
}

func craft_flow_control_data(status string, blocksize int, stmin int) []byte {
	f := 0
	if status == WAIT {
		f = 1
	} else if status == OVERFLOW {
		f = 2
	}

	return []byte{byte(0x30 | (f)&0xF), byte(blocksize & 0xFF), byte(stmin & 0xFF)}
}

func (t *Transport) pad_message_data(payload []byte) []byte {
	if len(payload) < t.data_length && t.tx_padding > 0 {
		a := []byte{}
		for {
			a = append(a, byte(t.tx_padding&0xFF))
			if len(a) == t.data_length-len(payload) {
				break
			}
		}
		return append(payload, a...)
	}
	return payload
}

func (t *Transport) make_tx_msg(arbitration_id int64, payload []byte) Message {
	data := t.pad_message_data(payload)
	m := NewMessage(arbitration_id, data)
	m.extended_id = t.address.is_29bits()
	return m
}
