//package gpu implements the Graphics Processor Unit
package gpu

import (
	"github.com/lbarrios/yesSGMB/mmu"
	"github.com/lbarrios/yesSGMB/logger"
	"sync"
)

const (
	DISPLAY_WIDTH  = 160
	DISPLAY_HEIGHT = 144
)

type gpu struct {
	irqHandler mmu.IRQHandler
	log        logger.Logger
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
	gpu.irqHandler.RequestInterrupt(0xFF)
}

func (gpu *gpu) Run(wg *sync.WaitGroup) {
	gpu.log.Println("GPU started.")
	wg.Done()
}
