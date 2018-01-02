// Package gpu implements the Graphics Processor Unit
package gpu

import (
	"github.com/lbarrios/yesSGMB/clock"
	"github.com/lbarrios/yesSGMB/display"
	"github.com/lbarrios/yesSGMB/logger"
	"github.com/lbarrios/yesSGMB/mmu"
	"github.com/lbarrios/yesSGMB/types"
	"sync"
)

const (
	VIDEO_RAM_START types.Word = 0x8000
	VIDEO_RAM_END   types.Word = 0x9FFF
	// Video Ram Detail
	TILEDATA1_START types.Word = 0x8000
	TILEDATA0_START types.Word = 0x8800
	TILEDATA_SIZE   types.Word = 0x1000 // 4KB = 256x256
	TILEMAP0_START  types.Word = 0x9800
	TILEMAP1_START  types.Word = 0x9C00
	TILEMAP_SIZE    types.Word = 0x400 // 32*32 = 1024 tiles
	// Sprite Attribute Table
	OAM_START types.Word = 0xFE00
	OAM_END   types.Word = 0xFE9F
	// Flags
	LCDC_ADDRESS types.Word = 0xFF40
	STAT_ADDRESS types.Word = 0xFF41
	SCY_ADDRESS  types.Word = 0xFF42
	SCX_ADDRESS  types.Word = 0xFF43
	LY_ADDRESS   types.Word = 0xFF44
	LYC_ADDRESS  types.Word = 0xFF45
)

const ( // Video modes
	STAT_MODE_MASK = 0x03 // bit-mask to obtain the mode from the the stat register
	HBLANK_MODE    = 0x00
	VBLANK_MODE    = 0x01
	OAM_MODE       = 0x02
	VRAM_MODE      = 0x03
)

const ( // Video modes cycles
	HBLANK_MODE_CYCLES = 204
	VBLANK_MODE_CYCLES = 4560
	OAM_MODE_CYCLES    = 80
	VRAM_MODE_CYCLES   = 172
)

const ( // Interruptions
	VBLANK_IRQ = 0
)

type gpu struct {
	clock       clock.ClockCounter
	irqHandler  mmu.IRQHandler
	display     *display.Display
	log         logger.Logger
	lcdControl  *byte // lcdc = 0xFF40
	stat        *byte // stat = 0xFF41
	scrollY     *byte // scy = 0xFF42
	scrollX     *byte // scx = 0xFF43
	currentLine *byte // ly = 0xFF44
	lyc         *byte // lyc = 0xFF45
	video_ram   [1 + VIDEO_RAM_END - VIDEO_RAM_START]*byte
	oam         [1 + OAM_END - OAM_START]*byte
	tileMap0    [TILEMAP_SIZE]*byte
	tileMap1    [TILEMAP_SIZE]*byte
	tileData0   [TILEDATA_SIZE]*byte
	tileData1   [TILEDATA_SIZE]*byte

	displayOn    bool
	backgroundOn bool
	windowOn     bool
	spritesOn    bool
}

func NewGpu(mmu mmu.IRQHandler, l *logger.Logger) *gpu {
	gpu := new(gpu)
	gpu.irqHandler = mmu
	gpu.log = *l
	gpu.log.SetPrefix("\033[0;35mGPU: ")
	return gpu
}

func (gpu *gpu) ConnectClock(clockWg *sync.WaitGroup, clock clock.Clock) chan uint64 {
	gpu.clock.Init(clockWg, make(chan uint64), clock)
	return gpu.clock.Channel
}

func (gpu *gpu) ConnectDisplay(d *display.Display) {
	gpu.display = d
}

func (gpu *gpu) GetName() string {
	return "gpu"
}

func (gpu *gpu) MapByte(logical_address types.Address, physical_address *byte) {
	addr := logical_address.AsWord()
	switch {
	case addr == LCDC_ADDRESS:
		gpu.lcdControl = physical_address
	case addr == STAT_ADDRESS:
		gpu.stat = physical_address
	case addr == SCY_ADDRESS:
		gpu.scrollY = physical_address
	case addr == SCX_ADDRESS:
		gpu.scrollX = physical_address
	case addr == LY_ADDRESS:
		gpu.currentLine = physical_address
	case addr == LYC_ADDRESS:
		gpu.lyc = physical_address
	case addr >= VIDEO_RAM_START && addr <= VIDEO_RAM_END:
		gpu.video_ram[addr-VIDEO_RAM_START] = physical_address
		if addr >= TILEDATA1_START && addr < TILEDATA1_START+TILEDATA_SIZE {
			gpu.tileData1[addr-TILEDATA1_START] = physical_address
		}
		if addr >= TILEDATA0_START && addr < TILEDATA0_START+TILEDATA_SIZE {
			gpu.tileData0[addr-TILEDATA0_START] = physical_address
		}
		if addr >= TILEMAP0_START && addr < TILEMAP0_START+TILEMAP_SIZE {
			gpu.tileMap0[addr-TILEMAP0_START] = physical_address
		}
		if addr >= TILEMAP1_START && addr < TILEMAP1_START+TILEMAP_SIZE {
			gpu.tileMap1[addr-TILEMAP1_START] = physical_address
		}
	case addr >= OAM_START && addr <= OAM_END:
		gpu.oam[addr-OAM_START] = physical_address
	default:
		gpu.log.Fatalf("Trying to map unexpected address: 0x%.4x", addr)
	}
}

func (gpu *gpu) Reset() {
	gpu.log.Println("GPU reset triggered.")
}

func (gpu *gpu) mode() byte {
	mode := *gpu.stat
	mode &= STAT_MODE_MASK
	return mode
}

func (gpu *gpu) setMode(mode byte) {
	*gpu.stat = ((*gpu.stat >> 2) << 2) | mode

	switch mode {
	case HBLANK_MODE:

	case OAM_MODE:

	case VBLANK_MODE:

	case VRAM_MODE:

	}
}

func (gpu *gpu) step() {
	gpu.log.Printf("mode: %.2x", gpu.mode())
	switch {
	case gpu.mode() == HBLANK_MODE:
		gpu.log.Println("HBLANK")
		// render the current line
		if gpu.backgroundOn {

		}
		if gpu.windowOn {

		}
		if gpu.spritesOn {

		}

		// go to the next line
		gpu.clock.Cycles += HBLANK_MODE_CYCLES
		*gpu.currentLine += 1
		if *gpu.currentLine < display.HEIGHT {
			gpu.setMode(OAM_MODE)
		} else {
			gpu.setMode(VBLANK_MODE)
		}

	case gpu.mode() == VBLANK_MODE:
		gpu.log.Println("VBLANK")
		gpu.display.Refresh()
		gpu.irqHandler.RequestInterrupt(VBLANK_IRQ)
		gpu.clock.Cycles += VBLANK_MODE_CYCLES
		*gpu.currentLine = 0
		gpu.setMode(OAM_MODE)

	case gpu.mode() == OAM_MODE:
		gpu.clock.Cycles += OAM_MODE_CYCLES
		gpu.setMode(VRAM_MODE)

	case gpu.mode() == VRAM_MODE:
		gpu.clock.Cycles += VRAM_MODE_CYCLES
		gpu.setMode(HBLANK_MODE)
	}
}

func (gpu *gpu) renderBackgroundOnLine() {

}

func (gpu *gpu) renderWindowOnLine() {

}

func (gpu *gpu) renderSpritesOnLine() {

}

func (gpu *gpu) Run(wg *sync.WaitGroup) {
	gpu.log.Println("GPU started.")

	for {
		gpu.clock.WaitNextCycle()

		if *gpu.lcdControl == 0xff {
			gpu.log.Println("lcdc=0xff, stopping GPU.")
			gpu.clock.Disconnect(gpu)
			break
		}

		gpu.step()
	}

	gpu.display.Disconnect(wg)
	wg.Done()
}
