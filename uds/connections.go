package isotp

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
