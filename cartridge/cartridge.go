/*
This package implements the Cartridge type, that abstracts the GameBoy cartridge.

It basically reads a rom from a file, and loads it in memory.
*/
package cartridge

import (
	//"fmt"
	"io/ioutil"
	"log"
)

type cartridge struct {
	filename    string
	fileContent []byte
}

func NewCartridge(filename string) (*cartridge, error) {
	c := new(cartridge)

	if err := c.LoadROMFile(filename); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *cartridge) LoadROMFile(filename string) error {
	c.filename = filename
	log.Printf("Loading cartridge %s", filename)

	if fileContent, err := ioutil.ReadFile(filename); err != nil {
		return err
	} else {
		c.fileContent = fileContent
	}

	log.Printf("File %s loaded.", filename)
	log.Printf("data as hex: %x", c.fileContent)

	return nil
}
