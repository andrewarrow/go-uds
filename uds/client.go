package uds

type Client struct {
	connection interface{}
	timeout    float32
}

func NewClient(connection interface{}, timeout float32) *Client {
	c := Client{}
	c.connection = connection
	return &c
}

func (c *Client) transfer_data(seqnum int, data []byte) *Response {
	r := service_transfer_make_request(seqnum, data)
	return r
}
