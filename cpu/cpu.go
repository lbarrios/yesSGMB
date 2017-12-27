// Package cpu implements the CPU and Registers type
package cpu

import (
	"github.com/lbarrios/yesSGMB/logger"
	"github.com/lbarrios/yesSGMB/mmu"
	"github.com/lbarrios/yesSGMB/types"
	"sync"
	"github.com/lbarrios/yesSGMB/clock"
)

type cpu struct {
	r                 Registers
	mmu               mmu.MMU
	interruptsEnabled bool
	halted            bool
	log               logger.Logger
	clock             clock.ClockCounter
}

func NewCPU(mmu mmu.MMU, l *logger.Logger) *cpu {
	cpu := new(cpu)
	cpu.mmu = mmu
	cpu.log = *l
	cpu.log.SetPrefix("\033[0;34mCPU: ")
	cpu.Reset()
	return cpu
}

func (cpu *cpu) ConnectClock(clockWg *sync.WaitGroup, clock clock.Clock) chan uint64 {
	cpu.clock.Init(clockWg, make(chan uint64), clock)
	return cpu.clock.Channel
}

func (cpu *cpu) GetName() string {
	return "cpu"
}

func (cpu *cpu) Reset() {
	cpu.log.Println("CPU reset triggered.")
	cpu.r.pc = 0x0100 // On power up, the GameBoy Program Counter is initialized to 0x0100
	cpu.r.sp = 0xFFFE // On power up, the GameBoy Stack Pointer is initialized to 0xFFFE

	// After power up, on GameBoy, AF = 0x01B0 = 0000 0001 1011 0000
	cpu.r.af.a = 0x01
	cpu.r.af.f = Flags{true, false, true, true, false, false, false, false}
	cpu.r.bc.b = 0x00
	cpu.r.bc.c = 0x13
	cpu.r.de.d = 0x00
	cpu.r.de.e = 0xD8
	cpu.r.hl.h = 0x01
	cpu.r.hl.l = 0x4D

	// Set the stack default values
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x05}, 0x00) // TIMA
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x06}, 0x00) // TMA
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x07}, 0x00) // TAC
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x10}, 0x80) // NR10
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x11}, 0xBF) // NR11
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x12}, 0xF3) // NR12
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x14}, 0xBF) // NR14
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x16}, 0x3F) // NR21
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x17}, 0x00) // NR22
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x19}, 0xBF) // NR24
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x1A}, 0x7F) // NR30
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x1B}, 0xFF) // NR31
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x1C}, 0x9F) // NR32
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x1E}, 0xBF) // NR33
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x20}, 0xFF) // NR41
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x21}, 0x00) // NR42
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x22}, 0x00) // NR43
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x23}, 0xBF) // NR30
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x24}, 0x77) // NR50
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x25}, 0xF3) // NR51
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x26}, 0xF1) // NR52
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x40}, 0x91) // LCDC
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x42}, 0x00) // SCY
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x43}, 0x00) // SCX
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x45}, 0x00) // LYC
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x47}, 0xFC) // BGP
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x48}, 0xFF) // OBP0
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x49}, 0xFF) // OBP1
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x4A}, 0x00) // WY
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0x4B}, 0x00) // WX
	cpu.mmu.WriteByte(types.Address{High: 0xFF, Low: 0xFF}, 0x00) // IE
}

func (cpu *cpu) Stop() {
	// TODO: to implement
	cpu.Reset()
}

func (cpu *cpu) Step() {
	op := cpu.fetch()
	instr := cpu.decode(op)
	cycles := cpu.execute(instr)
	cpu.clock.Cycles += uint64(cycles)
}

func (cpu *cpu) fetch() byte {
	address := types.Address{High: cpu.r.pc.High(), Low: cpu.r.pc.Low()}
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

func (cpu *cpu) Run(wg *sync.WaitGroup) {
	cpu.log.Println("CPU started.")
	for {
		cpu.clock.WaitNextCycle()

		cpu.StepDebug()
		if cpu.r.pc == 0x02b2 {
			cpu.log.Println("\033[1;31mBreakpoint at 0x02b2.")
			cpu.StepDebug()
			cpu.StepDebug()
			cpu.log.Println("writing 0xff to 0xff40 (this will stop the GPU)")
			cpu.mmu.WriteByte(types.Address{0xff, 0x40}, 0xff)
			cpu.log.Println("writing 0xff to 0xff07 (this will stop the timer)")
			cpu.mmu.WriteByte(types.Address{0xff, 0x07}, 0xff)
			cpu.clock.Disconnect(cpu)
			break
		}
	}
	wg.Done()
}