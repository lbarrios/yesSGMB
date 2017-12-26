package cpu

import (
	"github.com/lbarrios/yesSGMB/types"
)

func (cpu *cpu) StepDebug() {
	op := cpu.fetchDebug()
	instr := cpu.decodeDebug(op)
	cycles := cpu.executeDebug(instr)
	cpu.cycle += cycles
	cpu.log.Println(cpu.r)
	cpu.log.Println("")
}

func (cpu *cpu) fetchDebug() byte {
	address := types.Address{High: cpu.r.pc.High(), Low: cpu.r.pc.Low()}
	opcode := cpu.mmu.ReadByte(address)
	cpu.log.Printf("fetch(0x%.2x%.2x) = 0x%.4x.", cpu.r.pc.High(), cpu.r.pc.Low(), opcode)
	cpu.r.pc++
	return opcode
}

func (cpu *cpu) decodeDebug(op byte) instruction {
	instr := operations[op]
	cpu.log.Printf("decode = %s", instructions_debug[op])
	return instr
}

func (cpu *cpu) executeDebug(instr instruction) cycleCount {
	cycles := instr(cpu)
	cpu.log.Printf("execute = %d cycles", cycles)
	return cycles
}
