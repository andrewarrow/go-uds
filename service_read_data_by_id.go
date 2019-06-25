package uds

import "fmt"

func byteArrayWithZeros(data []byte) []byte {
	buff := []byte{}
	for _, b := range data {
		buff = append(buff, 0, b)
	}
	return buff
}

func service_read_data_by_id_make_request(data []byte) *Request {
	r := NewRequest(0x22, "read_data_by_id")
	r.data = byteArrayWithZeros(data)
	r.use_subfunction = false
	return r
}

func service_read_data_by_id_handle_response(r *Response) {

	fmt.Println(r)

}
