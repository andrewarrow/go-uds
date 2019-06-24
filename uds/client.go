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
	response := response_from_payload(payload)
	return response
}

func (c *Client) transfer_data(seqnum byte, data []byte) *Response {
	req := service_transfer_make_request(seqnum, data)

	ponse := c.send_request(req)
	/*

			     data_len = 0 if data is None else len(data)
		                self.logger.info('%s - Sending a block of data with SequenceNumber=%d that is %d bytes long .' % (self.service_log_prefix(services.TransferData), sequence_number, data_len))
		                if data is not None:
		                        self.logger.debug('Data to transfer : %s' % binascii.hexlify(data))

		                response = self.send_request(request)
		                if response is None:
		                        return
		                services.TransferData.interpret_response(response)
	*/
	return ponse
}
