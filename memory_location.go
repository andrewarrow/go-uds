package uds

type MemoryLocation struct {
	address           int
	length            int
	address_format    int
	memorysize_format int
}

func NewMemoryLocation(address, length, address_format, memorysize_format int) *MemoryLocation {
	m := MemoryLocation{}
	m.address = address
	m.address_format = address_format
	m.length = length
	m.memorysize_format = memorysize_format
	return &m
}
