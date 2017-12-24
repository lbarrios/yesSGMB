package cpu

import "github.com/lbarrios/yesSGMB/mmu"

type word uint16

type Flags struct {
	z bool // zero flag
	n bool // subtract flag
	h bool // half carry flag
	c bool // carry flag
	_ bool // to complete 1 byte
	_ bool // to complete 1 byte
	_ bool // to complete 1 byte
	_ bool // to complete 1 byte
}

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
	sp word // stack pointer
	pc word // program counter
}

func (r Registers) bcAsWord() word {
	res := word(r.bc.b) << 8
	res += word(r.bc.c)
	return res
}

func (r Registers) deAsWord() word {
	res := word(r.de.d) << 8
	res += word(r.de.e)
	return res
}

func (r Registers) hlAsWord() word {
	res := word(r.hl.h) << 8
	res += word(r.hl.l)
	return res
}

func (r Registers) hlAsAddress() mmu.Address {
	addr := mmu.Address{High: r.hl.h, Low: r.hl.l}
	return addr
}

func (r Registers) spLow() byte {
	return byte(r.sp & 0x0f)
}

func (r Registers) spHigh() byte {
	return byte(r.sp >> 8)
}

func (r Registers) pcLow() byte {
	return byte(r.sp & 0x0f)
}

func (r Registers) pcHigh() byte {
	return byte(r.sp >> 8)
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

func (r Registers) flagAsByte(flag bool) byte {
	if flag {
		return 1
	} else {
		return 0
	}
}
