// Package mmu implements the Memory Management Unit of the Gameboy
package mmu

import (
	"github.com/lbarrios/yesSGMB/cartridge"
	"github.com/lbarrios/yesSGMB/logger"
	"github.com/lbarrios/yesSGMB/types"
	"sync"
)

const (
	MIN_ADDRESS               = 0x0000
	MAX_ADDRESS               = 0xFFFF
	INTERRUPT_FLAG_ADDR       = types.Word(0xFF0F)
	INTERRUPT_ENABLE_REGISTER = types.Word(0xFFFF)
)

type mmu struct {
	bios       [0x100]byte
	cartridge  *cartridge.Cartridge
	memory     [MAX_ADDRESS + 1]byte
	memoryLock sync.Mutex
	log        logger.Logger
}

type MMU interface {
	ReadByte(address types.Address) byte
	WriteByte(address types.Address, value byte)
}

func NewMMU(l *logger.Logger) *mmu {
	mmu := new(mmu)
	mmu.log = *l
	mmu.log.SetPrefix("\033[0;33mMMU: ")
	return mmu
}

func (mmu *mmu) LoadCartridge(cart *cartridge.Cartridge) {
	mmu.cartridge = cart
}

const (
	//ROM_BANK_0_16KB             = 0x0000
	SWITCHABLE_ROM_BANK_16KB    = 0x4000
	VIDEO_RAM_8KB               = 0x8000
	SWITCHABLE_RAM_BANK_8KB     = 0xA000
	INTERNAL_RAM_8KB            = 0xC000
	ECHO_8KB_INTERNAL_RAM       = 0xE000
	SPRITE_ATTRIB_MEMORY_OAM    = 0xFE00
	EMPTY_BUT_UNUSABLE_FOR_IO_1 = 0xFEA0
	IO_PORTS                    = 0xFF00
	EMPTY_BUT_UNUSABLE_FOR_IO_2 = 0xFF4C
	HIGH_RAM                    = 0xFF80
)

func (mmu *mmu) ReadByte(address types.Address) byte {
	mmu.memoryLock.Lock()
	var ret byte

	switch {
	case address.AsWord() < SWITCHABLE_ROM_BANK_16KB:
		// ROM_BANK_0_16KB
		ret = mmu.cartridge.MBC.Read(address)

	case address.AsWord() >= SWITCHABLE_ROM_BANK_16KB && address.AsWord() < VIDEO_RAM_8KB:
		// SWITCHABLE_ROM_BANK_16KB
		ret = mmu.cartridge.MBC.Read(address)

	case address.AsWord() >= VIDEO_RAM_8KB && address.AsWord() < SWITCHABLE_RAM_BANK_8KB:
		// VIDEO_RAM_8KB
		mmu.log.Fatalf("Attemping to read from unimplemented address %x (VIDEO_RAM_8KB)", address.AsWord())

	case address.AsWord() >= SWITCHABLE_RAM_BANK_8KB && address.AsWord() < INTERNAL_RAM_8KB:
		// SWITCHABLE_RAM_BANK_8KB
		ret = mmu.cartridge.MBC.Read(address)

	case address.AsWord() >= INTERNAL_RAM_8KB && address.AsWord() < ECHO_8KB_INTERNAL_RAM:
		// INTERNAL_RAM_8KB
		ret = mmu.memory[address.AsWord()]

	case address.AsWord() >= ECHO_8KB_INTERNAL_RAM && address.AsWord() < SPRITE_ATTRIB_MEMORY_OAM:
		// ECHO_8KB_INTERNAL_RAM
		// The addresses from E000 to FE00 appear to access the internal RAM the same as C000-DE00.
		ret = mmu.memory[address.AsWord()-(ECHO_8KB_INTERNAL_RAM-INTERNAL_RAM_8KB)]

	case address.AsWord() >= SPRITE_ATTRIB_MEMORY_OAM && address.AsWord() < EMPTY_BUT_UNUSABLE_FOR_IO_1:
		// SPRITE_ATTRIB_MEMORY_OAM
		mmu.log.Fatalf("Attemping to read from unimplemented address %x (SPRITE_ATTRIB_MEMORY_OAM)", address.AsWord())

	case address.AsWord() >= EMPTY_BUT_UNUSABLE_FOR_IO_1 && address.AsWord() < IO_PORTS:
		// EMPTY_BUT_UNUSABLE_FOR_IO_1
		mmu.log.Fatalf("Attemping to read from unimplemented address %x (EMPTY_BUT_UNUSABLE_FOR_IO_1)", address.AsWord())

	case address.AsWord() >= IO_PORTS && address.AsWord() < EMPTY_BUT_UNUSABLE_FOR_IO_2:
		// IO_PORTS, this case write to memory that is mapped to peripherals
		ret = mmu.memory[address.AsWord()]

	case address.AsWord() >= EMPTY_BUT_UNUSABLE_FOR_IO_2 && address.AsWord() < HIGH_RAM:
		// EMPTY_BUT_UNUSABLE_FOR_IO_2
		mmu.log.Fatalf("Attemping to read from unimplemented address %x (EMPTY_BUT_UNUSABLE_FOR_IO_2)", address.AsWord())

	case address.AsWord() >= HIGH_RAM && address.AsWord() < INTERRUPT_ENABLE_REGISTER:
		// HIGH_RAM
		ret = mmu.memory[address.AsWord()]

	case address.AsWord() >= INTERRUPT_ENABLE_REGISTER:
		// INTERRUPT_ENABLE_REGISTER
		ret = mmu.memory[address.AsWord()]

	default:
		mmu.log.Fatalf("Attemping to read from invalid address %x", address.AsWord())
		ret = mmu.memory[address.AsWord()]
	}

	mmu.memoryLock.Unlock()
	return ret
}

func (mmu *mmu) WriteByte(address types.Address, value byte) {
	mmu.memoryLock.Lock()

	switch {
	case address.AsWord() < SWITCHABLE_ROM_BANK_16KB:
		// ROM_BANK_0_16KB
		mmu.cartridge.MBC.Write(address, value)

	case address.AsWord() >= SWITCHABLE_ROM_BANK_16KB && address.AsWord() < VIDEO_RAM_8KB:
		// SWITCHABLE_ROM_BANK_16KB
		mmu.cartridge.MBC.Write(address, value)

	case address.AsWord() >= VIDEO_RAM_8KB && address.AsWord() < SWITCHABLE_RAM_BANK_8KB:
		// VIDEO_RAM_8KB
		mmu.log.Fatalf("Attemping to write to unimplemented address %x (VIDEO_RAM_8KB)", address.AsWord())

	case address.AsWord() >= SWITCHABLE_RAM_BANK_8KB && address.AsWord() < INTERNAL_RAM_8KB:
		// SWITCHABLE_RAM_BANK_8KB
		mmu.cartridge.MBC.Write(address, value)

	case address.AsWord() >= INTERNAL_RAM_8KB && address.AsWord() < ECHO_8KB_INTERNAL_RAM:
		// INTERNAL_RAM_8KB
		mmu.memory[address.AsWord()] = value

	case address.AsWord() >= ECHO_8KB_INTERNAL_RAM && address.AsWord() < SPRITE_ATTRIB_MEMORY_OAM:
		// ECHO_8KB_INTERNAL_RAM
		// The addresses from E000 to FE00 appear to access the internal RAM the same as C000-DE00.
		mmu.memory[address.AsWord()-(ECHO_8KB_INTERNAL_RAM-INTERNAL_RAM_8KB)] = value

	case address.AsWord() >= SPRITE_ATTRIB_MEMORY_OAM && address.AsWord() < EMPTY_BUT_UNUSABLE_FOR_IO_1:
		// SPRITE_ATTRIB_MEMORY_OAM
		mmu.log.Fatalf("Attemping to write to unimplemented address %x (SPRITE_ATTRIB_MEMORY_OAM)", address.AsWord())

	case address.AsWord() >= EMPTY_BUT_UNUSABLE_FOR_IO_1 && address.AsWord() < IO_PORTS:
		// EMPTY_BUT_UNUSABLE_FOR_IO_1
		mmu.log.Fatalf("Attemping to write to unimplemented address %x (EMPTY_BUT_UNUSABLE_FOR_IO_1)", address.AsWord())

	case address.AsWord() >= IO_PORTS && address.AsWord() < EMPTY_BUT_UNUSABLE_FOR_IO_2:
		// IO_PORTS, this case write to memory that is mapped to peripherals
		mmu.memory[address.AsWord()] = value

	case address.AsWord() >= EMPTY_BUT_UNUSABLE_FOR_IO_2 && address.AsWord() < HIGH_RAM:
		// EMPTY_BUT_UNUSABLE_FOR_IO_2
		mmu.log.Fatalf("Attemping to write to unimplemented address %x (EMPTY_BUT_UNUSABLE_FOR_IO_2)", address.AsWord())

	case address.AsWord() >= HIGH_RAM && address.AsWord() < INTERRUPT_ENABLE_REGISTER:
		// INTERNAL_RAM
		mmu.memory[address.AsWord()] = value

	case address.AsWord() >= INTERRUPT_ENABLE_REGISTER:
		// INTERRUPT_ENABLE_REGISTER
		mmu.memory[address.AsWord()] = value
		mmu.log.Printf("Writing value %.2x to INTERRUPT_ENABLE_REGISTER", value)

	default:
		mmu.log.Fatalf("Attemping to write to invalid address %x", address.AsWord())
		mmu.memory[address.AsWord()] = value
	}

	mmu.memoryLock.Unlock()
}

func (mmu *mmu) RequestInterrupt(interrupt byte) {
	iflag := mmu.ReadByte(INTERRUPT_FLAG_ADDR.AsAddress())
	iflag |= interrupt
	mmu.WriteByte(INTERRUPT_FLAG_ADDR.AsAddress(), iflag)
	mmu.log.Printf("Interruption %x", interrupt)
}

func (mmu *mmu) MapMemoryRegion(p Peripheral, begin types.Address, end types.Address) {
	if end.AsWord() < begin.AsWord() {
		mmu.log.Fatalf("MapMemoryRegion expects a non-negative lenght interval")
	}
	if begin.AsWord() < MIN_ADDRESS || end.AsWord() > MAX_ADDRESS {
		mmu.log.Fatalf("MapMemoryRegion parameters are out-o-range")
	}
	for i := begin.AsWord(); i < end.AsWord(); i++ {
		mmu.MapMemoryAdress(p, i.AsAddress())
	}
}

func (mmu *mmu) MapMemoryAdress(p Peripheral, address types.Address) {
	mmu.memoryLock.Lock()
	p.MapByte(address, &mmu.memory[address.AsWord()])
	mmu.memoryLock.Unlock()
}
