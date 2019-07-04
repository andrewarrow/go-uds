package uds

var address_map = map[int]int{
	8:  1,
	16: 2,
	24: 3,
	32: 4,
	40: 5,
}
var memsize_map = map[int]int{
	8:  1,
	16: 2,
	24: 3,
	32: 4,
}

type MemoryLocation struct {
	address           int
	length            int
	address_format    int
	memorysize_format int
}

func NewMemoryLocation(address, length, address_format, memorysize_format int) *MemoryLocation {
	m := MemoryLocation{}
	m.address = address
	m.length = length
	m.address_format = address_format
	m.memorysize_format = memorysize_format
	return &m
}
func (m *MemoryLocation) AlfidByte() int {
	return ((memsize_map[m.memorysize_format] << 4) | (address_map[m.address_format])) & 0xFF
}

/*

	def set_format_if_none(self, address_format=None, memorysize_format=None):
		previous_address_format = self.address_format
		previous_memorysize_format = self.memorysize_format
		try:
			if address_format is not None:
				if self.address_format is None:
					self.address_format = address_format

			if memorysize_format is not None:
				if self.memorysize_format is None:
					self.memorysize_format=memorysize_format

			address_format = self.address_format if self.address_format is not None else self.autosize_address(self.address)
			memorysize_format = self.memorysize_format if self.memorysize_format is not None else self.autosize_memorysize(self.memorysize)

			self.alfid = AddressAndLengthFormatIdentifier(memorysize_format=memorysize_format, address_format=address_format)
		except:
			self.address_format = previous_address_format
			self.memorysize_format = previous_memorysize_format
			raise

	def autosize_address(self, val):
		fmt = math.ceil(val.bit_length()/8)*8
		if fmt > 40:
			raise ValueError("address size must be smaller or equal than 40 bits")
		return fmt

	def autosize_memorysize(self, val):
		fmt = math.ceil(val.bit_length()/8)*8
		if fmt > 32:
			raise ValueError("memory size must be smaller or equal than 32 bits")
		return fmt

	def get_address_bytes(self):
		n = AddressAndLengthFormatIdentifier.address_map[self.alfid.address_format]

		data = struct.pack('>q', self.address)
		return data[-n:]


	def get_memorysize_bytes(self):
		n = AddressAndLengthFormatIdentifier.memsize_map[self.alfid.memorysize_format]

		data = struct.pack('>q', self.memorysize)
		return data[-n:]

	def from_bytes(cls, address_bytes, memorysize_bytes):
		if not isinstance(address_bytes, bytes):
			raise ValueError('address_bytes must be a valid bytes object')

		if not isinstance(memorysize_bytes, bytes):
			raise ValueError('memorysize_bytes must be a valid bytes object')

		if len(address_bytes) > 5:
			raise ValueError('Address must be at most 40 bits long')

		if len(memorysize_bytes) > 4:
			raise ValueError('Memory size must be at most 32 bits long')

		address_bytes_padded = b'\x00' * (8-len(address_bytes)) + address_bytes
		memorysize_bytes_padded = b'\x00' * (8-len(memorysize_bytes)) + memorysize_bytes

		address = struct.unpack('>q', address_bytes_padded)[0]
		memorysize = struct.unpack('>q', memorysize_bytes_padded)[0]
		address_format = len(address_bytes) * 8
		memorysize_format = len(memorysize_bytes) * 8

		return cls(address=address, memorysize=memorysize, address_format=address_format, memorysize_format=memorysize_format)

	def __str__(self):
		return 'Address=0x%x (%d bits), Size=0x%x (%d bits)' % (self.address, self.alfid.address_format, self.memorysize, self.alfid.memorysize_format)

	def __repr__(self):
		return '<%s: %s at 0x%08x>' % (self.__class__.__name__, str(self), id(self))
*/
