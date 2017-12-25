package main

import (
	"flag"
	"sync"
	"github.com/lbarrios/yesSGMB/logger"
	"github.com/lbarrios/yesSGMB/cartridge"
	"github.com/lbarrios/yesSGMB/mmu"
	"github.com/lbarrios/yesSGMB/cpu"
	"github.com/lbarrios/yesSGMB/gpu"
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

	MMU := mmu.NewMMU(log)
	MMU.LoadCartridge(cart)

	CPU := cpu.NewCPU(MMU, log)
	wg.Add(1)
	go CPU.Run(&wg)

	GPU := gpu.NewGpu(MMU, log)
	wg.Add(1)
	go GPU.Run(&wg)

	wg.Wait()
}
