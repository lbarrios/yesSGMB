// Package gpu implements the Graphics Processor Unit
package gpu

import (
	"github.com/lbarrios/yesSGMB/mmu"
	"github.com/lbarrios/yesSGMB/logger"
	"github.com/lbarrios/yesSGMB/types"
	"sync"
	"time"
)

const (
	DISPLAY_WIDTH  = 160
	DISPLAY_HEIGHT = 144
)

type gpu struct {
	irqHandler mmu.IRQHandler
	log        logger.Logger
	lcdc       *byte // 0xFF40
	stat       *byte // 0xFF41
	scy        *byte // 0xFF42
	scx        *byte // 0xFF43
	ly         *byte // 0xFF44
	lyc        *byte // 0xFF45
}

func NewGpu(mmu mmu.IRQHandler, l *logger.Logger) *gpu {
	gpu := new(gpu)
	gpu.irqHandler = mmu
	gpu.log = *l
	gpu.log.SetPrefix("\033[0;35mGPU: ")
	gpu.reset()
	return gpu
}

func (gpu *gpu) reset() {
	gpu.log.Println("GPU reset triggered.")
}

func (gpu *gpu) Run(wg *sync.WaitGroup) {
	gpu.log.Println("GPU started.")
	for {
		if *gpu.lcdc == 0xff {
			gpu.log.Println("lcdc=0xff, stopping GPU.")
			break
		}
		time.Sleep(1 * time.Millisecond)
	}
	wg.Done()
}

func (gpu *gpu) MapByte(logical_address types.Address, physical_address *byte) {
	addr := logical_address.AsWord()
	switch {
	case addr == 0xFF40:
		gpu.log.Println("Mapping GPU lcdc to MMU address 0xFF40")
		gpu.lcdc = physical_address
	default:
		gpu.log.Fatalf("Trying to map unexpected address: $.4x", logical_address)
	}
}
