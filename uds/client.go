package uds

type Client struct {
	conn    AnyConn
	timeout float32
}

func NewClient(connection AnyConn, timeout float32) *Client {
	c := Client{}
	c.conn = connection
	return &c
}

func (c *Client) send_request(request *Request) *Response {
	c.conn.empty_rxqueue()
	payload := request.get_payload()
	c.conn.send(payload)
	data := c.conn.wait_frame()
	response := response_from_payload(data)
	return response
}

func (c *Client) transfer_data(seqnum byte, data []byte) *Response {
	req := service_transfer_make_request(seqnum, data)

	ponse := c.send_request(req)
	service_transfer_handle_response(ponse)
	return ponse
}
