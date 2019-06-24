package uds

type Response struct {
	service string
	data    []byte
}

func NewResponse(service string) *Response {
	r := Response{}
	r.service = service
	r.data = []byte{}
	return &r
}
