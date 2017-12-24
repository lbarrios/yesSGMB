// package mmu implements the Memory Management Unit of the Gameboy
package mmu

type Address struct {
	High byte
	Low  byte
}

func (address Address) NextAddress() Address {
	addr := address
	addr.Low += 1
	if addr.Low == 0x00 {
		addr.High += 1
	}
	return addr
}

type mmu struct {
	bios [0x100]byte
}

type MMU interface {
	ReadByte(address Address) byte
	WriteByte(address Address, value byte)
}

func NewMMU() *mmu {
	mmu := new(mmu)
	return mmu
}

func (mmu *mmu) ReadByte(address Address) byte {
	return 0x00
}

func (mmu *mmu) WriteByte(address Address, value byte) {
}
