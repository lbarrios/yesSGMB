package cpu

import (
	"github.com/lbarrios/yesSGMB/types"
)

func (cpu *cpu) StepDebug() {
	op := cpu.fetchDebug()
	instr := cpu.decodeDebug(op)
	cycles := cpu.executeDebug(instr)
	cpu.clock.Cycles += uint64(cycles)
	cpu.log.Println(cpu.r)
	cpu.log.Println("")
}

func (cpu *cpu) fetchDebug() byte {
	address := types.Address{High: cpu.r.pc.High(), Low: cpu.r.pc.Low()}
	opcode := cpu.mmu.ReadByte(address)
	cpu.log.Printf("%d: fetch(0x%.2x%.2x) = 0x%.4x.", cpu.clock.Cycles, cpu.r.pc.High(), cpu.r.pc.Low(), opcode)
	cpu.r.pc++
	return opcode
}

func (cpu *cpu) decodeDebug(op byte) instruction {
	instr := operations[op]
	cpu.log.Printf("%d: decode = %s", cpu.clock.Cycles, instructions_debug[op])
	return instr
}

func (cpu *cpu) executeDebug(instr instruction) cycleCount {
	cycles := instr(cpu)
	cpu.log.Printf("%d: execute = %d cycles", cpu.clock.Cycles, cycles)
	return cycles
}
