package uds

//import "fmt"

func service_security_access_make_request(level int, mode string, key []byte) *Request {
	r := NewRequest(0x27, "security")
	r.use_subfunction = true

	if mode == "request_seed" {
		if level%2 == 1 {
			r.subfunction = byte(level)
		} else {
			r.subfunction = byte(level) - 1
		}
	} else if mode == "send_key" {
		if level%2 == 0 {
			r.subfunction = byte(level)
		} else {
			r.subfunction = byte(level) + 1
		}
	}

	if mode == "send_key" {
		r.data = append(r.data, key...)
	}
	return r
}
