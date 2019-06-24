package uds

import "fmt"
import "strings"

type Client struct {
	conn                       interface{}
	timeout                    float32
	suppress_positive_response bool
}

func NewClient(connection interface{}, timeout float32) *Client {
	c := Client{}
	c.conn = connection
	c.suppress_positive_response = true
	return &c
}

func Empty_rxqueue(conn interface{}) {
	name := fmt.Sprintf("%v", conn)
	tokens := strings.Split(name[2:], " ")
	fmt.Println(tokens[0])
	if tokens[0] == "queue" {
		qc := conn.(*QueueConnection)
		qc.Empty_rxqueue()
	}
}
func Wait_frame(conn interface{}) []byte {
	name := fmt.Sprintf("%v", conn)
	tokens := strings.Split(name[2:], " ")

	if tokens[0] == "queue" {
		qc := conn.(*QueueConnection)
		return qc.Wait_frame()
	}

	b := []byte{}
	return b
}
func Send(conn interface{}, data []byte) {
	name := fmt.Sprintf("%v", conn)
	tokens := strings.Split(name[2:], " ")
	if tokens[0] == "queue" {
		qc := conn.(*QueueConnection)
		qc.Send(data)
	}
}

func (c *Client) send_request(request *Request) *Response {
	Empty_rxqueue(c.conn)
	payload := []byte{}
	//override_suppress_positive_response := false
	if c.suppress_positive_response && request.use_subfunction {
		payload = request.get_payload(true)
		//override_suppress_positive_response = true
	} else {
		payload = request.get_payload(false)
	}
	Send(c.conn, payload)
	data := Wait_frame(c.conn)
	response := response_from_payload(data)
	return response
}

func (c *Client) transfer_data(seqnum byte, data []byte) *Response {
	req := service_transfer_make_request(seqnum, data)
	ponse := c.send_request(req)
	service_transfer_handle_response(ponse)
	return ponse
}
