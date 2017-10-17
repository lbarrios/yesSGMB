// package cpu implements the CPU and Registers type
package cpu

import (
	"log"
)

type Flags struct {
	z bool // zero flag
	n bool // substract flag
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
	sp int // stack pointer
	pc int // program counter
}

type cpu struct {
	r Registers
}

func NewCPU() *cpu {
	cpu := new(cpu)
	cpu.Reset()
	return cpu
}

func (cpu *cpu) Reset() {
	log.Println("CPU reset triggered.")
	cpu.r.pc = 0x0100 // On power up, the GameBoy Program Counter is initialized to 0x0100
	cpu.r.sp = 0xFFFE // On power up, the GameBoy Stack Pointer is initialized to 0xFFFE
	cpu.r.af.a = 0
	cpu.r.af.f = Flags{false,false,false,false,false,false,false,false}
	cpu.r.bc.b = 0
	cpu.r.bc.c = 0
	cpu.r.de.d = 0
	cpu.r.de.e = 0
	cpu.r.hl.h = 0
	cpu.r.hl.l = 0
}
