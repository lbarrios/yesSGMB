// package mmu implements the Memory Management Unit of the Gameboy
package mmu

type address struct {
	high byte
	low  byte
}

type mmu struct {
	bios [0x100]byte
}

type MMU interface {
	ReadByte(address address) byte
	WriteByte(address address, value byte)
}

func NewMMU() *mmu {
	mmu := new(mmu)
	return mmu
}

func (mmu *mmu) ReadByte(address address) byte {
	return 0x00
}

func (mmu *mmu) WriteByte(address address, value byte) {
}