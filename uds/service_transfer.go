package uds

func service_transfer_data_make_request(i int, data []byte) *Request {
	r := NewRequest(0x36, "transfer_data")
	r.data = []byte{byte(i)}
	r.data = append(r.data, data...)
	r.use_subfunction = false
	return r
}

func service_transfer_data_handle_response(r *Response) {

	r.Service_data["sequence_number_echo"] = r.Data[0]
	r.Service_data["parameter_records"] = []byte{}
	if len(r.Data) > 1 {
		r.Service_data["parameter_records"] = r.Data[1:]
	}

}
