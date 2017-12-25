package mmu

type IRQHandler interface {
	RequestInterrupt(interrupt byte)
}