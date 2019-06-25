package uds

import "fmt"

func service_read_data_by_id_make_request(seqnum byte, data []byte) *Request {
	r := NewRequest(0x22, "read_data_by_id")
	r.data = []byte{seqnum}
	r.data = append(r.data, data...)
	r.use_subfunction = false
	return r
}

func service_read_data_by_id_handle_response(r *Response) {

	fmt.Println(r)

}
