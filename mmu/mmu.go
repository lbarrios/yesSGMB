// package mmu implements the Memory Management Unit of the Gameboy
package mmu

import (
	"log"
	"github.com/lbarrios/yesSGMB/cartridge"
	"github.com/lbarrios/yesSGMB/types"
)

type mmu struct {
	bios      [0x100]byte
	cartridge *cartridge.Cartridge
}

type MMU interface {
	ReadByte(address types.Address) byte
	WriteByte(address types.Address, value byte)
}

func NewMMU() *mmu {
	mmu := new(mmu)
	return mmu
}

func (mmu *mmu) LoadCartridge(cart *cartridge.Cartridge) {
	mmu.cartridge = cart
}

func (mmu *mmu) ReadByte(address types.Address) byte {
	if address.High < 0x80 {
		if address.High < 0x40 {
			// ROM Bank 0
			return mmu.cartridge.MBC.Read(address)
		} else {
			// ROM Bank 1 (switchable)
			return mmu.cartridge.MBC.Read(address)
		}
	} else {
		if address.High >= 0xA0 && address.High < 0xC0 {
			return mmu.cartridge.MBC.Read(address)
		}
		if address.High >= 0xC0 {
			if address.High < 0xE0 {
				// GB Internal RAM
			}
			if address.High >= 0xE0 && address.High < 0xFE {
				// GB Internal RAM shadow
			}
		}
	}
	log.Printf("MMU: Attemping to read from invalid address %x", address.AsWord())
	return 0x00
}

func (mmu *mmu) WriteByte(address types.Address, value byte) {
}
