package uds

type Response struct {
	flavor string
}

func NewResponse() *Response {
	r := Response{}
	return &r
}
