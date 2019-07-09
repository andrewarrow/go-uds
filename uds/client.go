package uds

import "fmt"
import "github.com/andrewarrow/go-uds/isotp"
import "encoding/binary"

type Client struct {
	conn                       isotp.AnyConn
	timeout                    float32
	suppress_positive_response bool
	Data_identifiers           map[int]int
	Release_init               func() int
}

func NewClient(connection isotp.AnyConn, timeout float32, ri func() int) *Client {
	c := Client{}
	c.conn = connection
	c.suppress_positive_response = true
	c.Data_identifiers = map[int]int{}
	c.Release_init = ri
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
		return string(data[5:])
	}
	return fmt.Sprintf("%v", data[5:])
}
func (c *Client) Request_download(ml MemoryLocation) int {
	request := service_request_download_make_request(ml)
	payload := request.get_payload(false)
	fmt.Println("client, Request_download", payload)
	data := c.conn.Send_and_wait_for_reply(payload)

	todecode := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	lfid := int(data[0] >> 4)
	i := 1
	for {
		if i > lfid+1 {
			break
		}
		todecode[7-(i-1)] = data[lfid+1-i]
		i += 1
	}

	return int(binary.BigEndian.Uint32(todecode))
}
func (c *Client) Transfer_data(i int, data []byte) string {
	request := service_transfer_data_make_request(i, data)
	payload := request.get_payload(false)
	response := c.conn.Send_and_wait_for_reply(payload)
	return fmt.Sprintf("%v", response)
}
func (c *Client) Request_transfer_exit(crc int) string {
	return fmt.Sprintf("%v", "")
}
func (c *Client) Change_session(session int) string {
	request := service_diagnostic_session_make_request(session)
	payload := request.get_payload(false)
	data := c.conn.Send_and_wait_for_reply(payload)
	return fmt.Sprintf("%v", data)
}
