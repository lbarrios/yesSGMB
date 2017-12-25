package cartridge

import (
	"github.com/lbarrios/yesSGMB/types"
	"github.com/lbarrios/yesSGMB/logger"
)

const (
	MBC_ROMONLY                     = 0x00
	MBC_1                           = 0x01
	MBC_1_RAM                       = 0x02
	MBC_1_RAM_BATTERY               = 0x03
	MBC_2                           = 0x05
	MBC_2_BATTERY                   = 0x06
	MBC_ROM_RAM                     = 0x08
	MBC_ROM_RAM_BATTERY             = 0x09
	MMM_01                          = 0x0B
	MMM_01_RAM                      = 0x0C
	MMM_01_RAM_BATTERY              = 0x0D
	MBC_3_TIMER_BATTERY             = 0x0F
	MBC_3_TIMER_RAM_BATTERY         = 0x10
	MBC_3                           = 0x11
	MBC_3_RAM                       = 0x12
	MBC_3_RAM_BATTERY               = 0x13
	MBC_5                           = 0x19
	MBC_5_RAM                       = 0x1A
	MBC_5_RAM_BATTERY               = 0x1B
	MBC_5_RUMBLE                    = 0x1C
	MBC_5_RUMBLE_RAM                = 0x1D
	MBC_5_RUMBLE_RAM_BATTERY        = 0x1E
	MBC_6                           = 0x20
	MBC_7_SENSOR_RUMBLE_RAM_BATTERY = 0x22
	POCKET_CAMERA                   = 0xFC
	BANDAI_TAMA5                    = 0xFD
	HuC3                            = 0xFE
	HuC1_RAM_BATTERY                = 0xFF
)

type MemoryBankController interface {
	Init([]byte)
	Write(address types.Address, value byte)
	Read(address types.Address) byte
}

type MBCRomOnly struct {
	romBank []byte
	log     logger.Logger
}

func (mbc *MBCRomOnly) Init(data []byte) {
	mbc.romBank = data[0x0000:0x8000]
}

func (mbc *MBCRomOnly) Write(address types.Address, value byte) {
	if address.High >= 0x80 {
		mbc.log.Fatalf("Cannot write to address: %.4x!", address.AsWord())
		return
	}
	// TODO: Check this
	mbc.log.Printf("Cannot write to address: %.4x!", address.AsWord())
}

func (mbc *MBCRomOnly) Read(address types.Address) byte {
	return mbc.romBank[address.AsWord()]
}
