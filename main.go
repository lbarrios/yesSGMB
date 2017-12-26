package main

import (
	"flag"
	"github.com/lbarrios/yesSGMB/cartridge"
	"github.com/lbarrios/yesSGMB/cpu"
	"github.com/lbarrios/yesSGMB/gpu"
	"github.com/lbarrios/yesSGMB/logger"
	"github.com/lbarrios/yesSGMB/mmu"
	"github.com/lbarrios/yesSGMB/timer"
	"sync"
)

var (
	romFile = flag.String("rom", "test.gb", "Path to rom file")
	log     = new(logger.Logger)
	wg      sync.WaitGroup
)

func main() {
	// Parsing the parameters
	flag.Parse()

	// Initialize the logging
	log.Init()

	// Loading the cartridge data
	cart, err := cartridge.NewCartridge(*romFile, log)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	// Initialize the Memory Management Unit
	MMU := mmu.NewMMU(log)
	MMU.LoadCartridge(cart)

	// Initialize the Central Processing Unit
	CPU := cpu.NewCPU(MMU, log)

	// Initialize the Graphics Processing Unit
	GPU := gpu.NewGpu(MMU, log)
	MMU.MapMemoryRegion(GPU, gpu.VIDEO_RAM_START.AsAddress(), gpu.VIDEO_RAM_END.AsAddress())
	MMU.MapMemoryRegion(GPU, gpu.OAM_START.AsAddress(), gpu.OAM_END.AsAddress())
	MMU.MapMemoryAdress(GPU, gpu.LCDC_ADDRESS.AsAddress())
	MMU.MapMemoryAdress(GPU, gpu.STAT_ADDRESS.AsAddress())
	MMU.MapMemoryAdress(GPU, gpu.SCY_ADDRESS.AsAddress())
	MMU.MapMemoryAdress(GPU, gpu.SCX_ADDRESS.AsAddress())
	MMU.MapMemoryAdress(GPU, gpu.LY_ADDRESS.AsAddress())
	MMU.MapMemoryAdress(GPU, gpu.LYC_ADDRESS.AsAddress())
	GPU.Reset()

	// Initialize the Timer
	Timer := timer.NewTimer(log)
	MMU.MapMemoryAdress(Timer, timer.DIV_ADDRESS.AsAddress())
	MMU.MapMemoryAdress(Timer, timer.TIMA_ADDRESS.AsAddress())
	MMU.MapMemoryAdress(Timer, timer.TMA_ADDRESS.AsAddress())
	MMU.MapMemoryAdress(Timer, timer.TAC_ADDRESS.AsAddress())
	Timer.Reset()

	// Initialize the Clock
	Clock := timer.Clock{Log: log}
	Clock.ConnectPeripheral(CPU)

	// Run all the components
	wg.Add(3)
	go CPU.Run(&wg)
	go GPU.Run(&wg)
	go Timer.Run(&wg)
	go Clock.Run(&wg)

	// Wait to exit the program
	wg.Wait()
}
