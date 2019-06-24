package isotp

import "time"

type Timer struct {
	timeout   float32
	startedAt int64
}

func NewTimer(timeout float32) *Timer {
	t := Timer{}
	t.timeout = timeout
	return &t
}

func (t *Timer) is_timed_out() bool {
	if t.startedAt == 0 {
		return false
	}
	if ((time.Now().UnixNano() / 1000 / 1000) - t.startedAt) > int64(t.timeout*1000.0) {
		return true
	}
	return false
}
func (t *Timer) stop() {
	t.startedAt = 0
}
func (t *Timer) start() {
	t.startedAt = time.Now().UnixNano() / 1000 / 1000
}
