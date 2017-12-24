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

const (
	ROM_BANK_0_16KB             = 0x0000
	SWITCHABLE_ROM_BANK_16KB    = 0x4000
	VIDEO_RAM_8KB               = 0x8000
	SWITCHABLE_RAM_BANK_8KB     = 0xA000
	INTERNAL_RAM_8KB            = 0xC000
	ECHO_8KB_INTERNAL_RAM       = 0xE000
	SPRITE_ATTRIB_MEMORY_OAM    = 0xFE00
	EMPTY_BUT_UNUSABLE_FOR_IO_1 = 0xFEA0
	IO_PORTS                    = 0xFF00
	EMPTY_BUT_UNUSABLE_FOR_IO_2 = 0xFF4C
	INTERNAL_RAM                = 0xFF80
	INTERRUPT_ENABLE_REGISTER   = 0xFFFF
)

func (mmu *mmu) ReadByte(address types.Address) byte {
	if address.AsWord() < SWITCHABLE_ROM_BANK_16KB {
		// ROM_BANK_0_16KB
		return mmu.cartridge.MBC.Read(address)
	}
	if address.AsWord() >= SWITCHABLE_ROM_BANK_16KB && address.AsWord() < VIDEO_RAM_8KB {
		// SWITCHABLE_ROM_BANK_16KB
		return mmu.cartridge.MBC.Read(address)
	}
	if address.AsWord() >= VIDEO_RAM_8KB && address.AsWord() < SWITCHABLE_RAM_BANK_8KB {
		// VIDEO_RAM_8KB
		log.Printf("MMU: Attemping to read from unimplemented address %x (VIDEO_RAM_8KB)", address.AsWord())
	}
	if address.AsWord() >= SWITCHABLE_RAM_BANK_8KB && address.AsWord() < INTERNAL_RAM_8KB {
		// SWITCHABLE_RAM_BANK_8KB
		return mmu.cartridge.MBC.Read(address)
	}
	if address.AsWord() >= INTERNAL_RAM_8KB && address.AsWord() < ECHO_8KB_INTERNAL_RAM {
		// INTERNAL_RAM_8KB
		log.Printf("MMU: Attemping to read from unimplemented address %x (INTERNAL_RAM_8KB)", address.AsWord())
	}
	if address.AsWord() >= ECHO_8KB_INTERNAL_RAM && address.AsWord() < SPRITE_ATTRIB_MEMORY_OAM {
		// ECHO_8KB_INTERNAL_RAM
		log.Printf("MMU: Attemping to read from unimplemented address %x (ECHO_8KB_INTERNAL_RAM)", address.AsWord())
	}
	if address.AsWord() >= SPRITE_ATTRIB_MEMORY_OAM && address.AsWord() < EMPTY_BUT_UNUSABLE_FOR_IO_1 {
		// SPRITE_ATTRIB_MEMORY_OAM
		log.Printf("MMU: Attemping to read from unimplemented address %x (SPRITE_ATTRIB_MEMORY_OAM)", address.AsWord())
	}
	if address.AsWord() >= EMPTY_BUT_UNUSABLE_FOR_IO_1 && address.AsWord() < IO_PORTS {
		// EMPTY_BUT_UNUSABLE_FOR_IO_1
		log.Printf("MMU: Attemping to read from unimplemented address %x (EMPTY_BUT_UNUSABLE_FOR_IO_1)", address.AsWord())
	}
	if address.AsWord() >= IO_PORTS && address.AsWord() < EMPTY_BUT_UNUSABLE_FOR_IO_2 {
		// IO_PORTS
		log.Printf("MMU: Attemping to read from unimplemented address %x (IO_PORTS)", address.AsWord())
	}
	if address.AsWord() >= EMPTY_BUT_UNUSABLE_FOR_IO_2 && address.AsWord() < INTERNAL_RAM {
		// EMPTY_BUT_UNUSABLE_FOR_IO_2
		log.Printf("MMU: Attemping to read from unimplemented address %x (EMPTY_BUT_UNUSABLE_FOR_IO_2)", address.AsWord())
	}
	if address.AsWord() >= INTERNAL_RAM && address.AsWord() < INTERRUPT_ENABLE_REGISTER {
		// INTERNAL_RAM
		log.Printf("MMU: Attemping to read from unimplemented address %x (INTERNAL_RAM)", address.AsWord())
	}
	if address.AsWord() >= INTERRUPT_ENABLE_REGISTER {
		// INTERRUPT_ENABLE_REGISTER
		log.Printf("MMU: Attemping to read from unimplemented address %x (INTERRUPT_ENABLE_REGISTER)", address.AsWord())
	}
	log.Printf("MMU: Attemping to read from invalid address %x", address.AsWord())
	return 0x00
}

func (mmu *mmu) WriteByte(address types.Address, value byte) {
}
