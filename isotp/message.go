package isotp

//import "fmt"

type Message struct {
	Id          int64
	Payload     []byte
	extended_id bool
}

func NewMessage(id int64, payload []byte) Message {
	m := Message{}
	m.Id = id
	m.Payload = payload
	return m
}

func (m Message) ToBytes() []byte {
	return append([]byte{byte(m.Id)}, m.Payload...)
}

type PDU struct {
	msg       Message
	flavor    string
	length    int
	payload   []byte
	blocksize int
	stmin     int
	stmin_sec int
	seqnum    int
	flow      string
}

func NewPDU(msg Message, start_of_data int, data_length int) PDU {
	pdu := PDU{}
	pdu.msg = msg

	if len(msg.Payload) > start_of_data {
		h := (msg.Payload[start_of_data] >> 4) & 0xF
		if h == 0 {
			pdu.flavor = SINGLE
		} else if h == 1 {
			pdu.flavor = FIRST
		} else if h == 2 {
			pdu.flavor = CONSECUTIVE
		} else if h == 3 {
			pdu.flavor = FLOW
		}
	}
	if pdu.flavor == SINGLE {
		pdu.length = int(msg.Payload[start_of_data]) & 0xF
		pdu.payload = msg.Payload[0+start_of_data : len(msg.Payload)][1 : pdu.length+1]
	} else if pdu.flavor == FIRST {
		pdu.length = ((int(msg.Payload[start_of_data]) & 0xF) << 8) | int(msg.Payload[start_of_data+1])
		pdu.payload = msg.Payload[2+start_of_data : len(msg.Payload)][0:min(pdu.length, data_length-2-start_of_data)]
	} else if pdu.flavor == CONSECUTIVE {
		pdu.seqnum = int(msg.Payload[start_of_data]) & 0xF
		pdu.payload = msg.Payload[start_of_data+1 : data_length]
	} else if pdu.flavor == FLOW {
		f := int(msg.Payload[start_of_data]) & 0xF
		if f == 0 {
			pdu.flow = SINGLE
		} else if f == 1 {
			pdu.flow = FIRST
		} else if f == 2 {
			pdu.flow = CONSECUTIVE
		} else if f == 3 {
			pdu.flow = FLOW
		}

		pdu.blocksize = int(msg.Payload[1+start_of_data])
		stmin_temp := int(msg.Payload[2+start_of_data])

		if stmin_temp >= 0 && stmin_temp <= 0x7F {
			pdu.stmin_sec = stmin_temp / 1000
		} else if stmin_temp >= 0xf1 && stmin_temp <= 0xF9 {
			pdu.stmin_sec = (stmin_temp - 0xF0) / 10000
		}
		pdu.stmin = stmin_temp

	}
	return pdu
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
