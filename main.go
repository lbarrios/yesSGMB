package main

import (
	"flag"
	"log"
	"github.com/lbarrios/yesSGMB/cartridge"
	"fmt"
)

func main() {
	// Parsing the parameters
	romFile := flag.String("rom", "test.gb", "Path to rom file")
	flag.Parse()

	// Loading the cartridge data
	_, err := cartridge.NewCartridge(*romFile)
	if err != nil {
		log.Fatal(fmt.Sprintf("ERROR: %s",err))
	}
}
