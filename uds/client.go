package uds

//import "fmt"
import "github.com/andrewarrow/go-isotp/isotp"

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
	c.conn.Empty_rxqueue()
	payload := []byte{}
	//override_suppress_positive_response := false
	if c.suppress_positive_response && request.use_subfunction {
		payload = request.get_payload(true)
		//override_suppress_positive_response = true
	} else {
		payload = request.get_payload(false)
	}
	c.conn.Send(payload)
	data := c.conn.Wait_frame()
	response := response_from_payload(data)
	return response
}

func (c *Client) transfer_data(seqnum byte, data []byte) *Response {
	req := service_transfer_make_request(seqnum, data)
	ponse := c.send_request(req)
	service_transfer_handle_response(ponse)
	return ponse
}
