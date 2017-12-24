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

func (f *Flags) asByte() byte {
	var r byte
	if f.z {
		r += 1 >> 8
	}
	if f.n {
		r += 1 >> 7
	}
	if f.h {
		r += 1 >> 6
	}
	if f.c {
		r += 1 >> 5
	}
	return r
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

/**
	Registers as Word
 */
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

/**
	Registers as Address
 */
func (r Registers) bcAsAddress() mmu.Address {
	addr := mmu.Address{High: r.bc.b, Low: r.bc.c}
	return addr
}
func (r Registers) deAsAddress() mmu.Address {
	addr := mmu.Address{High: r.de.d, Low: r.de.e}
	return addr
}
func (r Registers) hlAsAddress() mmu.Address {
	addr := mmu.Address{High: r.hl.h, Low: r.hl.l}
	return addr
}
func (r Registers) spAsAddress() mmu.Address {
	addr := mmu.Address{High: r.sp.high(), Low: r.sp.low()}
	return addr
}

/**
	Registers as Low or High nibbles
	(get low or high nibbles of a Byte)
 */
func lowNibble(b byte) byte {
	return byte(b & 0x0F)
}

func highNibble(b byte) byte {
	return byte((b & 0xF0) >> 4)
}

/**
	Registers as Low or High bytes
	(get low or high parts of a Word)
 */
func (w *word) low() byte {
	return byte(*w & 0xFF)
}
func (w *word) high() byte {
	return byte(*w >> 8)
}

/*
	Set Flags
 */
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

/*
	Get Flags
 */
func (r Registers) flagAsByte(flag bool) byte {
	if flag {
		return 1
	} else {
		return 0
	}
}
