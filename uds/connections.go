package uds

type AnyConn interface {
	empty_rxqueue()
	send(payload []byte)
}

type QueueConnection struct {
	name string
	mtu  int
}

func NewQueueConnection(name string, mtu int) *QueueConnection {
	q := QueueConnection{}
	q.name = name
	q.mtu = mtu
	return &q
}

func (q *QueueConnection) empty_rxqueue() {
}
func (q *QueueConnection) send(payload []byte) {
}
