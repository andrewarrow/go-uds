package uds

func service_request_download_make_request(ml MemoryLocation) *Request {
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

}
