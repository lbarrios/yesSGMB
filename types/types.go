package types

type Address struct {
	High byte
	Low  byte
}

func (address Address) NextAddress() Address {
	addr := address
	addr.Low += 1
	if addr.Low == 0x00 {
		addr.High += 1
	}
	return addr
}

func (address Address) AsWord() uint16 {
	var ret uint16
	ret = uint16(address.High)<<8 + uint16(address.Low)
	return ret
}

func (address Address) String() string {
	return string(address.AsWord())
}