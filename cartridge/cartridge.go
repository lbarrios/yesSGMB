// Package cartridge implements the Cartridge type, that abstracts the GameBoy cartridge
// it basically reads a rom from a file, and loads it in memory
// the format of the header can be found here: http://gbdev.gg8.se/wiki/articles/The_Cartridge_Header
package cartridge

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/lbarrios/yesSGMB/logger"
	"io/ioutil"
	"strings"
)

type Cartridge struct {
	Filename            string
	data                []byte
	nintendoLogo        []byte
	Title               string
	manufacturerCode    string
	colorFlag           byte
	licensee            []byte
	licenseeDescription string
	sgbFlag             byte
	Type                byte
	TypeDescription     string
	romBanks            int
	romSize             int
	ramSize             int
	destinationCode     byte
	isJapanese          bool
	versionNumber       int
	headerChecksum      byte
	globalChecksum      []byte
	MBC                 MemoryBankController
	log                 logger.Logger
}

const (
	nintendoLogoStart          = 0x0104
	nintendoLogoEnd            = 0x0133 + 1
	titleStart                 = 0x0134
	titleEnd                   = 0x0142 + 1
	manufacturerCodeStart      = 0x013F
	manufacturerCodeEnd        = 0x0142 + 1
	colorFlagPosition          = 0x143
	cgbCompatibleColorFlag     = 0x80
	cgbOnlyColorFlag           = 0xC0
	oldLicenseeCodePosition    = 0x014B
	newLicenseeCodeFlag        = 0x33
	newLicenseeCodeStart       = 0x0144
	newLicenseeCodeEnd         = 0x0145 + 1
	sgbFlagPosition            = 0x146
	typePosition               = 0x0147
	romBanksPosition           = 0x0148
	ramSizePosition            = 0x0149
	destinationCodePosition    = 0x014A
	isJapaneseDestinationValue = 0x00
	versionNumberPosition      = 0x014C
	headerChecksumPosition     = 0x014D
	globalChecksumStart        = 0x014E
	globalChecksumEnd          = 0x014F + 1
)

func NewCartridge(filename string, l *logger.Logger) (*Cartridge, error) {
	c := new(Cartridge)
	c.log = *l
	c.log.SetPrefix("\033[0;30mCART: ")

	if err := c.LoadROMFile(filename); err != nil {
		return nil, err
	}

	if err := c.ParseHeader(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Cartridge) ParseHeader() error {
	// As the documentation states, the minimum cartridge ROM size is when the cartridge has zero rom banks
	minimumRomSize := romSizeForBanks(0)
	if len(c.data) < minimumRomSize {
		return errors.New(fmt.Sprintf("The cartridge rom size is lower than the minimum expected (%dKB).", minimumRomSize/1024))
	}

	// Nintendo Logo
	// These bytes define the bitmap of the Nintendo logo that is displayed when the GameBoy gets turned on.
	// The GameBoy boot procedure verifies the content of this bitmap (after it has displayed it), and LOCKS ITSELF UP if these bytes are incorrect.
	// A CGB verifies only the first 18h bytes of the bitmap, but others (for example a Pocket GameBoy) verify all 30h bytes.
	c.nintendoLogo = c.data[nintendoLogoStart:nintendoLogoEnd]
	if !bytes.Equal(c.nintendoLogo, originalNintendoLogo) {
		c.log.Println(c.nintendoLogo)
		c.log.Println(originalNintendoLogo)
		return errors.New("The cartridge is not original! It's a pirate copy, Nintendo is losing money!")
	}

	// Title - The empty chars are filled with 0's
	c.Title = strings.Trim(string(c.data[titleStart:titleEnd]), "\0000")
	c.log.Printf("The game title is: %s", c.Title)

	// GBC - Manufacturer Code
	c.manufacturerCode = string(c.data[manufacturerCodeStart:manufacturerCodeEnd])

	// GBC - Color Flag
	// 		80h - Game supports CGB functions, but works on old gameboys also.
	// 		C0h - Game works on CGB only (physically the same as 80h).
	c.colorFlag = c.data[colorFlagPosition]
	if c.colorFlag == cgbOnlyColorFlag {
		return errors.New("The cartridge is not compatible with GameBoy Classic; requires GameBoy Color")
	}

	// Licensee Code
	// Specifies a two character ASCII licensee code, indicating the company or publisher of the game
	// If the old licensee code equals to 0x33, then the New Licencee Code is used
	if c.data[oldLicenseeCodePosition] != newLicenseeCodeFlag {
		c.licensee = []byte{c.data[oldLicenseeCodePosition]}
	} else {
		c.licensee = c.data[newLicenseeCodeStart:newLicenseeCodeEnd]
	}
	c.licenseeDescription = licenseeMap[fmt.Sprintf("%x", c.licensee)]
	c.log.Printf("The game vendor is: %s=%s", fmt.Sprintf("%x", c.licensee), c.licenseeDescription)

	// SGB - Super GameBoy Flag
	// 		00h = No SGB functions (Normal Gameboy or CGB only game)
	//		03h = Game supports SGB functions
	c.sgbFlag = c.data[sgbFlagPosition]

	// Cartridge Type
	c.Type = c.data[typePosition]
	if desc, ok := typeMap[c.Type]; !ok {
		return errors.New(fmt.Sprintf("Unknown cartridge type: %X", c.Type))
	} else {
		c.TypeDescription = desc
	}
	switch c.Type {
	case MBC_ROMONLY:
		c.MBC = &MBCRomOnly{log: c.log}
	case MBC_1:
		c.MBC = &MBC1{log: c.log}
	default:
		return errors.New(fmt.Sprintf("Unknown cartridge type for MBC: %X", c.Type))
	}
	c.MBC.Init(c.data)
	c.log.Printf("The cartridge type is: %s", c.TypeDescription)

	// ROM Banks and Size
	c.romBanks = romBanksMap[c.data[romBanksPosition]]
	c.romSize = romSizeForBanks(c.romBanks)
	c.log.Printf("ROM size is %d KB", c.romSize/1024)

	// RAM Size
	c.ramSize = ramSizeMap[c.data[ramSizePosition]]
	c.log.Printf("RAM size is %d KB", c.ramSize/1024)

	// Destination Code (ie Japanese or not)
	c.destinationCode = c.data[destinationCodePosition]
	c.isJapanese = c.destinationCode == isJapaneseDestinationValue

	// Game Version Number
	c.versionNumber = int(c.data[versionNumberPosition])

	// Header Checksum
	// Contains an 8 bit checksum across the cartridge header bytes 0134-014C.
	// The checksum is calculated as follows:
	// 		x=0:FOR i=0134h TO 014Ch:x=x-MEM[i]-1:NEXT
	// The lower 8 bits of the result must be the same than the value in this entry.
	// The GAME WON'T WORK if this checksum is incorrect.
	c.headerChecksum = c.data[headerChecksumPosition]

	// Global Checksum
	c.globalChecksum = c.data[globalChecksumStart:globalChecksumEnd]

	return nil
}

func (c *Cartridge) LoadROMFile(filename string) error {
	c.log.Printf("Loading file %s...", filename)
	c.Filename = filename

	if fileContent, err := ioutil.ReadFile(filename); err != nil {
		return err
	} else {
		c.data = fileContent
	}

	c.log.Printf("File %s loaded.", filename)
	return nil
}

func romSizeForBanks(romBanks int) int {
	var romSize int
	if romBanks == 0 {
		romSize = 32 * 1024
	} else {
		romSize = romBanks * 16 * 1024
	}
	return romSize
}

var typeMap = map[byte]string{
	0x00: "ROM ONLY",
	0x01: "MBC1",
	0x02: "MBC1+RAM",
	0x03: "MBC1+RAM+BATTERY",
	0x05: "MBC2",
	0x06: "MBC2+BATTERY",
	0x08: "ROM+RAM",
	0x09: "ROM+RAM+BATTERY",
	0x0B: "MMM01",
	0x0C: "MMM01+RAM",
	0x0D: "MMM01+RAM+BATTERY",
	0x0F: "MBC3+TIMER+BATTERY",
	0x10: "MBC3+TIMER+RAM+BATTERY",
	0x11: "MBC3",
	0x12: "MBC3+RAM",
	0x13: "MBC3+RAM+BATTERY",
	0x19: "MBC5",
	0x1A: "MBC5+RAM",
	0x1B: "MBC5+RAM+BATTERY",
	0x1C: "MBC5+RUMBLE",
	0x1D: "MBC5+RUMBLE+RAM",
	0x1E: "MBC5+RUMBLE+RAM+BATTERY",
	0x20: "MBC6",
	0x22: "MBC7+SENSOR+RUMBLE+RAM+BATTERY",
	0xFC: "POCKET CAMERA",
	0xFD: "BANDAI TAMA5",
	0xFE: "HuC3",
	0xFF: "HuC1+RAM+BATTERY",
}

var licenseeMap = map[string]string{
	"01": "Nintendo",
	"02": "Rocket Games",
	"08": "Capcom",
	"09": "Hot B Co.",
	"0A": "Jaleco",
	"0B": "Coconuts Japan",
	"0C": "Coconuts Japan/G.X.Media",
	"0H": "Starfish",
	"0L": "Warashi Inc.",
	"0N": "Nowpro",
	"0P": "Game Village",
	"13": "Electronic Arts Japan",
	"18": "Hudson Soft Japan",
	"19": "S.C.P.",
	"1A": "Yonoman",
	"1G": "SMDE",
	"1P": "Creatures Inc.",
	"1Q": "TDK Deep Impresion",
	"20": "Destination Software",
	"22": "VR 1 Japan",
	"25": "San-X",
	"28": "Kemco Japan",
	"29": "Seta",
	"2H": "Ubisoft Japan",
	"2K": "NEC InterChannel",
	"2L": "Tam",
	"2M": "Jordan",
	"2N": "Smilesoft",
	"2Q": "Mediakite",
	"36": "Codemasters",
	"37": "GAGA Communications",
	"38": "Laguna",
	"39": "Telstar Fun and Games",
	"41": "Ubi Soft Entertainment",
	"42": "Sunsoft",
	"47": "Spectrum Holobyte",
	"49": "IREM",
	"4D": "Malibu Games",
	"4F": "Eidos/U.S. Gold",
	"4J": "Fox Interactive",
	"4K": "Time Warner Interactive",
	"4Q": "Disney",
	"4S": "Black Pearl",
	"4X": "GT Interactive",
	"4Y": "RARE",
	"4Z": "Crave Entertainment",
	"50": "Absolute Entertainment",
	"51": "Acclaim",
	"52": "Activision",
	"53": "American Sammy Corp.",
	"54": "Take 2 Interactive",
	"55": "Hi Tech",
	"56": "LJN LTD.",
	"58": "Mattel",
	"5A": "Mindscape/Red Orb Ent.",
	"5C": "Taxan",
	"5D": "Midway",
	"5F": "American Softworks",
	"5G": "Majesco Sales Inc",
	"5H": "3DO",
	"5K": "Hasbro",
	"5L": "NewKidCo",
	"5M": "Telegames",
	"5N": "Metro3D",
	"5P": "Vatical Entertainment",
	"5Q": "LEGO Media",
	"5S": "Xicat Interactive",
	"5T": "Cryo Interactive",
	"5W": "Red Storm Ent./BKN Ent.",
	"5X": "Microids",
	"5Z": "Conspiracy Entertainment Corp.",
	"60": "Titus Interactive Studios",
	"61": "Virgin Interactive",
	"62": "Maxis",
	"64": "LucasArts Entertainment",
	"67": "Ocean",
	"69": "Electronic Arts",
	"6E": "Elite Systems Ltd.",
	"6F": "Electro Brain",
	"6G": "The Learning Company",
	"6H": "BBC",
	"6J": "Software 2000",
	"6L": "BAM! Entertainment",
	"6M": "Studio 3",
	"6Q": "Classified Games",
	"6S": "TDK Mediactive",
	"6U": "DreamCatcher",
	"6V": "JoWood Productions",
	"6W": "SEGA",
	"6X": "Wannado Edition",
	"6Y": "LSP",
	"6Z": "ITE Media",
	"70": "Infogrames",
	"71": "Interplay",
	"72": "JVC Musical Industries Inc",
	"73": "Parker Brothers",
	"75": "SCI",
	"78": "THQ",
	"79": "Accolade",
	"7A": "Triffix Ent. Inc.",
	"7C": "Microprose Software",
	"7D": "Universal Interactive Studios",
	"7F": "Kemco",
	"7G": "Rage Software",
	"7H": "Encore",
	"7J": "Zoo",
	"7K": "BVM",
	"7L": "Simon & Schuster Interactive",
	"7M": "Asmik Ace Entertainment Inc./AIA",
	"7N": "Empire Interactive",
	"7Q": "Jester Interactive",
	"7T": "Scholastic",
	"7U": "Ignition Entertainment",
	"7W": "Stadlbauer",
	"80": "Misawa",
	"83": "LOZC",
	"8B": "Bulletproof Software",
	"8C": "Vic Tokai Inc.",
	"8J": "General Entertainment",
	"8N": "Success",
	"8P": "SEGA Japan",
	"91": "Chun Soft",
	"92": "Video System",
	"93": "BEC",
	"96": "Yonezawa/S'pal",
	"97": "Kaneko",
	"99": "Victor Interactive Software",
	"9A": "Nichibutsu/Nihon Bussan",
	"9B": "Tecmo",
	"9C": "Imagineer",
	"9F": "Nova",
	"9H": "Bottom Up",
	"9L": "Hasbro Japan",
	"9N": "Marvelous Entertainment",
	"9P": "Keynet Inc.",
	"9Q": "Hands-On Entertainment",
	"A0": "Telenet",
	"A1": "Hori",
	"A4": "Konami",
	"A6": "Kawada",
	"A7": "Takara",
	"A9": "Technos Japan Corp.",
	"AA": "JVC",
	"AC": "Toei Animation",
	"AD": "Toho",
	"AF": "Namco",
	"AG": "Media Rings Corporation",
	"AH": "J-Wing",
	"AK": "KID",
	"AL": "MediaFactory",
	"AP": "Infogrames Hudson",
	"AQ": "Kiratto. Ludic Inc",
	"B0": "Acclaim Japan",
	"B1": "ASCII",
	"B2": "Bandai",
	"B4": "Enix",
	"B6": "HAL Laboratory",
	"B7": "SNK",
	"B9": "Pony Canyon Hanbai",
	"BA": "Culture Brain",
	"BB": "Sunsoft",
	"BD": "Sony Imagesoft",
	"BF": "Sammy",
	"BG": "Magical",
	"BJ": "Compile",
	"BL": "MTO Inc.",
	"BN": "Sunrise Interactive",
	"BP": "Global A Entertainment",
	"BQ": "Fuuki",
	"C0": "Taito",
	"C2": "Kemco",
	"C3": "Square Soft",
	"C5": "Data East",
	"C6": "Tonkin House",
	"C8": "Koei",
	"CA": "Konami/Palcom/Ultra",
	"CB": "Vapinc/NTVIC",
	"CC": "Use Co.,Ltd.",
	"CD": "Meldac",
	"CE": "FCI/Pony Canyon",
	"CF": "Angel",
	"CM": "Konami Computer Entertainment Osaka",
	"CP": "Enterbrain",
	"D1": "Sofel",
	"D2": "Quest",
	"D3": "Sigma Enterprises",
	"D4": "Ask Kodansa",
	"D6": "Naxat",
	"D7": "Copya System",
	"D9": "Banpresto",
	"DA": "TOMY",
	"DB": "LJN Japan",
	"DD": "NCS",
	"DF": "Altron Corporation",
	"DH": "Gaps Inc.",
	"DN": "ELF",
	"E2": "Yutaka",
	"E3": "Varie",
	"E5": "Epoch",
	"E7": "Athena",
	"E8": "Asmik Ace Entertainment Inc.",
	"E9": "Natsume",
	"EA": "King Records",
	"EB": "Atlus",
	"EC": "Epic/Sony Records",
	"EE": "IGS",
	"EL": "Spike",
	"EM": "Konami Computer Entertainment Tokyo",
	"EN": "Alphadream Corporation",
	"F0": "A Wave",
	"G1": "PCCW",
	"G4": "KiKi Co Ltd",
	"G5": "Open Sesame Inc.",
	"G6": "Sims",
	"G7": "Broccoli",
	"G8": "Avex",
	"G9": "D3 Publisher",
	"GB": "Konami Computer Entertainment Japan",
	"GD": "Square-Enix",
	"HY": "Sachen",
}

var romBanksMap = map[byte]int{
	0x00: 0,
	0x01: 4,
	0x02: 8,
	0x03: 16,
	0x04: 32,
	0x05: 64,
	0x06: 128,
	0x07: 256,
	0x08: 512,
	0x52: 72,
	0x53: 80,
	0x54: 96,
}

var ramSizeMap = map[byte]int{
	0x00: 0,
	0x01: 2 * 1024,
	0x02: 8 * 1024,
	0x03: 32 * 1024,
	0x04: 128 * 1024,
	0x05: 64 * 1024,
}

var originalNintendoLogo = []byte{
	0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B, 0x03, 0x73, 0x00, 0x83, 0x00, 0x0C, 0x00, 0x0D,
	0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E, 0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99,
	0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC, 0xCC, 0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E,
}
