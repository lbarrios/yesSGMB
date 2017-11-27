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

func (r *Registers) setFlagZ(condition bool) {
	// Put true on FLAG Z if condition is true, else false
	r.af.f.z = condition
}

func (r *Registers) setFlagN(condition bool) {
	// Put true on FLAG N if condition is true, else false
	r.af.f.n = condition
}

func (r *Registers) setFlagH(condition bool) {
	// Put true on FLAG H if condition is true, else false
	r.af.f.h = condition
}

func (r *Registers) setFlagC(condition bool) {
	// Put true on FLAG C if condition is true, else false
	r.af.f.c = condition
}

func (r Registers) flagAsInt(flag bool) uint8 {
	if flag {
		return 1
	} else {
		return 0
	}
}