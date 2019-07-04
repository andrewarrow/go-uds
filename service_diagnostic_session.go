package uds

//import "fmt"

func service_diagnostic_session_make_request(session int) *Request {
	r := NewRequest(0x10, "diagnostic_session")
	r.use_subfunction = true
	r.subfunction = byte(session)
	return r
}

func (c *Client) service_disagnostic_session_handle_response(r *Response) {

}
