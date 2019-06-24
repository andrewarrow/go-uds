package uds

func service_transfer_make_request(seqnum byte, data []byte) *Request {
	r := NewRequest(0x36, "transfer")
	r.data = []byte{seqnum}
	r.data = append(r.data, data...)
	r.use_subfunction = false
	return r
}

func service_transfer_handle_response(r *Response) {

	r.service_data["sequence_number_echo"] = r.data[0]
	r.service_data["parameter_records"] = []byte{}
	if len(r.data) > 1 {
		r.service_data["parameter_records"] = r.data[1:]
	}

}
