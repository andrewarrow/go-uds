package isotp

const Normal_11bits = 0
const Normal_29bits = 1
const NormalFixed_29bits = 2
const Extended_11bits = 3
const Extended_29bits = 4
const Mixed_11bits = 5
const Mixed_29bits = 6

const Physical = 0
const Functional = 1

type Address struct {
	rx_prefix_size    int
	addressing_mode   int
	txid              int
	rxid              int
	target_address    int
	source_address    int
	address_extension int
	tx_payload_prefix []byte
}

func NewAddress(rxid, txid int) Address {
	a := Address{}
	a.tx_payload_prefix = []byte{}
	a.addressing_mode = Normal_11bits
	//a.target_address = target_address
	//a.source_address = source_address
	//a.address_extension = address_extension
	a.rxid = rxid
	a.txid = txid

	return a
}

func (a Address) is_for_me(msg Message) bool {
	if a.is_29bits() == msg.extended_id {
		return msg.arbitration_id == a.rxid
	}
	return false
}

func (a Address) _get_tx_arbitraton_id(address_type int) int {
	if a.addressing_mode == Normal_11bits {
		return a.txid
	} else if a.addressing_mode == Normal_29bits {
		return a.txid
	} else if a.addressing_mode == NormalFixed_29bits {
		bits23_16 := 0xDB0000
		if address_type == Physical {
			bits23_16 = 0xDA0000
		}
		return 0x18000000 | bits23_16 | (a.target_address << 8) | a.source_address
	} else if a.addressing_mode == Extended_11bits {
		return a.txid
	} else if a.addressing_mode == Extended_29bits {
		return a.txid
	} else if a.addressing_mode == Mixed_11bits {
		return a.txid
	} else if a.addressing_mode == Mixed_29bits {
		bits23_16 := 0xCD0000
		if address_type == Physical {
			bits23_16 = 0xCE0000
		}
		return 0x18000000 | bits23_16 | (a.target_address << 8) | a.source_address
	}
	return 0
}

func (a Address) is_29bits() bool {
	if a.addressing_mode == Normal_29bits ||
		a.addressing_mode == NormalFixed_29bits ||
		a.addressing_mode == Extended_29bits ||
		a.addressing_mode == Mixed_29bits {
		return true
	}
	return false
}
