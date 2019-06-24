package isotp

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
	r := NewResponse()
	return r
}
