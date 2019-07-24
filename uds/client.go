package uds

import "fmt"
import "github.com/andrewarrow/go-uds/isotp"
import "encoding/binary"

type Client struct {
	conn                       isotp.AnyConn
	timeout                    float32
	suppress_positive_response bool
	Data_identifiers           map[int]int
	Release                    func() int
}

func NewClient(connection isotp.AnyConn, timeout float32, rf func() int) *Client {
	c := Client{}
	c.conn = connection
	c.suppress_positive_response = true
	c.Data_identifiers = map[int]int{}
	c.Release = rf
	return &c
}
func (c *Client) send_request(request *Request) *Response {
	//fmt.Println(111)
	c.conn.Empty_rxqueue()
	payload := []byte{}
	//override_suppress_positive_response := false
	if c.suppress_positive_response && request.use_subfunction {
		payload = request.get_payload(true)
		//override_suppress_positive_response = true
	} else {
		payload = request.get_payload(false)
	}
	//fmt.Println(888, payload)
	c.conn.Send(payload)
	data := c.conn.Wait_frame()
	//fmt.Println(777, data)
	response := response_from_payload(data)
	return response
}

func (c *Client) Read_data_by_id(data []int) *Response {
	req := service_read_data_by_id_make_request(data)
	response := c.send_request(req)
	c.service_read_data_by_id_handle_response(response)
	return response
}
func (c *Client) Simple_read_data_by_id(did, length int, flavor string) string {
	request := service_read_data_by_id_make_request([]int{did})
	payload := request.get_payload(false)
	data := c.conn.Send_and_grant_flow_request(payload, length)
	if flavor == "text" {
		if len(data) > 5 {
			return string(data[5:])
		}
		return string(data)
	}
	if len(data) > 5 {
		return fmt.Sprintf("%v", data[5:])
	}
	return fmt.Sprintf("%v", data)
}
func (c *Client) Request_download(ml MemoryLocation) int {
	request := service_request_download_make_request(ml)
	payload := request.get_payload(false)
	fmt.Println("client, Request_download", payload)
	data := c.conn.Send_and_wait_for_reply(payload)
	data = data[2:]

	fmt.Println(data)
	todecode := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	lfid := int(data[0] >> 4)
	i := 1
	for {
		if i >= lfid+1 {
			break
		}
		todecode[7-(i-1)] = data[lfid+1-i]
		i += 1
	}

	fmt.Println(lfid, todecode, int(binary.BigEndian.Uint64(todecode)))
	return int(binary.BigEndian.Uint64(todecode))
}
func (c *Client) Transfer_data(i int, data []byte) string {
	request := service_transfer_data_make_request(i, data)
	payload := request.get_payload(false)
	response := c.conn.Send_and_no_wait_for_reply(payload)
	return fmt.Sprintf("%v", response)
}
func (c *Client) Request_transfer_exit(crc int) string {
	r := NewRequest(0x37, "transfer_exit")
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, uint32(crc))
	fmt.Println(bs)
	r.data = bs
	r.use_subfunction = false
	payload := r.get_payload(false)
	c.conn.Send_and_no_wait_for_reply(payload)
	return fmt.Sprintf("%v", crc)
}
func (c *Client) Change_session(session int) string {
	request := service_diagnostic_session_make_request(session)
	payload := request.get_payload(false)
	data := c.conn.Send_and_wait_for_reply(payload)
	return fmt.Sprintf("%v", data)
}
func (c *Client) Unlock_security_access(level int, algo func(seed []byte, params int) []byte) string {
	request := service_security_access_make_request(level, "request_seed", []byte{})
	payload := request.get_payload(false)
	seed := c.conn.Send_and_wait_for_reply(payload)
	//TODO review this result
	fmt.Printf("%v\n", seed)

	request = service_security_access_make_request(level, "sendkey", algo(seed, 0))
	payload = request.get_payload(false)
	data := c.conn.Send_and_wait_for_reply(payload)
	fmt.Printf("%v\n", data)
	//TODO review this result
	return ""
}
