package uds

func service_transfer_make_request(seqnum byte, data []byte) *Response {
	r := NewResponse("transfer")
	r.data = []byte{seqnum}
	r.data = append(r.data, data...)
	return r
}
