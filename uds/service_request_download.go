package uds

import "fmt"

func service_request_download_make_request(ml MemoryLocation) *Request {

	r := NewRequest(0x34, "request_download")
	r.use_subfunction = false
	r.data = append(r.data, DataFormatId())
	r.data = append(r.data, ml.AlfidByte())
	r.data = append(r.data, ml.GetAddressBytes()...)
	r.data = append(r.data, ml.GetMemorySizeBytes()...)

	fmt.Println("request_download", r.data)
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
