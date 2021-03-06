package uds

//import "fmt"
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

func (c *Client) service_read_data_by_id_handle_response(r *Response) {

	/*
		did := 0
		//did := binary.BigEndian.Uint16(r.Data[0:1])
		offset := 0
		for {
			if len(r.Data) <= offset {
				break
			}
			if len(r.Data) <= offset+1 {
				if true && r.Data[len(r.Data)-1] == 0 {
					break
				}
			}

			if did == 0 && true {
				compare := []byte{}
				for {
					compare = append(compare, 0)
					if len(compare) == len(r.Data)-offset {
						break
					}
				}
				if fmt.Sprintf("%v", r.Data[offset:]) == fmt.Sprintf("%v", compare) {
					break
				}
			}
			offset += 2

			codec_length := c.Data_identifiers[int(did)]
			subpayload := r.Data[offset:]
			if offset+codec_length < len(r.Data) {
				subpayload = r.Data[offset : offset+codec_length]
			}
			offset += codec_length
			r.Service_data[fmt.Sprintf("%d", did)] = subpayload
		}
	*/

}
