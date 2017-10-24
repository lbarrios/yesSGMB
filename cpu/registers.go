package cpu

type Registers struct {
	af struct {
		a byte  // high
		f Flags // low - flags - not accessible by programmer
	}
	bc struct {
		b byte // high
		c byte // low
	}
	de struct {
		d byte // high
		e byte // low
	}
	hl struct {
		h byte // high
		l byte // low
	}
	sp uint16 // stack pointer
	pc uint16 // program counter
}

func (r Registers) bcAsInt() uint16 {
	res := uint16(r.bc.b) << 8
	res += uint16(r.bc.c)
	return res
}

func (r Registers) deAsInt() uint16 {
	res := uint16(r.de.d) << 8
	res += uint16(r.de.e)
	return res << 8
}

func (r Registers) hlAsInt() uint16 {
	res := uint16(r.hl.h) << 8
	res += uint16(r.hl.l)
	return res << 8
}