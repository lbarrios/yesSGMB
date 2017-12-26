// Package gpu implements the Graphics Processor Unit
package gpu

import (
	"github.com/lbarrios/yesSGMB/logger"
	"github.com/lbarrios/yesSGMB/mmu"
	"github.com/lbarrios/yesSGMB/types"
	"sync"
	"time"
)

const (
	DISPLAY_WIDTH  = 160
	DISPLAY_HEIGHT = 144
)

const (
	VIDEO_RAM_START = types.Word(0x8000)
	VIDEO_RAM_END   = types.Word(0x9FFF)
	OAM_START       = types.Word(0xFE00) // Sprite Attribute Table
	OAM_END         = types.Word(0xFE9F) // Sprite Attribute Table
	LCDC_ADDRESS    = types.Word(0xFF40)
	STAT_ADDRESS    = types.Word(0xFF41)
	SCY_ADDRESS     = types.Word(0xFF42)
	SCX_ADDRESS     = types.Word(0xFF43)
	LY_ADDRESS      = types.Word(0xFF44)
	LYC_ADDRESS     = types.Word(0xFF45)
)

type gpu struct {
	irqHandler mmu.IRQHandler
	log        logger.Logger
	lcdc       *byte
	stat       *byte
	scy        *byte
	scx        *byte
	ly         *byte
	lyc        *byte
	video_ram  [1 + VIDEO_RAM_END - VIDEO_RAM_START]*byte
	oam        [1 + OAM_END - OAM_START]*byte
}

func NewGpu(mmu mmu.IRQHandler, l *logger.Logger) *gpu {
	gpu := new(gpu)
	gpu.irqHandler = mmu
	gpu.log = *l
	gpu.log.SetPrefix("\033[0;35mGPU: ")
	return gpu
}

func (gpu *gpu) Reset() {
	gpu.log.Println("GPU reset triggered.")
}

func (gpu *gpu) Run(wg *sync.WaitGroup) {
	gpu.log.Println("GPU started.")
	for {
		if *gpu.lcdc == 0xff {
			gpu.log.Println("lcdc=0xff, stopping GPU.")
			break
		}
		gpu.log.Printf("lcdc=0x%.2x; the GPU will stop when this reaches 0xff..", *gpu.lcdc)
		time.Sleep(1 * time.Second)
	}
	wg.Done()
}

func (gpu *gpu) MapByte(logical_address types.Address, physical_address *byte) {
	addr := logical_address.AsWord()
	switch {
	case addr == LCDC_ADDRESS:
		gpu.lcdc = physical_address
	case addr == STAT_ADDRESS:
		gpu.stat = physical_address
	case addr == SCY_ADDRESS:
		gpu.scy = physical_address
	case addr == SCX_ADDRESS:
		gpu.scx = physical_address
	case addr == LY_ADDRESS:
		gpu.ly = physical_address
	case addr == LYC_ADDRESS:
		gpu.lyc = physical_address
	case addr >= VIDEO_RAM_START && addr <= VIDEO_RAM_END:
		gpu.video_ram[addr-VIDEO_RAM_START] = physical_address
	case addr >= OAM_START && addr <= OAM_END:
		gpu.oam[addr-OAM_START] = physical_address
	default:
		gpu.log.Fatalf("Trying to map unexpected address: 0x%.4x", addr)
	}
}
