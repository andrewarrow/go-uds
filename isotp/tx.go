package isotp

import "fmt"

/*

	# ======= Timeouts ======
	if self.timer_rx_fc.is_timed_out():
		self.trigger_error(isotp.errors.FlowControlTimeoutError('Reception of FlowControl timed out. Stopping transmission'))
		self.stop_sending()

	# ======= FSM ======

	# Check this first as we may have another isotp frame to send and we need to handle it right away without waiting for next "process()" call
	if self.tx_state != self.TxState.IDLE and len(self.tx_buffer) == 0:
		self.stop_sending()

	if self.tx_state == self.TxState.IDLE:
		read_tx_queue = True	# Read until we get non-empty frame to send
		while read_tx_queue:
			read_tx_queue = False
			if not self.tx_queue.empty():
				popped_object = self.tx_queue.get()
				if len(popped_object['data']) == 0:
					read_tx_queue = True	# Read another frame from tx_queue
				else:
					self.tx_buffer = bytearray(popped_object['data'])
					size_on_first_byte = True if len(self.tx_buffer) <= 7 else False
					size_offset = 1 if size_on_first_byte else 2
					if len(self.tx_buffer) <= self.params.ll_data_length-size_offset-len(self.address.tx_payload_prefix):	# Single frame
						if size_on_first_byte:
							msg_data 	= self.address.tx_payload_prefix + bytearray([0x0 | len(self.tx_buffer)]) + self.tx_buffer
						else:
							msg_data 	= self.address.tx_payload_prefix + bytearray([0x0, len(self.tx_buffer)]) + self.tx_buffer
						arbitration_id 	= self.address.get_tx_arbitraton_id(popped_object['target_address_type'])
						output_msg		= self.make_tx_msg(arbitration_id, msg_data)
					else:							# Multi frame
						self.tx_frame_length = len(self.tx_buffer)
						encode_length_on_2_first_bytes = True if self.tx_frame_length <= 4095 else False

						if encode_length_on_2_first_bytes:
							data_length = self.params.ll_data_length-2-len(self.address.tx_payload_prefix)
							msg_data 	= self.address.tx_payload_prefix + bytearray([0x10|((self.tx_frame_length >> 8) & 0xF), self.tx_frame_length&0xFF]) + self.tx_buffer[:data_length]
						else:
							data_length = self.params.ll_data_length-6-len(self.address.tx_payload_prefix)
							msg_data 	= self.address.tx_payload_prefix + bytearray([0x10, 0x00, (self.tx_frame_length>>24) & 0xFF, (self.tx_frame_length>>16) & 0xFF, (self.tx_frame_length>>8) & 0xFF, (self.tx_frame_length>>0) & 0xFF]) + self.tx_buffer[:data_length]

						arbitration_id 	= self.address.get_tx_arbitraton_id()
						output_msg 		= self.make_tx_msg(arbitration_id, msg_data)
						self.tx_buffer 	= self.tx_buffer[data_length:]
						self.tx_state 	= self.TxState.WAIT_FC
						self.tx_seqnum 	= 1
						self.start_rx_fc_timer()


	elif self.tx_state == self.TxState.WAIT_FC:
		pass # Nothing to do. Flow control will make the FSM switch state by calling init_tx_consecutive_frame

	elif self.tx_state == self.TxState.TRANSMIT_CF:
		if self.timer_tx_stmin.is_timed_out() or self.params.squash_stmin_requirement:
			data_length = self.params.ll_data_length-1-len(self.address.tx_payload_prefix)
			msg_data = self.address.tx_payload_prefix + bytearray([0x20 | self.tx_seqnum]) + self.tx_buffer[:data_length]
			arbitration_id 	= self.address.get_tx_arbitraton_id()
			output_msg = self.make_tx_msg(arbitration_id, msg_data)
			self.tx_buffer = self.tx_buffer[data_length:]
			self.tx_seqnum = (self.tx_seqnum + 1 ) & 0xF
			self.timer_tx_stmin.start()
			self.tx_block_counter+=1

		if self.remote_blocksize != 0 and self.tx_block_counter >= self.remote_blocksize:
			self.tx_state = self.TxState.WAIT_FC
			self.start_rx_fc_timer()

	return output_msg
*/

func (t *Transport) process_tx() (Message, bool) {
	m := Message{}
	//output_msg := false

	if t.pending_flow_control_tx {
		t.pending_flow_control_tx = false
		return t.make_flow_control(CONTINUE), true
	}

	flow_control_frame := t.last_flow_control_frame

	if flow_control_frame != nil {
		if flow_control_frame.flow == OVERFLOW {
			t.stop_sending()
			return m, false
		}
		if t.tx_state == IDLE {
			fmt.Println("Received a FlowControl message while transmission was Idle.")
		} else {
			if flow_control_frame.flow == WAIT {
				if t.wftmax == 0 {
					fmt.Println("Received a FlowControl requesting to wait, but fwtmax is set to 0")
				} else if t.wft_counter >= t.wftmax {
					fmt.Println("Received d wait frame which is the maximum set in params.wftmax")
					t.stop_sending()
				} else {
					t.wft_counter += 1
					if t.tx_state == WAIT || t.tx_state == TRANSMIT {
						t.tx_state = WAIT
						t.timer_rx_fc.start()
					}
				}
			} else if flow_control_frame.flow == CONTINUE && !t.timer_rx_fc.is_timed_out() {
				t.wft_counter = 0
				t.timer_rx_fc.stop()
				//t.timer_tx_stmin.set_timeout(flow_control_frame.stmin_sec)
				t.remote_blocksize = flow_control_frame.blocksize
				if t.tx_state == WAIT {
					t.tx_block_counter = 0
					t.timer_tx_stmin.start()
				}
				t.tx_state = TRANSMIT
			}
		}
	}

	if t.timer_rx_fc.is_timed_out() {
		fmt.Println("Reception of FlowControl timed out. Stopping transmission")
		t.stop_sending()
	}

	if t.tx_state != IDLE && len(t.tx_buffer) == 0 {
		t.stop_sending()
	}

	if t.tx_state == IDLE {
		if t.tx_queue.Len() > 0 {
			payload := t.tx_queue.Get()
			t.tx_buffer = payload
			size_on_first_byte := false
			if len(t.tx_buffer) <= 7 {
				size_on_first_byte = true
			}
			size_offset := 2
			if size_on_first_byte {
				size_offset = 1
			}
			msg_data := []byte{}
			if len(t.tx_buffer) <= t.data_length-size_offset-len(t.address.tx_payload_prefix) {
				if size_on_first_byte {
					msg_data = append([]byte{byte(0x0 | len(t.tx_buffer))}, t.tx_buffer...)
				} else {
					msg_data = append([]byte{0x0, byte(len(t.tx_buffer))}, t.tx_buffer...)
				}

			} else {
				t.tx_frame_length = len(t.tx_buffer)
				encode_length_on_2_first_bytes := false
				if t.tx_frame_length <= 4095 {
					encode_length_on_2_first_bytes = true
				}
				data_length := 0
				if encode_length_on_2_first_bytes {
					data_length = t.data_length - 2 - len(t.address.tx_payload_prefix)
					msg_data = append([]byte{0x10 | byte((t.tx_frame_length>>8)&0xF), byte(t.tx_frame_length & 0xFF)}, t.tx_buffer[:data_length]...)
				} else {
					data_length = t.data_length - 6 - len(t.address.tx_payload_prefix)
					msg_data = append([]byte{0x10, 0x00, byte(t.tx_frame_length>>24) & 0xFF, byte(t.tx_frame_length>>16) & 0xFF, byte(t.tx_frame_length>>8) & 0xFF, byte(t.tx_frame_length>>0) & 0xFF}, t.tx_buffer[:data_length]...)
				}
				t.tx_buffer = t.tx_buffer[data_length:]
				t.tx_state = WAIT
				t.tx_seqnum = 1
				t.timer_rx_fc.start()

			}
			m = t.make_tx_msg(t.address.txid, msg_data)
			return m, true
		}
	} else if t.tx_state == WAIT {
	} else if t.tx_state == TRANSMIT {
		if t.timer_tx_stmin.is_timed_out() || t.squash_stmin_requirement {
			/*
							data_length = self.params.ll_data_length-1-len(self.address.tx_payload_prefix)
				        msg_data = self.address.tx_payload_prefix + bytearray([0x20 | self.tx_seqnum]) + self.tx_buffer[:data_length]
				        arbitration_id  = self.address.get_tx_arbitraton_id()
				        output_msg = self.make_tx_msg(arbitration_id, msg_data)
				        self.tx_buffer = self.tx_buffer[data_length:]
				        self.tx_seqnum = (self.tx_seqnum + 1 ) & 0xF
				        self.timer_tx_stmin.start()
				        self.tx_block_counter+=1
			*/
		}
		if t.remote_blocksize != 0 && t.tx_block_counter >= t.remote_blocksize {
			t.tx_state = WAIT
			t.timer_rx_fc.start()
		}

	}
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
