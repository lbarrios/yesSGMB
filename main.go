package main

import (
	"flag"
	"log"
	"github.com/lbarrios/yesSGMB/cartridge"
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
		log.Fatalf("ERROR: %s",err)
	}
}
