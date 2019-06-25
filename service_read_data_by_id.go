package uds

import "fmt"
import "encoding/binary"

func convertIntArrayToByteArray(data []int) []byte {

	buff := []byte{}
	for _, short := range data {
		bs := make([]byte, 2)
		binary.BigEndian.PutUint16(bs, uint16(short))
		buff = append(buff, bs...)
	}
	return buff
}

func service_read_data_by_id_make_request(data []int) *Request {
	r := NewRequest(0x22, "read_data_by_id")
	r.data = convertIntArrayToByteArray(data)
	r.use_subfunction = false
	return r
}

func service_read_data_by_id_handle_response(r *Response) {

	fmt.Println(r)

}
