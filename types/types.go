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

func (address Address) AsWord() Word {
	var ret Word
	ret = Word(address.High)<<8 + Word(address.Low)
	return ret
}

func (address Address) String() string {
	return string(address.AsWord())
}

type Word uint16

/**
Return the Low byte of a Word
*/
func (w Word) Low() byte {
	return byte(w & 0xFF)
}

/**
Return the Low byte of a Word
*/
func (w Word) High() byte {
	return byte(w >> 8)
}

func (w Word) AsAddress() Address {
	return Address{High: w.High(), Low: w.Low()}
}

func WordFromBytes(high byte, low byte) Word {
	return (Word(high)<<8 + Word(low))
}

func BitIsSet(b byte, pos uint) bool {
	return b&(1<<pos) == (1 << pos)
}
