package main

import (
	"flag"
	"github.com/lbarrios/yesSGMB/cartridge"
	"github.com/lbarrios/yesSGMB/cpu"
	"log"
)

var (
	romFile = flag.String("rom", "test.gb", "Path to rom file")
)

func main() {
	// Parsing the parameters
	flag.Parse()

	// Loading the cartridge data
	_, err := cartridge.NewCartridge(*romFile)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	cpu.NewCPU()
}
