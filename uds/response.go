package uds

type Response struct {
	service      string
	data         []byte
	service_data map[string]interface{}
}

func NewResponse(service string) *Response {
	r := Response{}
	r.service = service
	r.data = []byte{}
	r.service_data = map[string]interface{}{}
	return &r
}
