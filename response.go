package uds

type Response struct {
	Service      string
	Data         []byte
	Service_data map[string]interface{}
}

func NewResponse(service string) *Response {
	r := Response{}
	r.Service = service
	r.Data = []byte{}
	r.Service_data = map[string]interface{}{}
	return &r
}
func response_from_payload(data []byte) *Response {
	r := Response{}
	r.Data = data[1:]
	r.Service_data = map[string]interface{}{}
	return &r
}
