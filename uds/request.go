package uds

type Request struct {
	service string
	data    []byte
}

func NewRequest(service string) *Request {
	r := Request{}
	r.service = service
	r.data = []byte{}
	return &r
}

func (r *Request) get_payload() []byte {
	return r.data
}
