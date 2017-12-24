package main

import (
	"flag"
	"github.com/lbarrios/yesSGMB/logger"
	"github.com/lbarrios/yesSGMB/cartridge"
	"github.com/lbarrios/yesSGMB/mmu"
	"github.com/lbarrios/yesSGMB/cpu"
)

var (
	romFile = flag.String("rom", "test.gb", "Path to rom file")
	log = logger.Logger("MAIN: ")
)

func main() {
	// Parsing the parameters
	flag.Parse()

	// Loading the cartridge data
	cart, err := cartridge.NewCartridge(*romFile)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	mmu := mmu.NewMMU()
	mmu.LoadCartridge(cart)
	cpu := cpu.NewCPU(mmu)
	log.Println("A")
	cpu.Run()
}
