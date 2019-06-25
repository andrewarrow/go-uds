package uds

import "fmt"
import "github.com/andrewarrow/go-uds/isotp"

type Client struct {
	conn                       isotp.AnyConn
	timeout                    float32
	suppress_positive_response bool
}

func NewClient(connection isotp.AnyConn, timeout float32) *Client {
	c := Client{}
	c.conn = connection
	c.suppress_positive_response = true
	return &c
}

func (c *Client) send_request(request *Request) *Response {
	fmt.Println(111)
	c.conn.Empty_rxqueue()
	payload := []byte{}
	//override_suppress_positive_response := false
	if c.suppress_positive_response && request.use_subfunction {
		payload = request.get_payload(true)
		//override_suppress_positive_response = true
	} else {
		payload = request.get_payload(false)
	}
	fmt.Println(888, payload)
	c.conn.Send(payload)
	data := c.conn.Wait_frame()
	fmt.Println(777, data)
	response := response_from_payload(data)
	return response
}

func (c *Client) Transfer_data(seqnum byte, data []byte) *Response {
	req := service_transfer_make_request(seqnum, data)
	fmt.Println(req)
	ponse := c.send_request(req)
	fmt.Println(999, ponse)
	service_transfer_handle_response(ponse)
	return ponse
}
func (c *Client) Read_data_by_id(data []byte) *Response {
	req := service_read_data_by_id_make_request(data)
	fmt.Println(req)
	ponse := c.send_request(req)
	fmt.Println(999, ponse)
	service_read_data_by_id_handle_response(ponse)
	return ponse
}
