package uds

func service_transfer_make_request(seqnum byte, data []byte) *Request {
	r := NewRequest("transfer")
	r.data = []byte{seqnum}
	r.data = append(r.data, data...)
	return r
}

func interpret_response() {
}
