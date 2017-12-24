// package cpu implements the CPU and Registers type
package cpu

import (
	"log"
	"github.com/lbarrios/yesSGMB/mmu"
)

type cpu struct {
	r                 Registers
	mmu               mmu.MMU
	cycle             cycleCount
	interruptsEnabled bool
	halted            bool
}

func NewCPU(mmu mmu.MMU) *cpu {
	cpu := new(cpu)
	cpu.mmu = mmu
	cpu.Reset()
	return cpu
}

func (cpu *cpu) Reset() {
	log.Println("CPU reset triggered.")
	cpu.r.pc = 0x0100 // On power up, the GameBoy Program Counter is initialized to 0x0100
	cpu.r.sp = 0xFFFE // On power up, the GameBoy Stack Pointer is initialized to 0xFFFE
	cpu.r.af.a = 0
	cpu.r.af.f = Flags{false, false, false, false, false, false, false, false}
	cpu.r.bc.b = 0
	cpu.r.bc.c = 0
	cpu.r.de.d = 0
	cpu.r.de.e = 0
	cpu.r.hl.h = 0
	cpu.r.hl.l = 0
}

func (cpu *cpu) Stop() {
	// TODO: to implement
	cpu.Reset()
}

func (cpu *cpu) Step() {
	op := cpu.fetch()
	instr := cpu.decode(op)
	cycles := cpu.execute(instr)
	cpu.cycle += cycles
}

func (cpu *cpu) fetch() byte {
	address := mmu.Address{High: cpu.r.pc.high(), Low: cpu.r.pc.low()}
	opcode := cpu.mmu.ReadByte(address)
	cpu.r.pc++
	return opcode
}

func (cpu *cpu) decode(op byte) instruction {
	instr := operations[op]
	return instr
}

func (cpu *cpu) execute(instr instruction) cycleCount {
	cycles := instr(cpu)
	return cycles
}
