package uds

type Request struct {
	service         string
	data            []byte
	sid             byte
	subfunction     byte
	use_subfunction bool
}

func NewRequest(sid byte, service string) *Request {
	r := Request{}
	r.sid = sid
	r.service = service
	r.data = []byte{}
	return &r
}

func (r *Request) get_payload(suppress_positive_response bool) []byte {
	payload := []byte{r.sid}
	if r.use_subfunction {
		r.data = append([]byte{r.subfunction}, r.data...)
	} else {
	}
	return append(payload, r.data...)
}
