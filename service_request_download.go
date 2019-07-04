package uds

func service_request_download_make_request(ml MemoryLocation) *Request {

	ml.AlfidByte()

	r := NewRequest(0x34, "request_download")
	r.use_subfunction = true

	/*
		request.data += memory_location.alfid.get_byte()        # AddressAndLengthFormatIdentifier
		request.data += memory_location.get_address_bytes()
		request.data += memory_location.get_memorysize_bytes()
	*/
	return r
}

func (c *Client) service_request_download_handle_response(r *Response) {

	//lfid = int(response.data[0]) >> 4
	/*
		todecode = bytearray(b'\x00\x00\x00\x00\x00\x00\x00\x00')
		for i in range(1,lfid+1):
		  todecode[-i] = response.data[lfid+1-i]
		response.service_data.max_length = struct.unpack('>q', todecode)[0]
	*/
}
