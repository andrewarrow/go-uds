package isotp

//import "fmt"

type Message struct {
	arbitration_id int
	dlc            int
	payload        []byte
	extended_id    bool
}

func NewMessage(arbitration_id int, payload []byte) Message {
	m := Message{}
	m.payload = payload
	m.arbitration_id = arbitration_id
	return m
}

func (m Message) GetData() []byte {
	return m.payload
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

	if len(msg.payload) > start_of_data {
		h := (msg.payload[start_of_data] >> 4) & 0xF
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
		pdu.length = int(msg.payload[start_of_data]) & 0xF
		pdu.payload = msg.payload[0+start_of_data : len(msg.payload)][1 : pdu.length+1]
	} else if pdu.flavor == FIRST {
		pdu.length = ((int(msg.payload[start_of_data]) & 0xF) << 8) | int(msg.payload[start_of_data+1])
		pdu.payload = msg.payload[2+start_of_data : len(msg.payload)][0:min(pdu.length, data_length-2-start_of_data)]
	} else if pdu.flavor == CONSECUTIVE {
		pdu.seqnum = int(msg.payload[start_of_data]) & 0xF
		pdu.payload = msg.payload[start_of_data+1 : data_length]
	} else if pdu.flavor == FLOW {
		f := int(msg.payload[start_of_data]) & 0xF
		if f == 0 {
			pdu.flow = SINGLE
		} else if f == 1 {
			pdu.flow = FIRST
		} else if f == 2 {
			pdu.flow = CONSECUTIVE
		} else if f == 3 {
			pdu.flow = FLOW
		}

		pdu.blocksize = int(msg.payload[1+start_of_data])
		stmin_temp := int(msg.payload[2+start_of_data])

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
