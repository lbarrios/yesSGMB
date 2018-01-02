package cpu

import "github.com/lbarrios/yesSGMB/types"

//interrupt handler addresses
const (
	V_BLANK_IR_ADDR        types.Word = 0x40
	LCD_IR_ADDR            types.Word = 0x48
	TIMER_OVERFLOW_IR_ADDR types.Word = 0x50
	JOYP_HILO_IR_ADDR      types.Word = 0x60
)

func (cpu *cpu) jumpToInterruptHandler(address types.Word) {
	cpu.r.sp--
	cpu.mmu.WriteByte(cpu.r.spAsAddress(), address.High())
	cpu.r.sp--
	cpu.mmu.WriteByte(cpu.r.spAsAddress(), address.Low())
}