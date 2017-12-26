package cpu

import (
	"fmt"
	"github.com/lbarrios/yesSGMB/types"
)

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

func (f Flags) String() string {
	z := 0
	n := 0
	h := 0
	c := 0
	if f.z {
		z = 1
	}
	if f.n {
		n = 1
	}
	if f.h {
		h = 1
	}
	if f.c {
		c = 1
	}
	return fmt.Sprintf("{z:%d n:%d h:%d c:%d}", z, n, h, c)
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
	sp types.Word // stack pointer
	pc types.Word // program counter
}

// String prints registers as string
func (r Registers) String() string {
	formatByte := func(b byte) string {
		if b < 0x10 {
			return fmt.Sprintf("0x0%X", b)
		}
		return fmt.Sprintf("0x%X", b)
	}
	return fmt.Sprintf("A: %s  B: %s  C: %s  D: %s  E: %s  H: %s  L: %s SP: %x PC: %x f: %s",
		formatByte(r.af.a), formatByte(r.bc.b), formatByte(r.bc.c),
		formatByte(r.de.d), formatByte(r.de.e), formatByte(r.hl.h), formatByte(r.hl.l),
		r.sp, r.pc, r.af.f)
}

/**
Registers as Word
*/
func (r Registers) bcAsWord() types.Word {
	res := types.Word(r.bc.b) << 8
	res += types.Word(r.bc.c)
	return res
}

func (r Registers) deAsWord() types.Word {
	res := types.Word(r.de.d) << 8
	res += types.Word(r.de.e)
	return res
}

func (r Registers) hlAsWord() types.Word {
	res := types.Word(r.hl.h) << 8
	res += types.Word(r.hl.l)
	return res
}

/**
Registers as Address
*/
func (r Registers) bcAsAddress() types.Address {
	addr := types.Address{High: r.bc.b, Low: r.bc.c}
	return addr
}
func (r Registers) deAsAddress() types.Address {
	addr := types.Address{High: r.de.d, Low: r.de.e}
	return addr
}
func (r Registers) hlAsAddress() types.Address {
	addr := types.Address{High: r.hl.h, Low: r.hl.l}
	return addr
}
func (r Registers) spAsAddress() types.Address {
	addr := types.Address{High: r.sp.High(), Low: r.sp.Low()}
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
