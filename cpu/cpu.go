// package cpu implements the CPU and Registers type
// the comments are based on http://marc.rawer.de/Gameboy/Docs/GBCPUman.pdf
package cpu

import (
	"log"
)

type Flags struct {
	z bool // zero flag
	n bool // subtract flag
	h bool // half carry flag
	c bool // carry flag
	_ bool // to complete 1 byte
	_ bool // to complete 1 byte
	_ bool // to complete 1 byte
	_ bool // to complete 1 byte
}

type cycleCount int
type instructions func(*cpu) cycleCount

type cpu struct {
	r Registers
}

func NewCPU() *cpu {
	cpu := new(cpu)
	cpu.Reset()
	return cpu
}

func (cpu *cpu) Reset() {
	log.Println("CPU reset triggered.")
	cpu.r.pc = 0x0100 // On power up, the GameBoy Program Counter is initialized to 0x0100
	cpu.r.sp = 0xFFFE // On power up, the GameBoy Stack Pointer is initialized to 0xFFFE
	cpu.r.af.a = 0
	cpu.r.af.f = Flags{false, false, false, false, false, false, false, false}
	cpu.r.bc.b = 0
	cpu.r.bc.c = 0
	cpu.r.de.d = 0
	cpu.r.de.e = 0
	cpu.r.hl.h = 0
	cpu.r.hl.l = 0
}

const (
	nopCycles         = 4  // 0x00
	ldBcNnCycles      = 12 // 0x01
	xZeroTwoCycles    = 0  // 0x02
	incBcCycles       = 8  // 0x03
	incBCycles        = 4  // 0x04
	decBCycles        = 4  // 0x05
	ldBNCycles        = 8  // 0x06
	rlcACycles        = 4  // 0x07
	ldMemNnSpCycles   = 20 // 0x08
	addHlBcCycles     = 8  // 0x09
	ldAMemBcCycles    = 8  // 0x0A
	decBcCycles       = 8  // 0x0B
	incCCycles        = 4  // 0x0C
	decCCycles        = 4  // 0x0D
	ldCNCycles        = 8  // 0x0E
	rrcACycles        = 4  // 0x0F
	stopCycles        = 4  // 0x10
	ldDeNnCycles      = 12 // 0x11
	xOneTwoCycles     = 0  // 0x12
	incDeCycles       = 8  // 0x13
	incDCycles        = 4  // 0x14
	decDCycles        = 4  // 0x15
	ldDNCycles        = 8  // 0x16
	rlACycles         = 4  // 0x17
	jrCycles          = 8  // 0x18
	addHlDeCycles     = 8  // 0x19
	ldAMemDeCycles    = 8  // 0x1A
	decDeCycles       = 8  // 0x1B
	incECycles        = 4  // 0x1C
	decECycles        = 4  // 0x1D
	ldENCycles        = 8  // 0x1E
	rrACycles         = 4  // 0x1F
	jrNZCycles        = 8  // 0x20
	ldHlNnCycles      = 12 // 0x21
	ldiMemHlACycles   = 8  // 0x22
	incHlCycles       = 8  // 0x23
	incHCycles        = 4  // 0x24
	decHCycles        = 4  // 0x25
	ldHNCycles        = 8  // 0x26
	daACycles         = 4  // 0x27
	jrZCycles         = 8  // 0x28
	addHlHlCycles     = 8  // 0x29
	ldiAMemHlCycles   = 8  // 0x2A
	decHlCycles       = 8  // 0x2B
	incLCycles        = 4  // 0x2C
	decLCycles        = 4  // 0x2D
	ldLNCycles        = 8  // 0x2E
	cplACycles        = 4  // 0x2F
	jrNCCycles        = 8  // 0x30
	ldSpNnCycles      = 12 // 0x31
	lddMemHlACycles   = 8  // 0x32
	incSpCycles       = 8  // 0x33
	incMemHlCycles    = 12 // 0x34
	decMemHlCycles    = 12 // 0x35
	ldMemHlNCycles    = 12 // 0x36
	scfCycles         = 4  // 0x37
	jrCCycles         = 8  // 0x38
	addHlSpCycles     = 8  // 0x39
	lddAMemHlCycles   = 8  // 0x3A
	decSpCycles       = 8  // 0x3B
	incACycles        = 4  // 0x3C
	decACycles        = 4  // 0x3D
	ldANCycles        = 0  // 0x3E
	ccfCycles         = 4  // 0x3F
	ldBBCycles        = 4  // 0x40
	ldBCCycles        = 4  // 0x41
	ldBDCycles        = 4  // 0x42
	ldBECycles        = 4  // 0x43
	ldBHCycles        = 4  // 0x44
	ldBLCycles        = 4  // 0x45
	ldBHlCycles       = 8  // 0x46
	ldBACycles        = 4  // 0x47
	ldCBCycles        = 4  // 0x48
	ldCCCycles        = 4  // 0x49
	ldCDCycles        = 4  // 0x4A
	ldCECycles        = 4  // 0x4B
	ldCHCycles        = 4  // 0x4C
	ldCLCycles        = 4  // 0x4D
	ldCHlCycles       = 8  // 0x4E
	ldCACycles        = 4  // 0x4F
	ldDBCycles        = 4  // 0x50
	ldDCCycles        = 4  // 0x51
	ldDDCycles        = 4  // 0x52
	ldDECycles        = 4  // 0x53
	ldDHCycles        = 4  // 0x54
	ldDLCycles        = 4  // 0x55
	ldDHlCycles       = 8  // 0x56
	ldDACycles        = 4  // 0x57
	ldEBCycles        = 4  // 0x58
	ldECCycles        = 4  // 0x59
	ldEDCycles        = 4  // 0x5A
	ldEECycles        = 4  // 0x5B
	ldEHCycles        = 4  // 0x5C
	ldELCycles        = 4  // 0x5D
	ldEHlCycles       = 8  // 0x5E
	ldEACycles        = 4  // 0x5F
	ldHBCycles        = 4  // 0x60
	ldHCCycles        = 4  // 0x61
	ldHDCycles        = 4  // 0x62
	ldHECycles        = 4  // 0x63
	ldHHCycles        = 4  // 0x64
	ldHLCycles        = 4  // 0x65
	ldHHlCycles       = 8  // 0x66
	xSixSevenCycles   = 0  // 0x67
	ldLBCycles        = 4  // 0x68
	ldLCCycles        = 4  // 0x69
	ldLDCycles        = 4  // 0x6A
	ldLECycles        = 4  // 0x6B
	ldLHCycles        = 4  // 0x6C
	ldLLCycles        = 4  // 0x6D
	ldLHlCycles       = 8  // 0x6E
	xSixFCycles       = 0  // 0x6F
	ldMemHlBCycles    = 8  // 0x70
	ldMemHlCCycles    = 8  // 0x71
	ldMemHlDCycles    = 8  // 0x72
	ldMemHlECycles    = 8  // 0x73
	ldMemHlHCycles    = 8  // 0x74
	ldMemHlLCycles    = 8  // 0x75
	haltCycles        = 4  // 0x76
	xSevenSevenCycles = 0  // 0x77
	ldABCycles        = 4  // 0x78
	ldACCycles        = 4  // 0x79
	ldADCycles        = 4  // 0x7A
	ldAECycles        = 4  // 0x7B
	ldAHCycles        = 4  // 0x7C
	ldALCycles        = 4  // 0x7D
	ldAMemHlCycles    = 8  // 0x7E
	ldAACycles        = 4  // 0x7F
	addABCycles       = 4  // 0x80
	addACCycles       = 4  // 0x81
	addADCycles       = 4  // 0x82
	addAECycles       = 4  // 0x83
	addAHCycles       = 4  // 0x84
	addALCycles       = 4  // 0x85
	addAMemHlCycles   = 8  // 0x86
	addAACycles       = 4  // 0x87
	xEightEightCycles = 0  // 0x88
	xEightNineCycles  = 0  // 0x89
	xEightACycles     = 0  // 0x8A
	xEightBCycles     = 0  // 0x8B
	xEightCCycles     = 0  // 0x8C
	xEightDCycles     = 0  // 0x8D
	xEightECycles     = 0  // 0x8E
	xEightFCycles     = 0  // 0x8F
	subABCycles       = 4  // 0x90
	subACCycles       = 4  // 0x91
	subADCycles       = 4  // 0x92
	subAECycles       = 4  // 0x93
	subAHCycles       = 4  // 0x94
	subALCycles       = 4  // 0x95
	subAMemHlCycles   = 8  // 0x96
	subAACycles       = 4  // 0x97
	sbcABCycles       = 4  // 0x98
	sbcACCycles       = 4  // 0x99
	sbcADCycles       = 4  // 0x9A
	sbcAECycles       = 4  // 0x9B
	sbcAHCycles       = 4  // 0x9C
	sbcALCycles       = 4  // 0x9D
	sbcAMemHlCycles   = 8  // 0x9E
	sbcAACycles       = 4  // 0x9F
	andABCycles       = 4  // 0xA0
	andACCycles       = 4  // 0xA1
	andADCycles       = 4  // 0xA2
	andAECycles       = 4  // 0xA3
	andAHCycles       = 4  // 0xA4
	andALCycles       = 4  // 0xA5
	andAMemHlCycles   = 8  // 0xA6
	andAACycles       = 8  // 0xA7
	xorABCycles       = 4  // 0xA8
	xorACCycles       = 4  // 0xA9
	xorADCycles       = 4  // 0xAA
	xorAECycles       = 4  // 0xAB
	xorAHCycles       = 4  // 0xAC
	xorALCycles       = 4  // 0xAD
	xorAMemHlCycles   = 8  // 0xAE
	xorAACycles       = 4  // 0xAF
	orABCycles        = 4  // 0xB0
	orACCycles        = 4  // 0xB1
	orADCycles        = 4  // 0xB2
	orAECycles        = 4  // 0xB3
	orAHCycles        = 4  // 0xB4
	orALCycles        = 4  // 0xB5
	orAMemHlCycles    = 8  // 0xB6
	orAACycles        = 4  // 0xB7
	cpABCycles        = 4  // 0xB8
	cpACCycles        = 4  // 0xB9
	cpADCycles        = 4  // 0xBA
	cpAECycles        = 4  // 0xBB
	cpAHCycles        = 4  // 0xBC
	cpALCycles        = 4  // 0xBD
	cpAMemHlCycles    = 8  // 0xBE
	cpAACycles        = 4  // 0xBF
	xCZeroCycles      = 0  // 0xC0
	popBcCycles       = 12 // 0xC1
	jpNZCycles        = 12 // 0xC2
	jpCycles          = 8  // 0xC3
	callNZCycles      = 12 // 0xC4
	pushBcCycles      = 16 // 0xC5
	addANnCycles      = 8  // 0xC6
	rst00HCycles      = 32 // 0xC7
	xCEightCycles     = 0  // 0xC8
	xCNineCycles      = 0  // 0xC9
	jpZCycles         = 12 // 0xCA
	rxNCycles         = 8  // 0xCB
	callZCycles       = 12 // 0xCC
	callCycles        = 12 // 0xCD
	xCECycles         = 0  // 0xCE
	rst08HCycles      = 32 // 0xCF
	xDZeroCycles      = 0  // 0xD0
	popDeCycles       = 12 // 0xD1
	jpNCCycles        = 12 // 0xD2
	xDThreeCycles     = 0  // 0xD3
	callNCCycles      = 12 // 0xD4
	pushDeCycles      = 16 // 0xD5
	subANnCycles      = 8  // 0xD6
	rst10HCycles      = 32 // 0xD7
	xDEightCycles     = 0  // 0xD8
	xDNineCycles      = 0  // 0xD9
	jpCCycles         = 12 // 0xDA
	xDBCycles         = 0  // 0xDB
	callCCycles       = 12 // 0xDC
	xDDCycles         = 0  // 0xDD
	xDECycles         = 0  // 0xDE
	rst18HCycles      = 32 // 0xDF
	ldStackNACycles   = 12 // 0xE0
	popHlCycles       = 12 // 0xE1
	ldStackCACycles   = 8  // 0xE2
	xEThreeCycles     = 0  // 0xE3
	xEFourCycles      = 0  // 0xE4
	pushHlCycles      = 16 // 0xE5
	andANCycles       = 8  // 0xE6
	rst20HCycles      = 32 // 0xE7
	addSpNCycles      = 16 // 0xE8
	jpMemHlCycles     = 4  // 0xE9
	xEACycles         = 0  // 0xEA
	xEBCycles         = 0  // 0xEB
	xECCycles         = 0  // 0xEC
	xEDCycles         = 0  // 0xED
	xorANCycles       = 8  // 0xEE
	rst28HCycles      = 32 // 0xEF
	ldAStackNCycles   = 12 // 0xF0
	popAfCycles       = 12 // 0xF1
	ldAStackCCycles   = 8  // 0xF2
	diCycles          = 4  // 0xF3
	xFFourCycles      = 0  // 0xF4
	pushAfCycles      = 16 // 0xF5
	orANCycles        = 8  // 0xF6
	rst30HCycles      = 32 // 0xF7
	ldHlSpNCycles     = 12 // 0xF8
	ldSpHlCycles      = 0  // 0xF9
	ldAMemNnCycles    = 16 // 0xFA
	eiCycles          = 4  // 0xFB
	xFCCycles         = 0  // 0xFC
	xFDCycles         = 0  // 0xFD
	cpANCycles        = 8  // 0xFE
	rst38HCycles      = 32 // 0xFF
)

var op = [0x100] instructions{
	nop,       //0x00
	ldBcNn,    //0x01
	TODO,      //0x02
	incBc,     //0x03
	incB,      //0x04
	decB,      //0x05
	ldBN,      //0x06
	rlcA,      //0x07
	ldMemNnSp, //0x08
	addHlBc,   //0x09
	ldAMemBc,  //0x0A
	decBc,     //0x0B
	incC,      //0x0C
	decC,      //0x0D
	ldCN,      //0x0E
	rrcA,      //0x0F
	stop,      //0x10
	ldDeNn,    //0x11
	TODO,      //0x12
	incDe,     //0x13
	incD,      //0x14
	decD,      //0x15
	ldDN,      //0x16
	rlA,       //0x17
	jr,        //0x18
	addHlDe,   //0x19
	ldAMemDe,  //0x1A
	decDe,     //0x1B
	incE,      //0x1C
	decE,      //0x1D
	ldEN,      //0x1E
	rrA,       //0x1F
	jrNZ,      //0x20
	ldHlNn,    //0x21
	ldiMemHlA, //0x22
	incHl,     //0x23
	incH,      //0x24
	decH,      //0x25
	ldHN,      //0x26
	daA,       //0x27
	jrZ,       //0x28
	addHlHl,   //0x29
	ldiAMemHl, //0x2A
	decHl,     //0x2B
	incL,      //0x2C
	decL,      //0x2D
	ldLN,      //0x2E
	cplA,      //0x2F
	jrNC,      //0x30
	ldSpNn,    //0x31
	lddMemHlA, //0x32
	incSp,     //0x33
	incMemHl,  //0x34
	decMemHl,  //0x35
	ldMemHlN,  //0x36
	scf,       //0x37
	jrC,       //0x38
	addHlSp,   //0x39
	lddAMemHl, //0x3A
	decSp,     //0x3B
	incA,      //0x3C
	decA,      //0x3D
	ldAN,      //0x3E
	ccf,       //0x3F
	ldBB,      //0x40
	ldBC,      //0x41
	ldBD,      //0x42
	ldBE,      //0x43
	ldBH,      //0x44
	ldBL,      //0x45
	ldBHl,     //0x46
	ldBA,      //0x47
	ldCB,      //0x48
	ldCC,      //0x49
	ldCD,      //0x4A
	ldCE,      //0x4B
	ldCH,      //0x4C
	ldCL,      //0x4D
	ldCHl,     //0x4E
	ldCA,      //0x4F
	ldDB,      //0x50
	ldDC,      //0x51
	ldDD,      //0x52
	ldDE,      //0x53
	ldDH,      //0x54
	ldDL,      //0x55
	ldDHl,     //0x56
	ldDA,      //0x57
	ldEB,      //0x58
	ldEC,      //0x59
	ldED,      //0x5A
	ldEE,      //0x5B
	ldEH,      //0x5C
	ldEL,      //0x5D
	ldEHl,     //0x5E
	ldEA,      //0x5F
	ldHB,      //0x60
	ldHC,      //0x61
	ldHD,      //0x62
	ldHE,      //0x63
	ldHH,      //0x64
	ldHL,      //0x65
	ldHHl,     //0x66
	TODO,      //0x67
	ldLB,      //0x68
	ldLC,      //0x69
	ldLD,      //0x6A
	ldLE,      //0x6B
	ldLH,      //0x6C
	ldLL,      //0x6D
	ldLHl,     //0x6E
	TODO,      //0x6F
	ldMemHlB,  //0x70
	ldMemHlC,  //0x71
	ldMemHlD,  //0x72
	ldMemHlE,  //0x73
	ldMemHlL,  //0x74
	ldMemHlH,  //0x75
	halt,      //0x76
	TODO,      //0x77
	ldAB,      //0x78
	ldAC,      //0x79
	ldAD,      //0x7A
	ldAE,      //0x7B
	ldAH,      //0x7C
	ldAL,      //0x7D
	ldAMemHl,  //0x7E
	ldAA,      //0x7F
	addAB,     //0x80
	addAC,     //0x81
	addAD,     //0x82
	addAE,     //0x83
	addAH,     //0x84
	addAL,     //0x85
	addAMemHl, //0x86
	TODO,      //0x87
	TODO,      //0x88
	TODO,      //0x89
	TODO,      //0x8A
	TODO,      //0x8B
	TODO,      //0x8C
	TODO,      //0x8D
	TODO,      //0x8E
	TODO,      //0x8F
	subAB,     //0x90
	subAC,     //0x91
	subAD,     //0x92
	subAE,     //0x93
	subAH,     //0x94
	subAL,     //0x95
	subAMemHl, //0x96
	subAA,     //0x97
	sbcAB,     //0x98
	sbcAC,     //0x99
	sbcAD,     //0x9A
	sbcAE,     //0x9B
	sbcAH,     //0x9C
	sbcAL,     //0x9D
	sbcAMemHl, //0x9E
	TODO,      //0x9F
	andAB,     //0xA0
	andAC,     //0xA1
	andAD,     //0xA2
	andAE,     //0xA3
	andAH,     //0xA4
	andAL,     //0xA5
	andAMemHl, //0xA6
	andAA,     //0xA7
	xorAB,     //0xA8
	xorAC,     //0xA9
	xorAD,     //0xAA
	xorAE,     //0xAB
	xorAH,     //0xAC
	xorAL,     //0xAD
	xorAMemHl, //0xAE
	xorAA,     //0xAF
	orAB,      //0xB0
	orAC,      //0xB1
	orAD,      //0xB2
	orAE,      //0xB3
	orAH,      //0xB4
	orAL,      //0xB5
	orAMemHl,  //0xB6
	orAN,      //0xB7
	cpAB,      //0xB8
	cpAC,      //0xB9
	cpAD,      //0xBA
	cpAE,      //0xBB
	cpAH,      //0xBC
	cpAL,      //0xBD
	cpAMemHl,  //0xBE
	cpAN,      //0xBF
	TODO,      //0xC0
	popBc,     //0xC1
	jpNZ,      //0xC2
	jp,        //0xC3
	callNZ,    //0xC4
	pushBc,    //0xC5
	addANn,    //0xC6
	rst00H,    //0xC7
	TODO,      //0xC8
	TODO,      //0xC9
	jpZ,       //0xCA
	rxN,       //0xCB
	callZ,     //0xCC
	call,      //0xCD
	TODO,      //0xCE
	rst08H,    //0xCF
	TODO,      //0xD0
	popDe,     //0xD1
	jpNC,      //0xD2
	TODO,      //0xD3
	callNC,    //0xD4
	pushDe,    //0xD5
	subANn,    //0xD6
	rst10H,    //0xD7
	TODO,      //0xD8
	TODO,      //0xD9
	jpC,       //0xDA
	TODO,      //0xDB
	callC,     //0xDC
	TODO,      //0xDD
	TODO,      //0xDE
	rst18H,    //0xDF
	ldStackNA, //0xE0
	popHl,     //0xE1
	ldStackCA, //0xE2
	TODO,      //0xE3
	TODO,      //0xE4
	pushHl,    //0xE5
	andAN,     //0xE6
	rst20H,    //0xE7
	addSpN,    //0xE8
	jpMemHl,   //0xE9
	TODO,      //0xEA
	TODO,      //0xEB
	TODO,      //0xEC
	TODO,      //0xED
	xorAN,     //0xEE
	rst28H,    //0xEF
	ldAStackN, //0xF0
	popAf,     //0xF1
	ldAStackC, //0xF2
	di,        //0xF3
	TODO,      //0xF4
	pushAf,    //0xF5
	orAN,      //0xF6
	rst30H,    //0xF7
	ldHlSpN,   //0xF8
	ldSpHl,    //0xF9
	ldAMemNn,  //0xFA
	ei,        //0xFB
	TODO,      //0xFC
	TODO,      //0xFD
	cpAN,      //0xFE
	rst38H,    //0xFF
}

func TODO(cpu *cpu) cycleCount {
	// This function is not defined!
	return 0
}

// 3.3. Instructions
// The GameBoy CPU is based on a subset of the Z80 micro-
// processor. A summary of these commands is given below.
// If 'Flags affected' is not given for a command then
// none are affected.

// 3.3.1. 8-Bit Loads

// 3.3.1.1. LD r1,n
// Description:
// 		Put value n into r1.
// Use with:
// 		r1 = B,C,D,E,H,L
// 		n = 8 bit immediate value

func ldBN(cpu *cpu) cycleCount {
	// Put value of the immediate value n into register B
	// TODO: to implement
	return ldBNCycles
}

func ldCN(cpu *cpu) cycleCount {
	// Put value of the immediate value n into register C
	// TODO: to implement
	return ldCNCycles
}

func ldDN(cpu *cpu) cycleCount {
	// Put value of the immediate value n into register D
	// TODO: to implement
	return ldDNCycles
}

func ldEN(cpu *cpu) cycleCount {
	// Put value of the immediate value n into register E
	// TODO: to implement
	return ldENCycles
}

func ldHN(cpu *cpu) cycleCount {
	// Put value of the immediate value n into register H
	// TODO: to implement
	return ldHNCycles
}

func ldLN(cpu *cpu) cycleCount {
	// Put value of the immediate value n into register L
	// TODO: to implement
	return ldLNCycles
}

// 3.3.1.2. LD r1,r2
// Description:
// 		Put value r2 into r1.
// Use with:
// 		r1 = A,B,C,D,E,H,L,(HL)
//		r2 = A,B,C,D,E,H,L,(HL)

func ldAA(cpu *cpu) cycleCount {
	// Put value of register A into register A
	cpu.r.af.a = cpu.r.af.a
	return ldAACycles
}

func ldAB(cpu *cpu) cycleCount {
	// Put value of register B into register A
	cpu.r.af.a = cpu.r.bc.b
	return ldABCycles
}

func ldAC(cpu *cpu) cycleCount {
	// Put value of register C into register A
	cpu.r.af.a = cpu.r.bc.c
	return ldACCycles
}

func ldAD(cpu *cpu) cycleCount {
	// Put value of register D into register A
	cpu.r.af.a = cpu.r.de.d
	return ldADCycles
}

func ldAE(cpu *cpu) cycleCount {
	// Put value of register E into register A
	cpu.r.af.a = cpu.r.de.e
	return ldAECycles
}

func ldAH(cpu *cpu) cycleCount {
	// Put value of register H into register A
	cpu.r.af.a = cpu.r.hl.h
	return ldAHCycles
}

func ldAL(cpu *cpu) cycleCount {
	// Put value of register L into register A
	cpu.r.af.a = cpu.r.hl.l
	return ldALCycles
}

func ldAHl(cpu *cpu) cycleCount {
	// Put value of the position of memory indicated by register HL into register A
	// TODO: how to implement this?
	//cpu.r.af.a = cpu.r.hl
	return ldAMemHlCycles
}

func ldBB(cpu *cpu) cycleCount {
	// Put value of register B into register B
	cpu.r.bc.b = cpu.r.bc.b
	return ldBBCycles
}

func ldBC(cpu *cpu) cycleCount {
	// Put value of register C into register B
	cpu.r.bc.b = cpu.r.bc.c
	return ldBCCycles
}

func ldBD(cpu *cpu) cycleCount {
	// Put value of register D into register B
	cpu.r.bc.b = cpu.r.de.d
	return ldBDCycles
}

func ldBE(cpu *cpu) cycleCount {
	// Put value of register E into register B
	cpu.r.bc.b = cpu.r.de.e
	return ldBECycles
}

func ldBH(cpu *cpu) cycleCount {
	// Put value of register H into register B
	cpu.r.bc.b = cpu.r.hl.h
	return ldBHCycles
}

func ldBL(cpu *cpu) cycleCount {
	// Put value of register L into register B
	cpu.r.bc.b = cpu.r.hl.l
	return ldBLCycles
}

func ldBHl(cpu *cpu) cycleCount {
	// Put value of the position of memory indicated by register HL into register B
	// TODO: how to implement this?
	return ldBHlCycles
}

func ldCB(cpu *cpu) cycleCount {
	// Put value of register B into register C
	cpu.r.bc.c = cpu.r.bc.b
	return ldCBCycles
}
func ldCC(cpu *cpu) cycleCount {
	// Put value of register C into register C
	cpu.r.bc.c = cpu.r.bc.c
	return ldCCCycles
}
func ldCD(cpu *cpu) cycleCount {
	// Put value of register D into register C
	cpu.r.bc.c = cpu.r.de.d
	return ldCDCycles
}
func ldCE(cpu *cpu) cycleCount {
	// Put value of register E into register C
	cpu.r.bc.c = cpu.r.de.e
	return ldCECycles
}
func ldCH(cpu *cpu) cycleCount {
	// Put value of register H into register C
	cpu.r.bc.c = cpu.r.hl.h
	return ldCHCycles
}
func ldCL(cpu *cpu) cycleCount {
	// Put value of register L into register C
	cpu.r.bc.c = cpu.r.hl.l
	return ldCLCycles
}
func ldCHl(cpu *cpu) cycleCount {
	// Put value of the position of memory indicated by register HL into register C
	// TODO: how to implement this?
	return ldCHlCycles
}
func ldDB(cpu *cpu) cycleCount {
	// Put value of register B into register D
	cpu.r.de.d = cpu.r.bc.b
	return ldDBCycles
}
func ldDC(cpu *cpu) cycleCount {
	// Put value of register C into register D
	cpu.r.de.d = cpu.r.bc.c
	return ldDCCycles
}
func ldDD(cpu *cpu) cycleCount {
	// Put value of register D into register D
	cpu.r.de.d = cpu.r.de.d
	return ldDDCycles
}
func ldDE(cpu *cpu) cycleCount {
	// Put value of register E into register D
	cpu.r.de.d = cpu.r.de.e
	return ldDECycles
}
func ldDH(cpu *cpu) cycleCount {
	// Put value of register H into register D
	cpu.r.de.d = cpu.r.hl.h
	return ldDHCycles
}
func ldDL(cpu *cpu) cycleCount {
	// Put value of register L into register D
	cpu.r.de.d = cpu.r.hl.l
	return ldDLCycles
}
func ldDHl(cpu *cpu) cycleCount {
	// Put value of the position of memory indicated by register HL into register D
	// TODO: how to implement this?
	return ldDHlCycles
}
func ldEB(cpu *cpu) cycleCount {
	// Put value of register B into register E
	cpu.r.de.e = cpu.r.bc.b
	return ldEBCycles
}
func ldEC(cpu *cpu) cycleCount {
	// Put value of register C into register E
	cpu.r.de.e = cpu.r.bc.c
	return ldECCycles
}
func ldED(cpu *cpu) cycleCount {
	// Put value of register D into register E
	cpu.r.de.e = cpu.r.de.d
	return ldEDCycles
}
func ldEE(cpu *cpu) cycleCount {
	// Put value of register E into register E
	cpu.r.de.e = cpu.r.de.e
	return ldEECycles
}
func ldEH(cpu *cpu) cycleCount {
	// Put value of register H into register E
	cpu.r.de.e = cpu.r.hl.h
	return ldEHCycles
}
func ldEL(cpu *cpu) cycleCount {
	// Put value of register L into register E
	cpu.r.de.e = cpu.r.hl.l
	return ldELCycles
}
func ldEHl(cpu *cpu) cycleCount {
	// Put value of the position of memory indicated by register HL into register E
	// TODO: to implement
	return ldEHlCycles
}
func ldHB(cpu *cpu) cycleCount {
	// Put value of register B into register H
	cpu.r.hl.h = cpu.r.bc.b
	return ldHBCycles
}
func ldHC(cpu *cpu) cycleCount {
	// Put value of register C into register H
	cpu.r.hl.h = cpu.r.bc.c
	return ldHCCycles
}
func ldHD(cpu *cpu) cycleCount {
	// Put value of register D into register H
	cpu.r.hl.h = cpu.r.de.d
	return ldHDCycles
}
func ldHE(cpu *cpu) cycleCount {
	// Put value of register H into register H
	cpu.r.hl.h = cpu.r.de.e
	return ldHECycles
}
func ldHH(cpu *cpu) cycleCount {
	// Put value of register H into register H
	cpu.r.hl.h = cpu.r.hl.h
	return ldHHCycles
}
func ldHL(cpu *cpu) cycleCount {
	// Put value of register L into register H
	cpu.r.hl.h = cpu.r.hl.l
	return ldHLCycles
}
func ldHHl(cpu *cpu) cycleCount {
	// Put value of the position of memory indicated by register HL into register H
	// TODO: to implement
	return ldHHlCycles
}
func ldLB(cpu *cpu) cycleCount {
	// Put value of register B into register L
	cpu.r.hl.l = cpu.r.bc.b
	return ldLBCycles
}
func ldLC(cpu *cpu) cycleCount {
	// Put value of register C into register L
	cpu.r.hl.l = cpu.r.bc.c
	return ldLCCycles
}
func ldLD(cpu *cpu) cycleCount {
	// Put value of register D into register L
	cpu.r.hl.l = cpu.r.de.d
	return ldLDCycles
}
func ldLE(cpu *cpu) cycleCount {
	// Put value of register E into register L
	cpu.r.hl.l = cpu.r.de.e
	return ldLECycles
}
func ldLH(cpu *cpu) cycleCount {
	// Put value of register H into register L
	cpu.r.hl.l = cpu.r.hl.h
	return ldLHCycles
}
func ldLL(cpu *cpu) cycleCount {
	// Put value of register L into register L
	cpu.r.hl.l = cpu.r.hl.l
	return ldLLCycles
}
func ldLHl(cpu *cpu) cycleCount {
	// Put value of the position of memory indicated by register HL into register L
	// TODO: to implement
	return ldLHlCycles
}
func ldMemHlB(cpu *cpu) cycleCount {
	// Put value of register B into the position of memory indicated by register HL
	// TODO: to implement
	return ldMemHlBCycles
}
func ldMemHlC(cpu *cpu) cycleCount {
	// Put value of register C into the position of memory indicated by register HL
	// TODO: to implement
	return ldMemHlCCycles
}
func ldMemHlD(cpu *cpu) cycleCount {
	// Put value of register D into the position of memory indicated by register HL
	// TODO: to implement
	return ldMemHlDCycles
}
func ldMemHlE(cpu *cpu) cycleCount {
	// Put value of register E into the position of memory indicated by register HL
	// TODO: to implement
	return ldMemHlECycles
}
func ldMemHlL(cpu *cpu) cycleCount {
	// Put value of register L into the position of memory indicated by register HL
	// TODO: to implement
	return ldMemHlLCycles
}
func ldMemHlH(cpu *cpu) cycleCount {
	// Put value of register H into the position of memory indicated by register HL
	// TODO: to implement
	return ldMemHlHCycles
}

func ldMemHlN(cpu *cpu) cycleCount {
	// Put the immediate value n into the position of memory indicated by register HL
	// TODO: to implement
	return ldMemHlNCycles
}

// 3.3.1.3. LD A,n
// Description:
// 		Put value n into A.
// Use with:
// 		n = A,B,C,D,E,H,L,(BC),(DE),(HL),(nn),#
// 		nn = two byte immediate value. (LS byte first.)

func ldAMemBc(cpu *cpu) cycleCount {
	// Put the value into the position of memory indicated by register BC into register A
	// TODO: to implement
	return ldAMemBcCycles
}

func ldAMemDe(cpu *cpu) cycleCount {
	// Put the value into the position of memory indicated by register BC into register A
	// TODO: to implement
	return ldAMemDeCycles
}

func ldAMemHl(cpu *cpu) cycleCount {
	// Put the value into the position of memory indicated by register BC into register A
	// TODO: to implement
	return ldAMemHlCycles
}

func ldAMemNn(cpu *cpu) cycleCount {
	// Put the value into the position of memory indicated by immediate value NN into register A
	// TODO: to implement
	return ldAMemNnCycles
}

func ldAN(cpu *cpu) cycleCount {
	// Put the immediate value N into register A
	// TODO: to implement
	return ldANCycles
}

// 3.3.1.4. LD n,A
// Description:
// 		Put value A into n.
// Use with:
// 		n = A,B,C,D,E,H,L,(BC),(DE),(HL),(nn)
// 		nn = two byte immediate value. (LS byte first.)

func ldBA(cpu *cpu) cycleCount {
	// Put the value of register A into register B
	cpu.r.bc.b = cpu.r.af.a
	return ldBACycles
}

func ldCA(cpu *cpu) cycleCount {
	// Put the value of register A into register B
	cpu.r.bc.c = cpu.r.af.a
	return ldCACycles
}
func ldDA(cpu *cpu) cycleCount {
	// Put the value of register A into register B
	cpu.r.de.d = cpu.r.af.a
	return ldDACycles
}
func ldEA(cpu *cpu) cycleCount {
	// Put the value of register A into register B
	cpu.r.de.e = cpu.r.af.a
	return ldEACycles
}

// 3.3.1.5. LD A,(C)
// Description:
// 		Put value at address $FF00 + register C into A.
// Same as: LD A,($FF00+C)

func ldAStackC(cpu *cpu) cycleCount {
	// Put the value from the position of memory (0xFF00+BC) into register A
	// TODO: to implement
	return ldAStackCCycles
}

// 3.3.1.6. LD (C),A
// Description:
//		Put A into address $FF00 + register C.

func ldStackCA(cpu *cpu) cycleCount {
	// Put the value from the register A into the position of memory (0xFF00+BC)
	// TODO: to implement
	return ldStackCACycles
}

// 3.3.1.7. LD A,(HLD)
// Description: Same as: LDD A,(HL)

// 3.3.1.8. LD A,(HL-)
// Description: Same as: LDD A,(HL)

// 3.3.1.9. LDD A,(HL)
// Description:
// 		Put value at address HL into A. Decrement HL.
// Same as: LD A,(HL) - DEC HL

func lddAMemHl(cpu *cpu) cycleCount {
	// Put the value from the position of memory in HL into the register A.
	// Then, decrement HL.
	// TODO: to implement
	return lddAMemHlCycles
}

// 3.3.1.10. LD (HLD), A
// Description: Same as: LDD (HL),A

// 3.3.1.11. LD (HLD), A
// Description: Same as: LDD (HL),A

// 3.3.1.12. LDD (HL), A
// Description:
// 		Put A into memory address HL. Decrement HL.
// Same as: LD (HL),A - DEC HL

func lddMemHlA(cpu *cpu) cycleCount {
	// Put value of the register A into the memory address pointed by HL.
	// Then, decrement HL.
	// TODO: to implement
	return lddMemHlACycles
}

// 3.3.1.13. LD A,(HLI)
// Description: Same as: LDI A,(HL)

// 3.3.1.14. LD A,(HL+)
// Description: Same as: LDI A,(HL)

// 3.3.1.15. LDI A,(HL)
// Description:
//		Put value at address HL into A. Increment HL.
// Same as: LD A,(HL) - INC HL

func ldiAMemHl(cpu *cpu) cycleCount {
	// Put the value from the position of memory pointed by HL into the register A.
	// Then, increment  HL.
	// TODO: to implement
	return ldiAMemHlCycles
}

// 3.3.1.16. LD (HLI),A
// Description: Same as: LDI (HL),A

// 3.3.1.17. LD (HL+),A
// Description: Same as: LDI (HL),A

// 3.3.1.18. LDI (HL),A
// Description:
//		Put A into memory address HL. Increment HL.
// Same as: LD (HL),A - INC HL

func ldiMemHlA(cpu *cpu) cycleCount {
	// Put value of the register A into the memory address pointed by HL.
	// Then, increment HL.
	// TODO: to implement
	return ldiMemHlACycles
}

// 3.3.1.19. LDH (n),A
// Description:
//		Put A into memory address $FF00+n.
// Use with:
//		n = one byte immediate value.
func ldStackNA(cpu *cpu) cycleCount {
	// Takes from the stack the value indexed by the immediate value N, and put it into register A.
	//TODO: to implement
	return ldStackNACycles
}

// 3.3.1.20. LDH A,(n)
// Description:
//		Put memory address $FF00+n into A.
// Use with:
//		n = one byte immediate value.
func ldAStackN(cpu *cpu) cycleCount {
	// Takes the value from the register A and put it into the stack the value indexed by the immediate value N.
	//TODO: to implement
	return ldAStackNCycles
}

// 3.3.2. 16-Bit Loads

// 3.3.2.1. LD n,nn
// Description:
// 		Put value nn into n.
// Use with:
// 		n = BC,DE,HL,SP
// 		nn = 16 bit immediate value

func ldBcNn(cpu *cpu) cycleCount {
	// Takes a 16-bit immediate value and put it into the register BC.
	//TODO: to implement
	return ldBcNnCycles
}

func ldDeNn(cpu *cpu) cycleCount {
	// Takes a 16-bit immediate value and put it into the register DE.
	//TODO: to implement
	return ldDeNnCycles
}

func ldHlNn(cpu *cpu) cycleCount {
	// Takes a 16-bit immediate value and put it into the register HL.
	//TODO: to implement
	return ldHlNnCycles
}

func ldSpNn(cpu *cpu) cycleCount {
	// Takes a 16-bit immediate value and put it into the register SP.
	//TODO: to implement
	return ldSpNnCycles
}

// 3.3.2.2. LD SP,HL
// Description:
// 		Put HL into Stack Pointer (SP).

func ldSpHl(cpu *cpu) cycleCount {
	// Put the value of the register HL into SP.
	cpu.r.sp = cpu.r.hlAsInt()
	return ldSpHlCycles
}

// 3.3.2.3. LD HL,SP+n
// Description: Same as: LDHL SP,n.

// 3.3.2.4. LDHL SP,n
// Description:
// 		Put SP + n effective address into HL.
// Use with:
// 		n = one byte signed immediate value.
// Flags affected:
// 		Z - Reset.
// 		N - Reset.
// 		H - Set or reset according to operation.
// 		C - Set or reset according to operation.

func ldHlSpN(cpu *cpu) cycleCount {
	// Put the value addressed by Stack Pointer (SP) + N, into register HL
	// cpu.r.hl = ??
	// TODO: to implement
	return ldHlSpNCycles
}

// 3.3.2.5. LD (nn),SP
// Description:
// 		Put Stack Pointer (SP) at address n.
// Use with:
// 		nn = two byte immediate address.
// Opcodes:
// 		Instruction 	Parameters 	Opcode 	Cycles
// 		LD 				(nn),SP 	08 		20

func ldMemNnSp(cpu *cpu) cycleCount {
	// Put the value of Stack Pointer into the memory position addressed by immediate value nn
	// TODO: to implement
	return ldMemNnSpCycles
}

// 3.3.2.6. PUSH nn
// Description:
// 		Push register pair nn onto stack. Decrement Stack Pointer (SP) twice.
// Use with:
// 		nn = AF,BC,DE,HL
// Opcodes:
// 		Instruction		Parameters 	Opcode  	Cycles
// 		PUSH			AF			F5        	16
// 		PUSH			BC			C5        	16
// 		PUSH			DE			D5        	16
// 		PUSH			HL			E5        	16

func pushAf(cpu *cpu) cycleCount {
	// Put the value of register AF into the stack.
	// Then, decrement SP twice
	// TODO: to implement
	return pushAfCycles
}

func pushBc(cpu *cpu) cycleCount {
	// Put the value of register BC into the stack.
	// Then, decrement SP twice
	// TODO: to implement
	return pushBcCycles
}

func pushDe(cpu *cpu) cycleCount {
	// Put the value of register DE into the stack.
	// Then, decrement SP twice
	// TODO: to implement
	return pushDeCycles
}

func pushHl(cpu *cpu) cycleCount {
	// Put the value of register HL into the stack.
	// Then, decrement SP twice
	// TODO: to implement
	return pushHlCycles
}

// 3.3.2.7. POP nn
// Description:
// Pop two bytes off stack into register pair nn. Increment Stack Pointer (SP) twice.
// Use with:
// nn = AF,BC,DE,HL
// Opcodes:
// 		Instruction 	Parameters 		Opcode		Cycles
// 		POP				AF 				F1			12
// 		POP 			BC 				C1			12
// 		POP 			DE 				D1			12
// 		POP 			HL				E1			12

func popAf(cpu *cpu) cycleCount {
	// Take two bytes from the stack into register AF
	// Then, increment SP twice
	// TODO: to implement
	return popAfCycles
}

func popBc(cpu *cpu) cycleCount {
	// Take two bytes from the stack into register BC
	// Then, increment SP twice
	// TODO: to implement
	return popBcCycles
}

func popDe(cpu *cpu) cycleCount {
	// Take two bytes from the stack into register DE
	// Then, increment SP twice
	// TODO: to implement
	return popDeCycles
}

func popHl(cpu *cpu) cycleCount {
	// Take two bytes from the stack into register HL
	// Then, increment SP twice
	// TODO: to implement
	return popHlCycles
}

// 3.3.3. 8-Bit ALU

// 3.3.3.1. ADD A,n
// Description:
// 		Add n to A.
// Use with:
// 		n = A,B,C,D,E,H,L,(HL),#
// Flags affected:
// 		Z - Set if result is zero.
// 		N - Reset.
// 		H - Set if carry from bit 3.
// 		C - Set if carry from bit 7.
// Opcodes:
// 		Instruction 	Parameters 		Opcode		Cycles
//		ADD				A, A			87			4
// 		ADD				A, B			80			4
// 		ADD				A, C			81			4
// 		ADD				A, D			82			4
// 		ADD				A, E			83			4
// 		ADD				A, H			84			4
// 		ADD				A, L			85			4
// 		ADD				A, (HL)			86			8
// 		ADD				A, #			C6			8

func addAA(cpu *cpu) cycleCount {
	// Add the value of register A into register A
	var oldA = cpu.r.af.a
	cpu.r.af.a += cpu.r.af.a
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH((cpu.r.af.a & 0xf) < (oldA & 0xf))
	cpu.r.setFlagC((cpu.r.af.a < oldA))
	return addAACycles
}

func addAB(cpu *cpu) cycleCount {
	// Add the value of register B into register A
	var oldA = cpu.r.af.a
	cpu.r.af.a += cpu.r.bc.b
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH((cpu.r.af.a & 0xf) < (oldA & 0xf))
	cpu.r.setFlagC((cpu.r.af.a < oldA))
	return addABCycles
}

func addAC(cpu *cpu) cycleCount {
	// Add the value of register C into register A
	var oldA = cpu.r.af.a
	cpu.r.af.a += cpu.r.bc.c
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH((cpu.r.af.a & 0xf) < (oldA & 0xf))
	cpu.r.setFlagC((cpu.r.af.a < oldA))
	return addACCycles
}

func addAD(cpu *cpu) cycleCount {
	// Add the value of register D into register A
	var oldA = cpu.r.af.a
	cpu.r.af.a += cpu.r.de.d
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH((cpu.r.af.a & 0xf) < (oldA & 0xf))
	cpu.r.setFlagC((cpu.r.af.a < oldA))
	return addADCycles
}

func addAE(cpu *cpu) cycleCount {
	// Add the value of register E into register A
	var oldA = cpu.r.af.a
	cpu.r.af.a += cpu.r.de.e
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH((cpu.r.af.a & 0xf) < (oldA & 0xf))
	cpu.r.setFlagC((cpu.r.af.a < oldA))
	return addAECycles
}

func addAH(cpu *cpu) cycleCount {
	// Add the value of register H into register A
	var oldA = cpu.r.af.a
	cpu.r.af.a += cpu.r.hl.h
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH((cpu.r.af.a & 0xf) < (oldA & 0xf))
	cpu.r.setFlagC((cpu.r.af.a < oldA))
	return addAHCycles
}

func addAL(cpu *cpu) cycleCount {
	// Add the value of register L into register A
	var oldA = cpu.r.af.a
	cpu.r.af.a += cpu.r.hl.l
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH((cpu.r.af.a & 0xf) < (oldA & 0xf))
	cpu.r.setFlagC((cpu.r.af.a < oldA))
	return addALCycles
}

func addAMemHl(cpu *cpu) cycleCount {
	// Add the value of the memory pointed by register HL into register A
	// TODO: To implement
	return addAMemHlCycles
}

func addANn(cpu *cpu) cycleCount {
	// Add the value of immediate value NN into register A
	// TODO: To implement
	return addANnCycles
}

// 3.3.3.3. SUB n
// Description:
// 		Subtract n from A.
// Use with:
// 		n = A,B,C,D,E,H,L,(HL),#
// Flags affected:
// 		Z - Set if result is zero.
// 		N - Set.
// 		H - Set if no borrow from bit 4.
// 		C - Set if no borrow.
// Opcodes:
// 		Instruction		Parameters		Opcode  	Cycles
// 		SUB				A, A			97			4
// 		SUB				A, B			90			4
// 		SUB				A, C			91			4
// 		SUB				A, D			92			4
// 		SUB				A, E			93			4
// 		SUB				A, H			94			4
// 		SUB				A, L			95			4
// 		SUB				A, (HL)			96			8
// 		SUB				A, #			D6			8

func subAA(cpu *cpu) cycleCount {
	// Subtract the value of register A to register A
	var oldA = cpu.r.af.a
	cpu.r.af.a -= cpu.r.af.a
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(true)
	cpu.r.setFlagH((oldA & 0xf) > (cpu.r.af.a & 0xf))
	cpu.r.setFlagC(oldA > cpu.r.af.a)
	return subAACycles
}

func subAB(cpu *cpu) cycleCount {
	// Subtract the value of register A to register B
	var oldA = cpu.r.af.a
	cpu.r.af.a -= cpu.r.bc.b
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(true)
	cpu.r.setFlagH((oldA & 0xf) > (cpu.r.af.a & 0xf))
	cpu.r.setFlagC(oldA > cpu.r.af.a)
	return subABCycles
}

func subAC(cpu *cpu) cycleCount {
	// Subtract the value of register C to register A
	var oldA = cpu.r.af.a
	cpu.r.af.a -= cpu.r.bc.c
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(true)
	cpu.r.setFlagH((oldA & 0xf) > (cpu.r.af.a & 0xf))
	cpu.r.setFlagC(oldA > cpu.r.af.a)
	return subACCycles
}

func subAD(cpu *cpu) cycleCount {
	// Subtract the value of register D to register A
	var oldA = cpu.r.af.a
	cpu.r.af.a -= cpu.r.de.d
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(true)
	cpu.r.setFlagH((oldA & 0xf) > (cpu.r.af.a & 0xf))
	cpu.r.setFlagC(oldA > cpu.r.af.a)
	return subADCycles
}

func subAE(cpu *cpu) cycleCount {
	// Subtract the value of register E to register A
	var oldA = cpu.r.af.a
	cpu.r.af.a -= cpu.r.de.e
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(true)
	cpu.r.setFlagH((oldA & 0xf) > (cpu.r.af.a & 0xf))
	cpu.r.setFlagC(oldA > cpu.r.af.a)
	return subAECycles
}

func subAH(cpu *cpu) cycleCount {
	// Subtract the value of register H to register A
	var oldA = cpu.r.af.a
	cpu.r.af.a -= cpu.r.hl.h
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(true)
	cpu.r.setFlagH((oldA & 0xf) > (cpu.r.af.a & 0xf))
	cpu.r.setFlagC(oldA > cpu.r.af.a)
	return subAHCycles
}

func subAL(cpu *cpu) cycleCount {
	// Subtract the value of register L to register A
	var oldA = cpu.r.af.a
	cpu.r.af.a -= cpu.r.hl.l
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(true)
	cpu.r.setFlagH((oldA & 0xf) > (cpu.r.af.a & 0xf))
	cpu.r.setFlagC(oldA > cpu.r.af.a)
	return subALCycles
}

func subAMemHl(cpu *cpu) cycleCount {
	// Subtract the value of the memory pointed by register HL to register A
	// TODO: to implement
	return subAMemHlCycles
}

func subANn(cpu *cpu) cycleCount {
	// Subtract the value of the immediate value NN to register A
	// TODO: to implement
	return subANnCycles
}

// 3.3.3.4. SBC A,n
// Description:
//		Subtract n + Carry flag from A.
// Use with:
//		n = A,B,C,D,E,H,L,(HL),#
// Flags affected:
//		Z - Set if result is zero.
//		N - Set.
//		H - Set if no borrow from bit 4.
//		C - Set if no borrow.
// Opcodes:
//		Instruction 	Parameters 		Opcode 	Cycles
//		SBC 			A,A 			9F 		4
//		SBC 			A,B 			98 		4
//		SBC 			A,C 			99 		4
//		SBC 			A,D 			9A 		4
//		SBC 			A,E 			9B 		4
//		SBC 			A,H 			9C 		4
//		SBC 			A,L 			9D 		4
//		SBC 			A,(HL) 			9E 		8
//		SBC 			A,# 			?? 		?

func sbcAA(cpu *cpu) cycleCount {
	// Subtract from register A the value of register A plus carry flag
	cpu.r.af.a -= byte(cpu.r.af.a)
	cpu.r.af.a -= cpu.r.flagAsInt(cpu.r.af.f.c)
	// TODO: set flags
	return sbcAACycles
}

func sbcAB(cpu *cpu) cycleCount {
	// Subtract from register A the value of register B plus carry flag
	cpu.r.af.a -= byte(cpu.r.bc.b)
	cpu.r.af.a -= cpu.r.flagAsInt(cpu.r.af.f.c)
	// TODO: set flags
	return sbcABCycles
}

func sbcAC(cpu *cpu) cycleCount {
	// Subtract from register A the value of register C plus carry flag
	cpu.r.af.a -= byte(cpu.r.bc.c)
	cpu.r.af.a -= cpu.r.flagAsInt(cpu.r.af.f.c)
	// TODO: set flags
	return sbcACCycles
}

func sbcAD(cpu *cpu) cycleCount {
	// Subtract from register A the value of register D plus carry flag
	cpu.r.af.a -= byte(cpu.r.de.d)
	cpu.r.af.a -= cpu.r.flagAsInt(cpu.r.af.f.c)
	// TODO: set flags
	return sbcADCycles
}

func sbcAE(cpu *cpu) cycleCount {
	// Subtract from register A the value of register E plus carry flag
	cpu.r.af.a -= byte(cpu.r.de.e)
	cpu.r.af.a -= cpu.r.flagAsInt(cpu.r.af.f.c)
	// TODO: set flags
	return sbcAECycles
}

func sbcAH(cpu *cpu) cycleCount {
	// Subtract from register A the value of register H plus carry flag
	cpu.r.af.a -= byte(cpu.r.hl.h)
	cpu.r.af.a -= cpu.r.flagAsInt(cpu.r.af.f.c)
	// TODO: set flags
	return sbcAHCycles
}

func sbcAL(cpu *cpu) cycleCount {
	// Subtract from register A the value of register L plus carry flag
	cpu.r.af.a -= byte(cpu.r.hl.l)
	cpu.r.af.a -= cpu.r.flagAsInt(cpu.r.af.f.c)
	// TODO: set flags
	return sbcALCycles
}

func sbcAMemHl(cpu *cpu) cycleCount {
	// Subtract from register A the value of memory pointed by HL plus carry flag
	//cpu.r.af.a -= ...
	//TODO: to implement
	return sbcAMemHlCycles
}

// There is missing documentation about this funcion
/*
func sbcAN(cpu *cpu) cycleCount {
	// Subtract from register A the value of register # plus carry flag
	// cpu.r.af.a -= ...
	//TODO: to implement
	return sbcANCycles
}
*/

// 3.3.3.5. AND n
// Description:
//		Logically AND n with A, result in A.
// Use with:
//		n = A,B,C,D,E,H,L,(HL),#
// Flags affected:
//		Z - Set if result is zero.
//		N - Reset.
//		H - Set.
//		C - Reset.
// Opcodes:
//		Instruction		Parameters		Opcode		Cycles
//		AND 			A 				A7 			4
//		AND 			B 				A0 			4
//		AND 			C 				A1 			4
//		AND 			D 				A2 			4
//		AND 			E 				A3 			4
//		AND 			H 				A4 			4
//		AND 			L 				A5 			4
//		AND 			(HL) 			A6 			8
//		AND 			# 				E6 			8

func andAA(cpu *cpu) cycleCount {
	// Stores into register A the result of (A & A) (bitwise AND)
	cpu.r.af.a = cpu.r.af.a & cpu.r.af.a
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	cpu.r.setFlagC(false)
	return andAACycles
}
func andAB(cpu *cpu) cycleCount {
	// Stores into register A the result of (A & B) (bitwise AND)
	cpu.r.af.a = cpu.r.af.a & cpu.r.bc.b
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	cpu.r.setFlagC(false)
	return andABCycles
}
func andAC(cpu *cpu) cycleCount {
	// Stores into register A the result of (A & C) (bitwise AND)
	cpu.r.af.a = cpu.r.af.a & cpu.r.bc.c
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	cpu.r.setFlagC(false)
	return andACCycles
}
func andAD(cpu *cpu) cycleCount {
	// Stores into register A the result of (A & D) (bitwise AND)
	cpu.r.af.a = cpu.r.af.a & cpu.r.de.d
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	cpu.r.setFlagC(false)
	return andADCycles
}
func andAE(cpu *cpu) cycleCount {
	// Stores into register A the result of (A & E) (bitwise AND)
	cpu.r.af.a = cpu.r.af.a & cpu.r.de.e
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	cpu.r.setFlagC(false)
	return andAECycles
}
func andAH(cpu *cpu) cycleCount {
	// Stores into register A the result of (A & H) (bitwise AND)
	cpu.r.af.a = cpu.r.af.a & cpu.r.hl.h
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	cpu.r.setFlagC(false)
	return andAHCycles
}
func andAL(cpu *cpu) cycleCount {
	// Stores into register A the result of (A & L) (bitwise AND)
	cpu.r.af.a = cpu.r.af.a & cpu.r.hl.l
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	cpu.r.setFlagC(false)
	return andALCycles
}
func andAMemHl(cpu *cpu) cycleCount {
	// Stores into register A the result of (A & the value of memory pointed by HL) (bitwise AND)
	// cpu.r.af.a = cpu.r.af.a & cpu.r.
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	cpu.r.setFlagC(false)
	// TODO: to implement
	return andAMemHlCycles
}
func andAN(cpu *cpu) cycleCount {
	// Stores into register A the result of (A & an immediate value) (bitwise AND)
	// cpu.r.af.a = cpu.r.af.a & cpu.r.
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	cpu.r.setFlagC(false)
	// TODO: to implement
	return andANCycles
}

// 3.3.3.6. OR n
// Description:
//	Logical OR n with register A, result in A.
// Use with:
// 	n = A,B,C,D,E,H,L,(HL),#
// Flags affected:
// 	Z - Set if result is zero.
// 	N - Reset.
// 	H - Reset.
// 	C - Reset.
// Opcodes:
//		Instruction 	Parameters 		Opcode 		Cycles
// 		OR 				A 				B7 			4
// 		OR 				B 				B0 			4
// 		OR 				C 				B1 			4
// 		OR 				D 				B2 			4
// 		OR 				E 				B3 			4
// 		OR 				H 				B4 			4
// 		OR 				L 				B5 			4
// 		OR 				(HL) 			B6 			8
// 		OR 				# 				F6 			8

func orAA(cpu *cpu) cycleCount {
	// Store into register A the result of (A | A) (bitwise OR)
	cpu.r.af.a = cpu.r.af.a | cpu.r.af.a
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(false)
	return orAACycles
}
func orAB(cpu *cpu) cycleCount {
	// Store into register A the result of (A | B) (bitwise OR)
	cpu.r.af.a = cpu.r.af.a | cpu.r.bc.b
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(false)
	return orABCycles
}
func orAC(cpu *cpu) cycleCount {
	// Store into register A the result of (A | C) (bitwise OR)
	cpu.r.af.a = cpu.r.af.a | cpu.r.bc.c
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(false)
	return orACCycles
}
func orAD(cpu *cpu) cycleCount {
	// Store into register A the result of (A | D) (bitwise OR)
	cpu.r.af.a = cpu.r.af.a | cpu.r.de.d
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(false)
	return orADCycles
}
func orAE(cpu *cpu) cycleCount {
	// Store into register A the result of (A | E) (bitwise OR)
	cpu.r.af.a = cpu.r.af.a | cpu.r.de.e
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(false)
	return orAECycles
}
func orAH(cpu *cpu) cycleCount {
	// Store into register A the result of (A | H) (bitwise OR)
	cpu.r.af.a = cpu.r.af.a | cpu.r.hl.h
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(false)
	return orAHCycles
}
func orAL(cpu *cpu) cycleCount {
	// Store into register A the result of (A | L) (bitwise OR)
	cpu.r.af.a = cpu.r.af.a | cpu.r.hl.l
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(false)
	return orALCycles
}
func orAMemHl(cpu *cpu) cycleCount {
	// Store into register A the result of (A | the memory position pointed by HL) (bitwise OR)
	//cpu.r.af.a = cpu.r.af.a | cpu.r.
	// TODO: to implement
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(false)
	return orAMemHlCycles
}
func orAN(cpu *cpu) cycleCount {
	// Store into register A the result of (A | an immediate value) (bitwise OR)
	//cpu.r.af.a = cpu.r.af.a | cpu.r.
	// TODO: to implement
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(false)
	return orANCycles
}

// 3.3.3.7. XOR n
// Description:
//	Logical exclusive OR n with register A, result in A.
// Use with:
//	n = A,B,C,D,E,H,L,(HL),#
// Flags affected:
//	Z - Set if result is zero.
//	N - Reset.
//	H - Reset.
//	C - Reset.
// Opcodes:
//		Instruction 	Parameters 		Opcode 		Cycles
//		XOR 			A 				AF 			4
//		XOR 			B 				A8 			4
//		XOR 			C 				A9 			4
//		XOR 			D 				AA 			4
//		XOR 			E 				AB 			4
//		XOR 			H 				AC 			4
//		XOR 			L 				AD 			4
//		XOR 			(HL) 			AE 			8
//		XOR 			* 				EE 			8

func xorAA(cpu *cpu) cycleCount {
	// Stores into register A the result of (A ^ A) (bitwise XOR)
	cpu.r.af.a ^= cpu.r.af.a
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(false)
	return xorAACycles
}
func xorAB(cpu *cpu) cycleCount {
	// Stores into register A the result of (A ^ B) (bitwise XOR)
	cpu.r.af.a ^= cpu.r.bc.b
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(false)
	return xorABCycles
}
func xorAC(cpu *cpu) cycleCount {
	// Stores into register A the result of (A ^ C) (bitwise XOR)
	cpu.r.af.a ^= cpu.r.bc.c
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(false)
	return xorACCycles
}
func xorAD(cpu *cpu) cycleCount {
	// Stores into register A the result of (A ^ D) (bitwise XOR)
	cpu.r.af.a ^= cpu.r.de.d
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(false)
	return xorADCycles
}
func xorAE(cpu *cpu) cycleCount {
	// Stores into register A the result of (A ^ E) (bitwise XOR)
	cpu.r.af.a ^= cpu.r.de.e
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(false)
	return xorAECycles
}
func xorAH(cpu *cpu) cycleCount {
	// Stores into register A the result of (A ^ H) (bitwise XOR)
	cpu.r.af.a ^= cpu.r.hl.h
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(false)
	return xorAHCycles
}
func xorAL(cpu *cpu) cycleCount {
	// Stores into register A the result of (A ^ L) (bitwise XOR)
	cpu.r.af.a ^= cpu.r.hl.l
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(false)
	return xorALCycles
}
func xorAMemHl(cpu *cpu) cycleCount {
	// Stores into register A the result of (A ^ the memory position pointed by HL) (bitwise XOR)
	//cpu.r.af.a ^= cpu.r.
	//TODO: to implement
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(false)
	return xorAMemHlCycles
}
func xorAN(cpu *cpu) cycleCount {
	// Stores into register A the result of (A ^ an immediate value) (bitwise XOR)
	// cpu.r.af.a ^= cpu.r.
	//TODO: to implement
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(false)
	return xorANCycles
}

// 3.3.3.8. CP n
// Description:
//	Compare A with n. This is basically an A - n
// 	subtraction instruction but the results are thrown
// 	away.
// Use with:
//	n = A,B,C,D,E,H,L,(HL),#
// Flags affected:
//	Z - Set if result is zero. (Set if A = n.)
//	N - Set.
//	H - Set if no borrow from bit 4.
//	C - Set for no borrow. (Set if A < n.)
// Opcodes:
//		Instruction 	Parameters 		Opcode 		Cycles
//		CP 				A 				BF 			4
//		CP 				B 				B8 			4
//		CP 				C 				B9 			4
//		CP 				D 				BA 			4
//		CP 				E 				BB 			4
//		CP 				H 				BC 			4
//		CP 				L 				BD 			4
//		CP 				(HL) 			BE 			8
//		CP 				# 				FE 			8

func cpAA(cpu *cpu) cycleCount {
	// Compares A to A. The result is not stored; this function only affects flags.
	cpu.r.setFlagZ(cpu.r.af.a == cpu.r.af.a)
	cpu.r.setFlagN(true)
	// TODO: cpu.r.setFlagH()
	cpu.r.setFlagC(cpu.r.af.a < cpu.r.af.a)
	return cpAACycles
}
func cpAB(cpu *cpu) cycleCount {
	// Compares A to B. The result is not stored; this function only affects flags.
	cpu.r.setFlagZ(cpu.r.af.a == cpu.r.af.a)
	cpu.r.setFlagN(true)
	// TODO: cpu.r.setFlagH()
	cpu.r.setFlagC(cpu.r.af.a < cpu.r.bc.b)
	return cpABCycles
}
func cpAC(cpu *cpu) cycleCount {
	// Compares A to C. The result is not stored; this function only affects flags.
	cpu.r.setFlagZ(cpu.r.af.a == cpu.r.bc.c)
	cpu.r.setFlagN(true)
	// TODO: cpu.r.setFlagH()
	cpu.r.setFlagC(cpu.r.af.a < cpu.r.af.a)
	return cpACCycles
}
func cpAD(cpu *cpu) cycleCount {
	// Compares A to D. The result is not stored; this function only affects flags.
	cpu.r.setFlagZ(cpu.r.af.a == cpu.r.af.a)
	cpu.r.setFlagN(true)
	// TODO: cpu.r.setFlagH()
	cpu.r.setFlagC(cpu.r.af.a < cpu.r.af.a)
	return cpADCycles
}
func cpAE(cpu *cpu) cycleCount {
	// Compares A to E. The result is not stored; this function only affects flags.
	cpu.r.setFlagZ(cpu.r.af.a == cpu.r.af.a)
	cpu.r.setFlagN(true)
	// TODO: cpu.r.setFlagH()
	cpu.r.setFlagC(cpu.r.af.a < cpu.r.af.a)
	return cpAECycles
}
func cpAH(cpu *cpu) cycleCount {
	// Compares A to H. The result is not stored; this function only affects flags.
	cpu.r.setFlagZ(cpu.r.af.a == cpu.r.af.a)
	cpu.r.setFlagN(true)
	// TODO: cpu.r.setFlagH()
	cpu.r.setFlagC(cpu.r.af.a < cpu.r.af.a)
	return cpAHCycles
}
func cpAL(cpu *cpu) cycleCount {
	// Compares A to L. The result is not stored; this function only affects flags.
	cpu.r.setFlagZ(cpu.r.af.a == cpu.r.af.a)
	cpu.r.setFlagN(true)
	// TODO: cpu.r.setFlagH()
	cpu.r.setFlagC(cpu.r.af.a < cpu.r.af.a)
	return cpALCycles
}
func cpAMemHl(cpu *cpu) cycleCount {
	// Compares A to (. The result is not stored; this function only affects flags.
	cpu.r.setFlagZ(cpu.r.af.a == cpu.r.af.a)
	cpu.r.setFlagN(true)
	// TODO: cpu.r.setFlagH()
	cpu.r.setFlagC(cpu.r.af.a < cpu.r.af.a)
	return cpAMemHlCycles
}
func cpAN(cpu *cpu) cycleCount {
	// Compares A to #. The result is not stored; this function only affects flags.
	cpu.r.setFlagZ(cpu.r.af.a == cpu.r.af.a)
	cpu.r.setFlagN(true)
	// TODO: cpu.r.setFlagH()
	cpu.r.setFlagC(cpu.r.af.a < cpu.r.af.a)
	return cpANCycles
}

// 3.3.3.9. INC n
// Description:
//	Increment register n.
// Use with:
//	n = A,B,C,D,E,H,L,(HL)
// Flags affected:
//	Z - Set if result is zero.
//	N - Reset.
//	H - Set if carry from bit 3.
//	C - Not affected.
// Opcodes:
//		Instruction 	Parameters 		Opcode 		Cycles
//		INC 			A 				3C 			4
//		INC 			B 				04 			4
//		INC 			C 				0C 			4
//		INC 			D 				14 			4
//		INC 			E 				1C 			4
//		INC 			H 				24 			4
//		INC 			L 				2C 			4
//		INC 			(HL) 			34 			12

func incA(cpu *cpu) cycleCount {
	// Increment register A 
	var oldA = cpu.r.af.a
	cpu.r.af.a++
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH((cpu.r.af.a & 0xf) < (oldA & 0xf))
	return incACycles
}
func incB(cpu *cpu) cycleCount {
	// Increment register B 
	var oldB = cpu.r.bc.b
	cpu.r.bc.b++
	cpu.r.setFlagZ(cpu.r.bc.b == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH((cpu.r.bc.b & 0xf) < (oldB & 0xf))
	return incBCycles
}
func incC(cpu *cpu) cycleCount {
	// Increment register C 
	var oldC = cpu.r.bc.c
	cpu.r.bc.c++
	cpu.r.setFlagZ(cpu.r.bc.c == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH((cpu.r.bc.c & 0xf) < (oldC & 0xf))
	return incCCycles
}
func incD(cpu *cpu) cycleCount {
	// Increment register D 
	var oldD = cpu.r.de.d
	cpu.r.de.d++
	cpu.r.setFlagZ(cpu.r.de.d == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH((cpu.r.de.d & 0xf) < (oldD & 0xf))
	return incDCycles
}
func incE(cpu *cpu) cycleCount {
	// Increment register E 
	var oldE = cpu.r.de.e
	cpu.r.de.e++
	cpu.r.setFlagZ(cpu.r.de.e == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH((cpu.r.de.e & 0xf) < (oldE & 0xf))
	return incECycles
}
func incH(cpu *cpu) cycleCount {
	// Increment register H 
	var oldH = cpu.r.hl.h
	cpu.r.hl.h++
	cpu.r.setFlagZ(cpu.r.hl.h == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH((cpu.r.hl.h & 0xf) < (oldH & 0xf))
	return incHCycles
}
func incL(cpu *cpu) cycleCount {
	// Increment register L 
	var oldL = cpu.r.hl.l
	cpu.r.hl.l++
	cpu.r.setFlagZ(cpu.r.hl.l == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH((cpu.r.hl.l & 0xf) < (oldL & 0xf))
	return incLCycles
}
func incMemHl(cpu *cpu) cycleCount {
	// Increment register (HL) 
	//cpu.r.
	//TODO: to implement
	return incMemHlCycles
}

// 3.3.3.10. DEC n
// Description:
//	Decrement register n.
// Use with:
//	n = A,B,C,D,E,H,L,(HL)
// Flags affected:
//	Z - Set if reselt is zero.
//	N - Set.
//	H - Set if no borrow from bit 4.
//	C - Not affected.
// Opcodes:
//		Instruction 	Parameters 		Opcode 		Cycles
//		DEC 			A 				3D 			4
//		DEC 			B 				05 			4
//		DEC 			C 				0D 			4
//		DEC 			D 				15 			4
//		DEC 			E 				1D 			4
//		DEC 			H 				25 			4
//		DEC 			L 				2D 			4
//		DEC 			(HL) 			35 			12

func decA(cpu *cpu) cycleCount {
	// Decrement register A
	var oldA = cpu.r.af.a
	cpu.r.af.a--
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(true)
	cpu.r.setFlagH((cpu.r.af.a & 0xf) > (oldA & 0xf))
	return decACycles
}
func decB(cpu *cpu) cycleCount {
	// Decrement register B
	var oldB = cpu.r.bc.b
	cpu.r.bc.b--
	cpu.r.setFlagZ(cpu.r.bc.b == 0)
	cpu.r.setFlagN(true)
	cpu.r.setFlagH((cpu.r.bc.b & 0xf) > (oldB & 0xf))
	return decBCycles
}
func decC(cpu *cpu) cycleCount {
	// Decrement register C
	var oldC = cpu.r.bc.c
	cpu.r.bc.c--
	cpu.r.setFlagZ(cpu.r.bc.c == 0)
	cpu.r.setFlagN(true)
	cpu.r.setFlagH((cpu.r.bc.c & 0xf) > (oldC & 0xf))
	return decCCycles
}
func decD(cpu *cpu) cycleCount {
	// Decrement register D
	var oldD = cpu.r.de.d
	cpu.r.de.d--
	cpu.r.setFlagZ(cpu.r.de.d == 0)
	cpu.r.setFlagN(true)
	cpu.r.setFlagH((cpu.r.de.d & 0xf) > (oldD & 0xf))
	return decDCycles
}
func decE(cpu *cpu) cycleCount {
	// Decrement register E
	var oldE = cpu.r.de.e
	cpu.r.de.e--
	cpu.r.setFlagZ(cpu.r.de.e == 0)
	cpu.r.setFlagN(true)
	cpu.r.setFlagH((cpu.r.de.e & 0xf) > (oldE & 0xf))
	return decECycles
}
func decH(cpu *cpu) cycleCount {
	// Decrement register H
	var oldH = cpu.r.hl.h
	cpu.r.hl.h--
	cpu.r.setFlagZ(cpu.r.hl.h == 0)
	cpu.r.setFlagN(true)
	cpu.r.setFlagH((cpu.r.hl.h & 0xf) > (oldH & 0xf))
	return decHCycles
}
func decL(cpu *cpu) cycleCount {
	// Decrement register L
	var oldL = cpu.r.hl.l
	cpu.r.hl.l--
	cpu.r.setFlagZ(cpu.r.hl.l == 0)
	cpu.r.setFlagN(true)
	cpu.r.setFlagH((cpu.r.hl.l & 0xf) > (oldL & 0xf))
	return decLCycles
}
func decMemHl(cpu *cpu) cycleCount {
	// Decrement memory position pointed by register HL
	//TODO: To implement
	return decMemHlCycles
}

// 3.3.4. 16-Bit Arithmetic

// 3.3.4.1. ADD HL,n
// Description:
//	Add n to HL.
// Use with:
//	n = BC,DE,HL,SP
// Flags affected:
//	Z - Not affected.
//	N - Reset.
//	H - Set if carry from bit 11.
//	C - Set if carry from bit 15.
// Opcodes:
//		Instruction 	Parameters 		Opcode 		Cycles
//		ADD 			HL,BC 			09 			8
//		ADD 			HL,DE 			19 			8
//		ADD 			HL,HL 			29 			8
//		ADD 			HL,SP 			39 			8

func addHlBc(cpu *cpu) cycleCount {
	// Add the value of register BC into register HL
	var carry = byte((uint16(cpu.r.hl.l) + uint16(cpu.r.bc.c)) >> 8)
	var oldH = cpu.r.hl.h
	cpu.r.hl.l += cpu.r.bc.c
	cpu.r.hl.h += carry + cpu.r.bc.b
	cpu.r.setFlagN(false)
	cpu.r.setFlagH((cpu.r.hl.h & 0xf) < (oldH & 0xf))
	cpu.r.setFlagC(cpu.r.hl.h < oldH)
	return addHlBcCycles
}

func addHlDe(cpu *cpu) cycleCount {
	// Add the value of register DE into register HL
	var carry = byte((uint16(cpu.r.hl.l) + uint16(cpu.r.de.e)) >> 8)
	var oldH = cpu.r.hl.h
	cpu.r.hl.l += cpu.r.de.e
	cpu.r.hl.h += carry + cpu.r.de.d
	cpu.r.setFlagN(false)
	cpu.r.setFlagH((cpu.r.hl.h & 0xf) < (oldH & 0xf))
	cpu.r.setFlagC(cpu.r.hl.h < oldH)
	return addHlDeCycles
}

func addHlHl(cpu *cpu) cycleCount {
	// Add the value of register HL into register HL
	var carry = byte((uint16(cpu.r.hl.l) + uint16(cpu.r.hl.l)) >> 8)
	var oldH = cpu.r.hl.h
	cpu.r.hl.l += cpu.r.hl.l
	cpu.r.hl.h += carry + cpu.r.hl.h
	cpu.r.setFlagN(false)
	cpu.r.setFlagH((cpu.r.hl.h & 0xf) < (oldH & 0xf))
	cpu.r.setFlagC(cpu.r.hl.h < oldH)
	return addHlHlCycles
}

func addHlSp(cpu *cpu) cycleCount {
	// Add the value of register SP into register HL
	var carry = byte((uint16(cpu.r.hl.l) + uint16(byte(cpu.r.spLow()))) >> 8)
	var oldH = cpu.r.hl.h
	cpu.r.hl.l += cpu.r.spLow()
	cpu.r.hl.h += carry + cpu.r.spHigh()
	cpu.r.setFlagN(false)
	cpu.r.setFlagH((cpu.r.hl.h & 0xf) < (oldH & 0xf))
	cpu.r.setFlagC(cpu.r.hl.h < oldH)
	return addHlSpCycles
}

// 3.3.4.2. ADD SP,n
// Description:
//	Add n to Stack Pointer (SP).
// Use with:
//	n = one byte signed immediate value (#).
// Flags affected:
//	Z - Reset.
//	N - Reset.
//	H - Set or reset according to operation.
//	C - Set or reset according to operation.
// Opcodes:
//		Instruction 	Parameters 		Opcode 		Cycles
//		ADD 			SP,# 			E8 			16

func addSpN(cpu *cpu) cycleCount {
	// Add the immediate value N to Stack Pointer (SP)
	var oldSpHigh = cpu.r.spHigh()
	//cpu.r.sp += N
	// TODO: To implement
	cpu.r.setFlagZ(false)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH((cpu.r.spHigh() & 0xf) < (oldSpHigh & 0xf))
	cpu.r.setFlagC(cpu.r.spHigh() < oldSpHigh)
	return addSpNCycles
}

// 3.3.4.3. INC nn
// Description:
//	Increment register nn.
// Use with:
//	nn = BC,DE,HL,SP
// Flags affected:
//	None.
// Opcodes:
//		Instruction 	Parameters 		Opcode 		Cycles
//		INC 			BC 				03 			8
//		INC 			DE 				13 			8
//		INC 			HL 				23 			8
//		INC 			SP 				33 			8

func incBc(cpu *cpu) cycleCount {
	// Increment register BC
	var bc = cpu.r.bcAsInt()
	bc++
	cpu.r.bc.c = byte(bc & 0xff)
	cpu.r.bc.b = byte(bc >> 8)
	return incBcCycles
}
func incDe(cpu *cpu) cycleCount {
	// Increment register DE
	var de = cpu.r.deAsInt()
	de++
	cpu.r.de.e = byte(de & 0xff)
	cpu.r.de.d = byte(de >> 8)
	return incDeCycles
}
func incHl(cpu *cpu) cycleCount {
	// Increment register HL
	var hl = cpu.r.hlAsInt()
	hl++
	cpu.r.hl.l = byte(hl & 0xff)
	cpu.r.hl.h = byte(hl >> 8)
	return incHlCycles
}
func incSp(cpu *cpu) cycleCount {
	// Increment register SP
	cpu.r.sp++
	return incSpCycles
}

// 3.3.4.4. DEC nn
// Description:
//	Decrement register nn.
// Use with:
//	nn = BC,DE,HL,SP
// Flags affected:
//	None.
// Opcodes:
//		Instruction 	Parameters 		Opcode 		Cycles
//		DEC 			BC 				0B 			8
//		DEC 			DE 				1B 			8
//		DEC 			HL 				2B 			8
//		DEC 			SP 				3B 			8

func decBc(cpu *cpu) cycleCount {
	// Decrement register BC
	var bc = cpu.r.bcAsInt()
	bc--
	cpu.r.bc.c = byte(bc & 0xff)
	cpu.r.bc.b = byte(bc >> 8)
	return decBcCycles
}
func decDe(cpu *cpu) cycleCount {
	// Decrement register DE
	var de = cpu.r.deAsInt()
	de--
	cpu.r.de.e = byte(de & 0xff)
	cpu.r.de.d = byte(de >> 8)
	return decDeCycles
}
func decHl(cpu *cpu) cycleCount {
	// Decrement register HL
	var hl = cpu.r.hlAsInt()
	hl--
	cpu.r.hl.l = byte(hl & 0xff)
	cpu.r.hl.h = byte(hl >> 8)
	return decHlCycles
}
func decSp(cpu *cpu) cycleCount {
	// Decrement register SP
	cpu.r.sp++
	return decSpCycles
}

// 3.3.5. Miscellaneous

// 3.3.5.1. SWAP n
// Description:
// 	Swap upper & lower nibles of n.
// Use with:
// 	n = A,B,C,D,E,H,L,(HL)
// Flags affected:
// 	Z - Set if result is zero.
// 	N - Reset.
// 	H - Reset.
// 	C - Reset.
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		SWAP 			A 				CB 37 			8
// 		SWAP 			B 				CB 30 			8
// 		SWAP 			C 				CB 31 			8
// 		SWAP 			D 				CB 32 			8
// 		SWAP 			E 				CB 33 			8
// 		SWAP 			H 				CB 34 			8
// 		SWAP 			L 				CB 35 			8
// 		SWAP 			(HL) 			CB 6 			16

var swapInstructions = map[byte]instructions{
	0x37: swapA,
	0x30: swapB,
	0x31: swapC,
	0x32: swapD,
	0x33: swapE,
	0x34: swapH,
	0x35: swapL,
	0x06: swapMemHl,
}

func swap(cpu *cpu) cycleCount {
	// Swap upper & lower nibles
	// This function gets the next opcode from memory,
	// and calls to the corresponding function
	var nextOpcode = byte(0)
	// TODO: Read the next opcode from memory
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(false)
	return swapInstructions[nextOpcode](cpu)
}

func swapA(cpu *cpu) cycleCount {
	// Swap upper & lower nibles of register A
	var oldA = cpu.r.af.a
	cpu.r.af.a = byte(oldA<<4) + byte(oldA<<4)
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	return 8
}
func swapB(cpu *cpu) cycleCount {
	// Swap upper & lower nibles of register B
	var oldB = cpu.r.bc.b
	cpu.r.bc.b = byte(oldB<<4) + byte(oldB<<4)
	cpu.r.setFlagZ(cpu.r.bc.b == 0)
	return 8
}
func swapC(cpu *cpu) cycleCount {
	// Swap upper & lower nibles of register C
	var oldC = cpu.r.bc.c
	cpu.r.bc.c = byte(oldC<<4) + byte(oldC<<4)
	cpu.r.setFlagZ(cpu.r.bc.c == 0)
	return 8
}
func swapD(cpu *cpu) cycleCount {
	// Swap upper & lower nibles of register D
	var oldD = cpu.r.de.d
	cpu.r.de.d = byte(oldD<<4) + byte(oldD<<4)
	cpu.r.setFlagZ(cpu.r.de.d == 0)
	return 8
}
func swapE(cpu *cpu) cycleCount {
	// Swap upper & lower nibles of register E
	var oldE = cpu.r.de.e
	cpu.r.de.e = byte(oldE<<4) + byte(oldE<<4)
	cpu.r.setFlagZ(cpu.r.de.e == 0)
	return 8
}
func swapH(cpu *cpu) cycleCount {
	// Swap upper & lower nibles of register H
	var oldH = cpu.r.hl.h
	cpu.r.hl.h = byte(oldH<<4) + byte(oldH<<4)
	cpu.r.setFlagZ(cpu.r.hl.h == 0)
	return 8
}
func swapL(cpu *cpu) cycleCount {
	// Swap upper & lower nibles of register L
	var oldL = cpu.r.hl.l
	cpu.r.hl.l = byte(oldL<<4) + byte(oldL<<4)
	cpu.r.setFlagZ(cpu.r.hl.l == 0)
	return 8
}
func swapMemHl(cpu *cpu) cycleCount {
	// Swap upper & lower nibles of the position of memory pointed by register HL
	// TODO: To implement
	return 16
}

// 3.3.5.2. DAA
// Description:
// 	Decimal adjust register A.
// 	This instruction adjusts register A so that the
// 	correct representation of Binary Coded Decimal (BCD)
// 	is obtained.
// Flags affected:
// 	Z - Set if register A is zero.
// 	N - Not affected.
// 	H - Reset.
// 	C - Set or reset according to operation.
//
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		DAA 			-/- 			27 			4

func daA(cpu *cpu) cycleCount {
	// TODO: To implement
	return daACycles
}

// 3.3.5.3. CPL
// Description:
// 	Complement A register. (Flip all bits.)
// Flags affected:
// 	Z - Not affected.
// 	N - Set.
// 	H - Set.
// 	C - Not affected.
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		CPL 			-/- 			2F 			4

func cplA(cpu *cpu) cycleCount {
	// Revert all bits of register A (ie. bitwise XOR with 0xFF)
	cpu.r.af.a ^= 0xFF
	return cplACycles
}

// 3.3.5.4. CCF
// Description:
// 	Complement carry flag.
// 	If C flag is set, then reset it.
// 	If C flag is reset, then set it.
// Flags affected:
// 	Z - Not affected.
// 	N - Reset.
// 	H - Reset.
// 	C - Complemented.
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		CCF 			-/- 			3F 			4
func ccf(cpu *cpu) cycleCount {
	// Toggle flag C
	cpu.r.setFlagC(!cpu.r.af.f.c)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	return ccfCycles
}

// 3.3.5.5. SCF
// Description:
// 	Set Carry flag.
// Flags affected:
// 	Z - Not affected.
// 	N - Reset.
// 	H - Reset.
// 	C - Set.
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		SCF 			-/- 			37 			4
func scf(cpu *cpu) cycleCount {
	// Sets the Carry Flag
	cpu.r.setFlagC(true)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	return scfCycles
}

// 3.3.5.6. NOP
// Description:
// 	No operation.
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		NOP 			-/- 			00 			4
func nop(cpu *cpu) cycleCount {
	// Does nothing
	return nopCycles
}

// 3.3.5.7. HALT
// Description:
// 	Power down CPU until an interrupt occurs. Use this
// 	when ever possible to reduce energy consumption.
// Opcodes:
// 	Instruction 	Parameters 		Opcode 		Cycles
// 	HALT 			-/- 			76 			4
func halt(cpu *cpu) cycleCount {
	// Power down CPU until an interrupt occurs
	// TODO: To implement
	return haltCycles
}

// 3.3.5.8. STOP
// Description:
// 	Halt CPU & LCD display until button pressed.
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		STOP 			-/- 			10 00 		4
func stop(cpu *cpu) cycleCount {
	// Stops the CPU until an interrupt occurs
	// TODO: To implement
	return stopCycles
}

// 3.3.5.9. DI
// Description:
// 	This instruction disables interrupts but not
// 	immediately. Interrupts are disabled after
// 	instruction after DI is executed.
// Flags affected:
// 	None.
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		DI 				-/- 			F3 			4
func di(cpu *cpu) cycleCount {
	// Disables Interruptions
	// TODO: To implement
	return diCycles
}

// 3.3.5.10. EI
// Description:
// 	Enable interrupts. This intruction enables interrupts
// 	but not immediately. Interrupts are enabled after
// 	instruction after EI is executed.
// Flags affected:
// 	None.
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		EI 				-/- 			FB 			4
func ei(cpu *cpu) cycleCount {
	// Enables Interruptions
	// TODO: To implement
	return eiCycles
}

// 3.3.6. Rotates & Shifts

// 3.3.6.1. RLCA
// Description:
// 	Rotate A left. Old bit 7 to Carry flag.
// Flags affected:
// 	Z - Set if result is zero.
// 	N - Reset.
// 	H - Reset.
// 	C - Contains old bit 7 data.
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		RLCA 			-/- 			07 			4
func rlcA(cpu *cpu) cycleCount {
	// Rotate A left 1 bit; A[0] = pre(A)[7]
	var bit7 = bool(cpu.r.af.a&0x80 == 0x80)

	cpu.r.af.a = cpu.r.af.a << 1
	if bit7 {
		cpu.r.af.a |= 0x01 // set bit 0 to 1
	}

	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit7)
	return rlcACycles
}

// 3.3.6.2. RLA
// Description:
// 	Rotate A left through Carry flag.
// Flags affected:
// 	Z - Set if result is zero.
// 	N - Reset.
// 	H - Reset.
// 	C - Contains old bit 7 data.
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		RLA 			-/- 			17 			4
func rlA(cpu *cpu) cycleCount {
	// Rotate A left 1 bit, but through Carry Flag
	// Note that Carry = A[7] and A[0] = Carry
	var bit7 = bool(cpu.r.af.a&0x80 == 0x80)

	cpu.r.af.a = cpu.r.af.a << 1
	if cpu.r.af.f.c {
		cpu.r.af.a |= 0x1 // set bit 0 to 1
	}

	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit7)
	return rlACycles
}

// 3.3.6.3. RRCA
// Description:
// 	Rotate A right. Old bit 0 to Carry flag.
// Flags affected:
// 	Z - Set if result is zero.
// 	N - Reset.
// 	H - Reset.
// 	C - Contains old bit 0 data.
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		RRCA 			-/- 			0F 			4
func rrcA(cpu *cpu) cycleCount {
	// Rotate A right 1 bit
	// Old bit 0 goes to Carry Flag
	var bit0 = bool(cpu.r.af.a&0x01 == 0x01)

	cpu.r.af.a = cpu.r.af.a >> 1
	if bit0 {
		cpu.r.af.a |= 0x80 // set bit 7 to 1
	}

	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit0)
	return rrcACycles
}

// 3.3.6.4. RRA
// Description:
// 	Rotate A right through Carry flag.
// Flags affected:
// 	Z - Set if result is zero.
// 	N - Reset.
// 	H - Reset.
// 	C - Contains old bit 0 data.
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		RRA 			-/- 			1F 			4
func rrA(cpu *cpu) cycleCount {
	// Rotate A right 1 bit, but through Carry Flag
	// Old bit 0 goes to Carry Flag, and old Carry Flag goes to bit 7
	var bit0 = bool(cpu.r.af.a&0x01 == 0x01)

	cpu.r.af.a = cpu.r.af.a >> 1
	if cpu.r.af.f.c {
		cpu.r.af.a |= 0x80 // set bit 7 to 1
	}

	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit0)
	return rrACycles
}

// 3.3.6.5. RLC n
// Description:
// 	Rotate n left. Old bit 7 to Carry flag.
// Use with:
// 	n = A,B,C,D,E,H,L,(HL)
// Flags affected:
// 	Z - Set if result is zero.
// 	N - Reset.
// 	H - Reset.
// 	C - Contains old bit 7 data.
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		RLC 			A 				CB 07 		8
// 		RLC 			B 				CB 00 		8
// 		RLC 			C 				CB 01 		8
// 		RLC 			D 				CB 02 		8
// 		RLC 			E 				CB 03 		8
// 		RLC 			H 				CB 04 		8
// 		RLC 			L 				CB 05 		8
// 		RLC 			(HL) 			CB 06 		16

var rxNInstructions = map[byte]instructions{
	// RLC
	0x07: rlcA,
	0x00: rlcB,
	0x01: rlcC,
	0x02: rlcD,
	0x03: rlcE,
	0x04: rlcH,
	0x05: rlcL,
	0x06: rlcMemHl,
	// RL
	0x17: rlA,
	0x10: rlB,
	0x11: rlC,
	0x12: rlD,
	0x13: rlE,
	0x14: rlH,
	0x15: rlL,
	0x16: rlMemHl,
	// RRC
	0x0F: rrcA,
	0x08: rrcB,
	0x09: rrcC,
	0x0A: rrcD,
	0x0B: rrcE,
	0x0C: rrcH,
	0x0D: rrcL,
	0x0E: rrcMemHl,
	// RR
	0x1F: rrA,
	0x18: rrB,
	0x19: rrC,
	0x1A: rrD,
	0x1B: rrE,
	0x1C: rrH,
	0x1D: rrL,
	0x1E: rrMemHl,
	// SLA
	0x27: slaA,
	0x20: slaB,
	0x21: slaC,
	0x22: slaD,
	0x23: slaE,
	0x24: slaH,
	0x25: slaL,
	0x26: slaMemHl,
	// SRA
	0x2F: sraA,
	0x28: sraB,
	0x29: sraC,
	0x2A: sraD,
	0x2B: sraE,
	0x2C: sraH,
	0x2D: sraL,
	0x2E: sraMemHl,
	// SRL
	0x3F: srlA,
	0x38: srlB,
	0x39: srlC,
	0x3A: srlD,
	0x3B: srlE,
	0x3C: srlH,
	0x3D: srlL,
	0x3E: srlMemHl,
	// BIT0
	0x47: bit0A,
	0x40: bit0B,
	0x41: bit0C,
	0x42: bit0D,
	0x43: bit0E,
	0x44: bit0H,
	0x45: bit0L,
	0x46: bit0MemHl,
	// BIT1
	0x4F: bit1A,
	0x48: bit1B,
	0x49: bit1C,
	0x4A: bit1D,
	0x4B: bit1E,
	0x4C: bit1H,
	0x4D: bit1L,
	0x4E: bit1MemHl,
	// BIT2
	0x57: bit2A,
	0x50: bit2B,
	0x51: bit2C,
	0x52: bit2D,
	0x53: bit2E,
	0x54: bit2H,
	0x55: bit2L,
	0x56: bit2MemHl,
	// BIT3
	0x5F: bit3A,
	0x58: bit3B,
	0x59: bit3C,
	0x5A: bit3D,
	0x5B: bit3E,
	0x5C: bit3H,
	0x5D: bit3L,
	0x5E: bit3MemHl,
	// BIT4
	0x67: bit4A,
	0x60: bit4B,
	0x61: bit4C,
	0x62: bit4D,
	0x63: bit4E,
	0x64: bit4H,
	0x65: bit4L,
	0x66: bit4MemHl,
	// BIT5
	0x6F: bit5A,
	0x68: bit5B,
	0x69: bit5C,
	0x6A: bit5D,
	0x6B: bit5E,
	0x6C: bit5H,
	0x6D: bit5L,
	0x6E: bit5MemHl,
	// BIT6
	0x77: bit6A,
	0x70: bit6B,
	0x71: bit6C,
	0x72: bit6D,
	0x73: bit6E,
	0x74: bit6H,
	0x75: bit6L,
	0x76: bit6MemHl,
	// BIT7
	0x7F: bit7A,
	0x78: bit7B,
	0x79: bit7C,
	0x7A: bit7D,
	0x7B: bit7E,
	0x7C: bit7H,
	0x7D: bit7L,
	0x7E: bit7MemHl,
	// RES0
	0x87: res0A,
	0x80: res0B,
	0x81: res0C,
	0x82: res0D,
	0x83: res0E,
	0x84: res0H,
	0x85: res0L,
	0x86: res0MemHl,
	// RES1
	0x8F: res1A,
	0x88: res1B,
	0x89: res1C,
	0x8A: res1D,
	0x8B: res1E,
	0x8C: res1H,
	0x8D: res1L,
	0x8E: res1MemHl,
	// RES2
	0x97: res2A,
	0x90: res2B,
	0x91: res2C,
	0x92: res2D,
	0x93: res2E,
	0x94: res2H,
	0x95: res2L,
	0x96: res2MemHl,
	// RES3
	0x9F: res3A,
	0x98: res3B,
	0x99: res3C,
	0x9A: res3D,
	0x9B: res3E,
	0x9C: res3H,
	0x9D: res3L,
	0x9E: res3MemHl,
	// RES4
	0xA7: res4A,
	0xA0: res4B,
	0xA1: res4C,
	0xA2: res4D,
	0xA3: res4E,
	0xA4: res4H,
	0xA5: res4L,
	0xA6: res4MemHl,
	// RES5
	0xAF: res5A,
	0xA8: res5B,
	0xA9: res5C,
	0xAA: res5D,
	0xAB: res5E,
	0xAC: res5H,
	0xAD: res5L,
	0xAE: res5MemHl,
	// RES6
	0xB7: res6A,
	0xB0: res6B,
	0xB1: res6C,
	0xB2: res6D,
	0xB3: res6E,
	0xB4: res6H,
	0xB5: res6L,
	0xB6: res6MemHl,
	// RES7
	0xBF: res7A,
	0xB8: res7B,
	0xB9: res7C,
	0xBA: res7D,
	0xBB: res7E,
	0xBC: res7H,
	0xBD: res7L,
	0xBE: res7MemHl,
	// SET0
	0xC7: set0A,
	0xC0: set0B,
	0xC1: set0C,
	0xC2: set0D,
	0xC3: set0E,
	0xC4: set0H,
	0xC5: set0L,
	0xC6: set0MemHl,
	// SET1
	0xCF: set1A,
	0xC8: set1B,
	0xC9: set1C,
	0xCA: set1D,
	0xCB: set1E,
	0xCC: set1H,
	0xCD: set1L,
	0xCE: set1MemHl,
	// SET2
	0xD7: set2A,
	0xD0: set2B,
	0xD1: set2C,
	0xD2: set2D,
	0xD3: set2E,
	0xD4: set2H,
	0xD5: set2L,
	0xD6: set2MemHl,
	// SET3
	0xDF: set3A,
	0xD8: set3B,
	0xD9: set3C,
	0xDA: set3D,
	0xDB: set3E,
	0xDC: set3H,
	0xDD: set3L,
	0xDE: set3MemHl,
	// SET4
	0xE7: set4A,
	0xE0: set4B,
	0xE1: set4C,
	0xE2: set4D,
	0xE3: set4E,
	0xE4: set4H,
	0xE5: set4L,
	0xE6: set4MemHl,
	// SET5
	0xEF: set5A,
	0xE8: set5B,
	0xE9: set5C,
	0xEA: set5D,
	0xEB: set5E,
	0xEC: set5H,
	0xED: set5L,
	0xEE: set5MemHl,
	// SET6
	0xF7: set6A,
	0xF0: set6B,
	0xF1: set6C,
	0xF2: set6D,
	0xF3: set6E,
	0xF4: set6H,
	0xF5: set6L,
	0xF6: set6MemHl,
	// SET7
	0xFF: set7A,
	0xF8: set7B,
	0xF9: set7C,
	0xFA: set7D,
	0xFB: set7E,
	0xFC: set7H,
	0xFD: set7L,
	0xFE: set7MemHl,
}

const (
	rlcBCycles     = 8
	rlcCCycles     = 8
	rlcDCycles     = 8
	rlcECycles     = 8
	rlcHCycles     = 8
	rlcLCycles     = 8
	rlcMemHlCycles = 16
)

func rxN(cpu *cpu) cycleCount {
	// Reads one opcode from memory,
	// and decides which RL/RLC/RR/RRC function to call
	var nextOpcode = byte(0)
	// TODO: Read the next opcode from memory
	var cycles = rxNInstructions[nextOpcode](cpu)
	if cycles < rxNCycles {
		return rxNCycles
	}
	return cycles
}

// func rlcA(cpu *cpu) cycleCount
// already implemented

func rlcB(cpu *cpu) cycleCount {
	// Rotate B left 1 bit; B[0] = pre(B)[7]
	var bit7 = bool(cpu.r.bc.b&0x80 == 0x80)

	cpu.r.bc.b = cpu.r.bc.b << 1
	if bit7 {
		cpu.r.bc.b |= 0x01 // set bit 0 to 1
	}

	cpu.r.setFlagZ(cpu.r.bc.b == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit7)
	return rlcBCycles
}

func rlcC(cpu *cpu) cycleCount {
	// Rotate C left 1 bit; C[0] = pre(C)[7]
	var bit7 = bool(cpu.r.bc.c&0x80 == 0x80)

	cpu.r.bc.c = cpu.r.bc.c << 1
	if bit7 {
		cpu.r.bc.c |= 0x01 // set bit 0 to 1
	}

	cpu.r.setFlagZ(cpu.r.bc.c == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit7)
	return rlcCCycles
}

func rlcD(cpu *cpu) cycleCount {
	// Rotate D left 1 bit; D[0] = pre(D)[7]
	var bit7 = bool(cpu.r.de.d&0x80 == 0x80)

	cpu.r.de.d = cpu.r.de.d << 1
	if bit7 {
		cpu.r.de.d |= 0x01 // set bit 0 to 1
	}

	cpu.r.setFlagZ(cpu.r.de.d == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit7)
	return rlcDCycles
}

func rlcE(cpu *cpu) cycleCount {
	// Rotate E left 1 bit; E[0] = pre(E)[7]
	var bit7 = bool(cpu.r.de.e&0x80 == 0x80)

	cpu.r.de.e = cpu.r.de.e << 1
	if bit7 {
		cpu.r.de.e |= 0x01 // set bit 0 to 1
	}

	cpu.r.setFlagZ(cpu.r.de.e == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit7)
	return rlcECycles
}

func rlcH(cpu *cpu) cycleCount {
	// Rotate H left 1 bit; H[0] = pre(H)[7]
	var bit7 = bool(cpu.r.hl.h&0x80 == 0x80)

	cpu.r.hl.h = cpu.r.hl.h << 1
	if bit7 {
		cpu.r.hl.h |= 0x01 // set bit 0 to 1
	}

	cpu.r.setFlagZ(cpu.r.hl.h == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit7)
	return rlcHCycles
}

func rlcL(cpu *cpu) cycleCount {
	// Rotate L left 1 bit; L[0] = pre(L)[7]
	var bit7 = bool(cpu.r.hl.l&0x80 == 0x80)

	cpu.r.hl.l = cpu.r.hl.l << 1
	if bit7 {
		cpu.r.hl.l |= 0x01 // set bit 0 to 1
	}

	cpu.r.setFlagZ(cpu.r.hl.l == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit7)
	return rlcLCycles
}

func rlcMemHl(cpu *cpu) cycleCount {
	// TODO: To implement
	return rlcMemHlCycles
}

// 3.3.6.6. RL n
// Description:
// 	Rotate n left through Carry flag.
// Use with:
// 	n = A,B,C,D,E,H,L,(HL)
// Flags affected:
// 	Z - Set if result is zero.
// 	N - Reset.
// 	H - Reset.
// 	C - Contains old bit 7 data.
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		RL 				A 				CB 17 		8
// 		RL 				B 				CB 10 		8
// 		RL 				C 				CB 11 		8
// 		RL 				D 				CB 12 		8
// 		RL 				E 				CB 13 		8
// 		RL 				H 				CB 14 		8
// 		RL 				L 				CB 15 		8
// 		RL 				(HL) 			CB 16 		16

const (
	rlBCycles     = 8
	rlCCycles     = 8
	rlDCycles     = 8
	rlECycles     = 8
	rlHCycles     = 8
	rlLCycles     = 8
	rlMemHlCycles = 16
)

// func rlA(cpu *cpu) cycleCount
// already implemented

func rlB(cpu *cpu) cycleCount {
	// Rotate B left 1 bit, but through Carry Flag
	// Note that Carry = B[7] and B[0] = Carry
	var bit7 = bool(cpu.r.bc.b&0x80 == 0x80)

	cpu.r.bc.b = cpu.r.bc.b << 1
	if cpu.r.af.f.c {
		cpu.r.bc.b |= 0x1 // set bit 0 to 1
	}

	cpu.r.setFlagZ(cpu.r.bc.b == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit7)
	return rlBCycles
}

func rlC(cpu *cpu) cycleCount {
	// Rotate C left 1 bit, but through Carry Flag
	// Note that Carry = C[7] and C[0] = Carry
	var bit7 = bool(cpu.r.bc.c&0x80 == 0x80)

	cpu.r.bc.c = cpu.r.bc.c << 1
	if cpu.r.af.f.c {
		cpu.r.bc.c |= 0x1 // set bit 0 to 1
	}

	cpu.r.setFlagZ(cpu.r.bc.c == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit7)
	return rlCCycles
}

func rlD(cpu *cpu) cycleCount {
	// Rotate D left 1 bit, but through Carry Flag
	// Note that Carry = D[7] and D[0] = Carry
	var bit7 = bool(cpu.r.de.d&0x80 == 0x80)

	cpu.r.de.d = cpu.r.de.d << 1
	if cpu.r.af.f.c {
		cpu.r.de.d |= 0x1 // set bit 0 to 1
	}

	cpu.r.setFlagZ(cpu.r.de.d == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit7)
	return rlDCycles
}

func rlE(cpu *cpu) cycleCount {
	// Rotate E left 1 bit, but through Carry Flag
	// Note that Carry = E[7] and E[0] = Carry
	var bit7 = bool(cpu.r.de.e&0x80 == 0x80)

	cpu.r.de.e = cpu.r.de.e << 1
	if cpu.r.af.f.c {
		cpu.r.de.e |= 0x1 // set bit 0 to 1
	}

	cpu.r.setFlagZ(cpu.r.de.e == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit7)
	return rlECycles
}

func rlH(cpu *cpu) cycleCount {
	// Rotate H left 1 bit, but through Carry Flag
	// Note that Carry = H[7] and H[0] = Carry
	var bit7 = bool(cpu.r.hl.h&0x80 == 0x80)

	cpu.r.hl.h = cpu.r.hl.h << 1
	if cpu.r.af.f.c {
		cpu.r.hl.h |= 0x1 // set bit 0 to 1
	}

	cpu.r.setFlagZ(cpu.r.hl.h == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit7)
	return rlHCycles
}

func rlL(cpu *cpu) cycleCount {
	// Rotate L left 1 bit, but through Carry Flag
	// Note that Carry = L[7] and L[0] = Carry
	var bit7 = bool(cpu.r.hl.l&0x80 == 0x80)

	cpu.r.hl.l = cpu.r.hl.l << 1
	if cpu.r.af.f.c {
		cpu.r.hl.l |= 0x1 // set bit 0 to 1
	}

	cpu.r.setFlagZ(cpu.r.hl.l == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit7)
	return rlLCycles
}

func rlMemHl(cpu *cpu) cycleCount {
	// TODO: To implement
	return rlMemHlCycles
}

// 3.3.6.7. RRC n
// Description:
// 	Rotate n right. Old bit 0 to Carry flag.
// Use with:
// 	n = A,B,C,D,E,H,L,(HL)
// Flags affected:
// 	Z - Set if result is zero.
// 	N - Reset.
// 	H - Reset.
// 	C - Contains old bit 0 data.
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		RRC 			A 				CB 0F 		8
// 		RRC 			B 				CB 08 		8
// 		RRC 			C 				CB 09 		8
// 		RRC 			D 				CB 0A 		8
// 		RRC 			E 				CB 0B 		8
// 		RRC 			H 				CB 0C 		8
// 		RRC 			L 				CB 0D 		8
// 		RRC 			(HL) 			CB 0E 		16

// func rrcA(cpu *cpu) cycleCount
// already implemented

const (
	rrcBCycles     = 8
	rrcCCycles     = 8
	rrcDCycles     = 8
	rrcECycles     = 8
	rrcHCycles     = 8
	rrcLCycles     = 8
	rrcMemHlCycles = 16
)

func rrcB(cpu *cpu) cycleCount {
	// Rotate B right 1 bit
	// Old bit 0 goes to Carry Flag
	var bit0 = bool(cpu.r.bc.b&0x01 == 0x01)

	cpu.r.bc.b = cpu.r.bc.b >> 1
	if bit0 {
		cpu.r.bc.b |= 0x80 // set bit 7 to 1
	}

	cpu.r.setFlagZ(cpu.r.bc.b == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit0)
	return rrcBCycles
}

func rrcC(cpu *cpu) cycleCount {
	// Rotate C right 1 bit
	// Old bit 0 goes to Carry Flag
	var bit0 = bool(cpu.r.bc.c&0x01 == 0x01)

	cpu.r.bc.c = cpu.r.bc.c >> 1
	if bit0 {
		cpu.r.bc.c |= 0x80 // set bit 7 to 1
	}

	cpu.r.setFlagZ(cpu.r.bc.c == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit0)
	return rrcCCycles
}

func rrcD(cpu *cpu) cycleCount {
	// Rotate D right 1 bit
	// Old bit 0 goes to Carry Flag
	var bit0 = bool(cpu.r.de.d&0x01 == 0x01)

	cpu.r.de.d = cpu.r.de.d >> 1
	if bit0 {
		cpu.r.de.d |= 0x80 // set bit 7 to 1
	}

	cpu.r.setFlagZ(cpu.r.de.d == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit0)
	return rrcDCycles
}

func rrcE(cpu *cpu) cycleCount {
	// Rotate E right 1 bit
	// Old bit 0 goes to Carry Flag
	var bit0 = bool(cpu.r.de.e&0x01 == 0x01)

	cpu.r.de.e = cpu.r.de.e >> 1
	if bit0 {
		cpu.r.de.e |= 0x80 // set bit 7 to 1
	}

	cpu.r.setFlagZ(cpu.r.de.e == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit0)
	return rrcECycles
}

func rrcH(cpu *cpu) cycleCount {
	// Rotate H right 1 bit
	// Old bit 0 goes to Carry Flag
	var bit0 = bool(cpu.r.hl.h&0x01 == 0x01)

	cpu.r.hl.h = cpu.r.hl.h >> 1
	if bit0 {
		cpu.r.hl.h |= 0x80 // set bit 7 to 1
	}

	cpu.r.setFlagZ(cpu.r.hl.h == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit0)
	return rrcHCycles
}

func rrcL(cpu *cpu) cycleCount {
	// Rotate L right 1 bit
	// Old bit 0 goes to Carry Flag
	var bit0 = bool(cpu.r.hl.l&0x01 == 0x01)

	cpu.r.hl.l = cpu.r.hl.l >> 1
	if bit0 {
		cpu.r.hl.l |= 0x80 // set bit 7 to 1
	}

	cpu.r.setFlagZ(cpu.r.hl.l == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit0)
	return rrcLCycles
}

func rrcMemHl(cpu *cpu) cycleCount {
	// TODO: To implement
	return rrcMemHlCycles
}

// 3.3.6.8. RR n
// Description:
// 	Rotate n right through Carry flag.
// Use with:
// 	n = A,B,C,D,E,H,L,(HL)
// Flags affected:
// 	Z - Set if result is zero.
// 	N - Reset.
// 	H - Reset.
// 	C - Contains old bit 0 data.
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		RR 				A 				CB 1F 		8
// 		RR 				B 				CB 18 		8
// 		RR 				C 				CB 19 		8
// 		RR 				D 				CB 1A 		8
// 		RR 				E 				CB 1B 		8
// 		RR 				H 				CB 1C 		8
// 		RR 				L 				CB 1D 		8
// 		RR 				(HL) 			CB 1E 		16

// func rrA(cpu *cpu) cycleCount
// already implemented

const (
	rrBCycles     = 8
	rrCCycles     = 8
	rrDCycles     = 8
	rrECycles     = 8
	rrHCycles     = 8
	rrLCycles     = 8
	rrMemHlCycles = 16
)

func rrB(cpu *cpu) cycleCount {
	// Rotate B right 1 bit, but through Carry Flag
	// Old bit 0 goes to Carry Flag, and old Carry Flag goes to bit 7
	var bit0 = bool(cpu.r.bc.b&0x01 == 0x01)

	cpu.r.bc.b = cpu.r.bc.b >> 1
	if cpu.r.af.f.c {
		cpu.r.bc.b |= 0x80 // set bit 7 to 1
	}

	cpu.r.setFlagZ(cpu.r.bc.b == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit0)
	return rrBCycles
}

func rrC(cpu *cpu) cycleCount {
	// Rotate C right 1 bit, but through Carry Flag
	// Old bit 0 goes to Carry Flag, and old Carry Flag goes to bit 7
	var bit0 = bool(cpu.r.bc.c&0x01 == 0x01)

	cpu.r.bc.c = cpu.r.bc.c >> 1
	if cpu.r.af.f.c {
		cpu.r.bc.c |= 0x80 // set bit 7 to 1
	}

	cpu.r.setFlagZ(cpu.r.bc.c == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit0)
	return rrCCycles
}

func rrD(cpu *cpu) cycleCount {
	// Rotate D right 1 bit, but through Carry Flag
	// Old bit 0 goes to Carry Flag, and old Carry Flag goes to bit 7
	var bit0 = bool(cpu.r.de.d&0x01 == 0x01)

	cpu.r.de.d = cpu.r.de.d >> 1
	if cpu.r.af.f.c {
		cpu.r.de.d |= 0x80 // set bit 7 to 1
	}

	cpu.r.setFlagZ(cpu.r.de.d == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit0)
	return rrDCycles
}

func rrE(cpu *cpu) cycleCount {
	// Rotate E right 1 bit, but through Carry Flag
	// Old bit 0 goes to Carry Flag, and old Carry Flag goes to bit 7
	var bit0 = bool(cpu.r.de.e&0x01 == 0x01)

	cpu.r.de.e = cpu.r.de.e >> 1
	if cpu.r.af.f.c {
		cpu.r.de.e |= 0x80 // set bit 7 to 1
	}

	cpu.r.setFlagZ(cpu.r.de.e == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit0)
	return rrECycles
}

func rrH(cpu *cpu) cycleCount {
	// Rotate H right 1 bit, but through Carry Flag
	// Old bit 0 goes to Carry Flag, and old Carry Flag goes to bit 7
	var bit0 = bool(cpu.r.hl.h&0x01 == 0x01)

	cpu.r.hl.h = cpu.r.hl.h >> 1
	if cpu.r.af.f.c {
		cpu.r.hl.h |= 0x80 // set bit 7 to 1
	}

	cpu.r.setFlagZ(cpu.r.hl.h == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit0)
	return rrHCycles
}

func rrL(cpu *cpu) cycleCount {
	// Rotate L right 1 bit, but through Carry Flag
	// Old bit 0 goes to Carry Flag, and old Carry Flag goes to bit 7
	var bit0 = bool(cpu.r.hl.l&0x01 == 0x01)

	cpu.r.hl.l = cpu.r.hl.l >> 1
	if cpu.r.af.f.c {
		cpu.r.hl.l |= 0x80 // set bit 7 to 1
	}

	cpu.r.setFlagZ(cpu.r.hl.l == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	cpu.r.setFlagC(bit0)
	return rrLCycles
}

func rrMemHl(cpu *cpu) cycleCount {
	return rrMemHlCycles
}

// 3.3.6.9. SLA n
// Description:
// 	Shift n left into Carry. LSB of n set to 0.
// Use with:
// 	n = A,B,C,D,E,H,L,(HL)
// Flags affected:
// 	Z - Set if result is zero.
// 	N - Reset.
// 	H - Reset.
// 	C - Contains old bit 7 data.
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		SLA 			A 				CB 27 		8
// 		SLA 			B 				CB 20 		8
// 		SLA 			C 				CB 21 		8
// 		SLA 			D 				CB 22 		8
// 		SLA 			E 				CB 23 		8
// 		SLA 			H 				CB 24 		8
// 		SLA 			L 				CB 25 		8
// 		SLA 			(HL) 			CB 26 		16

const (
	slaACycles     = 8
	slaBCycles     = 8
	slaCCycles     = 8
	slaDCycles     = 8
	slaECycles     = 8
	slaHCycles     = 8
	slaLCycles     = 8
	slaMemHlCycles = 16
)

func slaA(cpu *cpu) cycleCount {
	// Shift A left into Carry
	// A[0] set to 0
	cpu.r.setFlagC((cpu.r.af.a & 0x80) == 0x80)
	cpu.r.af.a = cpu.r.af.a << 1
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	return slaACycles
}

func slaB(cpu *cpu) cycleCount {
	// Shift B left into Carry
	// B[0] set to 0
	cpu.r.setFlagC((cpu.r.bc.b & 0x80) == 0x80)
	cpu.r.bc.b = cpu.r.bc.b << 1
	cpu.r.setFlagZ(cpu.r.bc.b == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	return slaBCycles
}

func slaC(cpu *cpu) cycleCount {
	// Shift C left into Carry
	// C[0] set to 0
	cpu.r.setFlagC((cpu.r.bc.c & 0x80) == 0x80)
	cpu.r.bc.c = cpu.r.bc.c << 1
	cpu.r.setFlagZ(cpu.r.bc.c == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	return slaCCycles
}

func slaD(cpu *cpu) cycleCount {
	// Shift D left into Carry
	// D[0] set to 0
	cpu.r.setFlagC((cpu.r.de.d & 0x80) == 0x80)
	cpu.r.de.d = cpu.r.de.d << 1
	cpu.r.setFlagZ(cpu.r.de.d == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	return slaDCycles
}

func slaE(cpu *cpu) cycleCount {
	// Shift E left into Carry
	// E[0] set to 0
	cpu.r.setFlagC((cpu.r.de.e & 0x80) == 0x80)
	cpu.r.de.e = cpu.r.de.e << 1
	cpu.r.setFlagZ(cpu.r.de.e == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	return slaECycles
}

func slaH(cpu *cpu) cycleCount {
	// Shift H left into Carry
	// H[0] set to 0
	cpu.r.setFlagC((cpu.r.hl.h & 0x80) == 0x80)
	cpu.r.hl.h = cpu.r.hl.h << 1
	cpu.r.setFlagZ(cpu.r.hl.h == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	return slaHCycles
}

func slaL(cpu *cpu) cycleCount {
	// Shift L left into Carry
	// L[0] set to 0
	cpu.r.setFlagC((cpu.r.hl.l & 0x80) == 0x80)
	cpu.r.hl.l = cpu.r.hl.l << 1
	cpu.r.setFlagZ(cpu.r.hl.l == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	return slaLCycles
}

func slaMemHl(cpu *cpu) cycleCount {
	// Shift MemHl left into Carry
	// MemHl[0] set to 0
	// TODO: To implement
	return slaMemHlCycles
}

// 3.3.6.10. SRA n
// Description:
// 	Shift n right into Carry. MSB doesn't change.
// Use with:
// 	n = A,B,C,D,E,H,L,(HL)
// Flags affected:
// 	Z - Set if result is zero.
// 	N - Reset.
// 	H - Reset.
// 	C - Contains old bit 0 data.
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		SRA 			A 				CB 2F 		8
// 		SRA 			B 				CB 28 		8
// 		SRA 			C 				CB 29 		8
// 		SRA 			D 				CB 2A 		8
// 		SRA 			E 				CB 2B 		8
// 		SRA 			H 				CB 2C 		8
// 		SRA 			L 				CB 2D 		8
// 		SRA 			(HL) 			CB 2E		16

const (
	sraACycles     = 8
	sraBCycles     = 8
	sraCCycles     = 8
	sraDCycles     = 8
	sraECycles     = 8
	sraHCycles     = 8
	sraLCycles     = 8
	sraMemHlCycles = 16
)

func sraA(cpu *cpu) cycleCount {
	// Shift A right into carry
	// A[7] doesn't change
	var msb = cpu.r.af.a & 0x80
	cpu.r.setFlagC((cpu.r.af.a & 0x01) == 0x01)
	cpu.r.af.a = (cpu.r.af.a >> 1) | msb // restores MSB
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	return sraACycles
}

func sraB(cpu *cpu) cycleCount {
	// Shift B right into carry
	// B[7] doesn't change
	var msb = cpu.r.bc.b & 0x80
	cpu.r.setFlagC((cpu.r.bc.b & 0x01) == 0x01)
	cpu.r.bc.b = (cpu.r.bc.b >> 1) | msb // restores MSB
	cpu.r.setFlagZ(cpu.r.bc.b == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	return sraBCycles
}

func sraC(cpu *cpu) cycleCount {
	// Shift C right into carry
	// C[7] doesn't change
	var msb = cpu.r.bc.c & 0x80
	cpu.r.setFlagC((cpu.r.bc.c & 0x01) == 0x01)
	cpu.r.bc.c = (cpu.r.bc.c >> 1) | msb // restores MSB
	cpu.r.setFlagZ(cpu.r.bc.c == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	return sraCCycles
}

func sraD(cpu *cpu) cycleCount {
	// Shift D right into carry
	// D[7] doesn't change
	var msb = cpu.r.de.d & 0x80
	cpu.r.setFlagC((cpu.r.de.d & 0x01) == 0x01)
	cpu.r.de.d = (cpu.r.de.d >> 1) | msb // restores MSB
	cpu.r.setFlagZ(cpu.r.de.d == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	return sraDCycles
}

func sraE(cpu *cpu) cycleCount {
	// Shift E right into carry
	// E[7] doesn't change
	var msb = cpu.r.de.e & 0x80
	cpu.r.setFlagC((cpu.r.de.e & 0x01) == 0x01)
	cpu.r.de.e = (cpu.r.de.e >> 1) | msb // restores MSB
	cpu.r.setFlagZ(cpu.r.de.e == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	return sraECycles
}

func sraH(cpu *cpu) cycleCount {
	// Shift H right into carry
	// H[7] doesn't change
	var msb = cpu.r.hl.h & 0x80
	cpu.r.setFlagC((cpu.r.hl.h & 0x01) == 0x01)
	cpu.r.hl.h = (cpu.r.hl.h >> 1) | msb // restores MSB
	cpu.r.setFlagZ(cpu.r.hl.h == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	return sraHCycles
}

func sraL(cpu *cpu) cycleCount {
	// Shift L right into carry
	// L[7] doesn't change
	var msb = cpu.r.hl.l & 0x80
	cpu.r.setFlagC((cpu.r.hl.l & 0x01) == 0x01)
	cpu.r.hl.l = (cpu.r.hl.l >> 1) | msb // restores MSB
	cpu.r.setFlagZ(cpu.r.hl.l == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	return sraLCycles
}

func sraMemHl(cpu *cpu) cycleCount {
	// Shift MemHl right into carry
	// MemHl[7] doesn't change
	return sraMemHlCycles
}

// 3.3.6.11. SRL n
// Description:
// 	Shift n right into Carry. MSB set to 0.
// Use with:
// 	n = A,B,C,D,E,H,L,(HL)
// Flags affected:
// 	Z - Set if result is zero.
// 	N - Reset.
// 	H - Reset.
// 	C - Contains old bit 0 data.
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		SRL 			A 				CB 3F 		8
// 		SRL 			B 				CB 38 		8
// 		SRL 			C 				CB 39 		8
// 		SRL 			D 				CB 3A 		8
// 		SRL 			E 				CB 3B 		8
// 		SRL 			H 				CB 3C 		8
// 		SRL 			L 				CB 3D 		8
// 		SRL 			(HL) 			CB 3E 		16

const (
	srlACycles     = 8
	srlBCycles     = 8
	srlCCycles     = 8
	srlDCycles     = 8
	srlECycles     = 8
	srlHCycles     = 8
	srlLCycles     = 8
	srlMemHlCycles = 16
)

func srlA(cpu *cpu) cycleCount {
	// Shift A right into Carry.
	// A[7] = 0
	cpu.r.setFlagC((cpu.r.af.a & 0x01) == 0x01)
	cpu.r.af.a = cpu.r.af.a >> 1
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	return srlACycles
}

func srlB(cpu *cpu) cycleCount {
	// Shift B right into Carry.
	// B[7] = 0
	cpu.r.setFlagC((cpu.r.bc.b & 0x01) == 0x01)
	cpu.r.bc.b = cpu.r.bc.b >> 1
	cpu.r.setFlagZ(cpu.r.bc.b == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	return srlBCycles
}

func srlC(cpu *cpu) cycleCount {
	// Shift C right into Carry.
	// C[7] = 0
	cpu.r.setFlagC((cpu.r.bc.c & 0x01) == 0x01)
	cpu.r.bc.c = cpu.r.bc.c >> 1
	cpu.r.setFlagZ(cpu.r.bc.c == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	return srlCCycles
}

func srlD(cpu *cpu) cycleCount {
	// Shift D right into Carry.
	// D[7] = 0
	cpu.r.setFlagC((cpu.r.de.d & 0x01) == 0x01)
	cpu.r.de.d = cpu.r.de.d >> 1
	cpu.r.setFlagZ(cpu.r.de.d == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	return srlDCycles
}

func srlE(cpu *cpu) cycleCount {
	// Shift E right into Carry.
	// E[7] = 0
	cpu.r.setFlagC((cpu.r.de.e & 0x01) == 0x01)
	cpu.r.de.e = cpu.r.de.e >> 1
	cpu.r.setFlagZ(cpu.r.de.e == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	return srlECycles
}

func srlH(cpu *cpu) cycleCount {
	// Shift H right into Carry.
	// H[7] = 0
	cpu.r.setFlagC((cpu.r.hl.h & 0x01) == 0x01)
	cpu.r.hl.h = cpu.r.hl.h >> 1
	cpu.r.setFlagZ(cpu.r.hl.h == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	return srlHCycles
}

func srlL(cpu *cpu) cycleCount {
	// Shift L right into Carry.
	// L[7] = 0
	cpu.r.setFlagC((cpu.r.hl.l & 0x01) == 0x01)
	cpu.r.hl.l = cpu.r.hl.l >> 1
	cpu.r.setFlagZ(cpu.r.hl.l == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(false)
	return srlLCycles
}

func srlMemHl(cpu *cpu) cycleCount {
	// Shift MemHl right into Carry.
	// MemHl[7] = 0
	return srlMemHlCycles
}

// 3.3.7. Bit Opcodes

// 3.3.7.1. BIT b,r
// Description:
// 	Test bit b in register r.
// Use with:
// 	b = 0 - 7, r = A,B,C,D,E,H,L,(HL)
// Flags affected:
// 	Z - Set if bit b of register r is 0.
// 	N - Reset.
// 	H - Set.
// 	C - Not affected.
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		BIT 			b,A 			CB 47 		8
// 		BIT 			b,B 			CB 40 		8
// 		BIT 			b,C 			CB 41 		8
// 		BIT 			b,D 			CB 42 		8
// 		BIT 			b,E 			CB 43 		8
// 		BIT 			b,H 			CB 44 		8
// 		BIT 			b,L 			CB 45 		8
// 		BIT 			b,(HL) 			CB 46 		16

const (
	bitACycles     = 8
	bitBCycles     = 8
	bitCCycles     = 8
	bitDCycles     = 8
	bitECycles     = 8
	bitHCycles     = 8
	bitLCycles     = 8
	bitMemHlCycles = 16
)

func bit0A(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.af.a & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitACycles
}
func bit1A(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.af.a & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitACycles
}
func bit2A(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.af.a & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitACycles
}
func bit3A(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.af.a & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitACycles
}
func bit4A(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.af.a & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitACycles
}
func bit5A(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.af.a & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitACycles
}
func bit6A(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.af.a & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitACycles
}
func bit7A(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.af.a & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitACycles
}

func bit0B(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.bc.b & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitBCycles
}
func bit1B(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.bc.b & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitBCycles
}
func bit2B(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.bc.b & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitBCycles
}
func bit3B(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.bc.b & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitBCycles
}
func bit4B(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.bc.b & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitBCycles
}
func bit5B(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.bc.b & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitBCycles
}
func bit6B(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.bc.b & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitBCycles
}
func bit7B(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.bc.b & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitBCycles
}

func bit0C(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.bc.c & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitCCycles
}
func bit1C(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.bc.c & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitCCycles
}
func bit2C(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.bc.c & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitCCycles
}
func bit3C(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.bc.c & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitCCycles
}
func bit4C(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.bc.c & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitCCycles
}
func bit5C(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.bc.c & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitCCycles
}
func bit6C(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.bc.c & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitCCycles
}
func bit7C(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.bc.c & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitCCycles
}

func bit0D(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.de.d & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitDCycles
}
func bit1D(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.de.d & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitDCycles
}
func bit2D(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.de.d & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitDCycles
}
func bit3D(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.de.d & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitDCycles
}
func bit4D(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.de.d & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitDCycles
}
func bit5D(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.de.d & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitDCycles
}
func bit6D(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.de.d & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitDCycles
}
func bit7D(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.de.d & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitDCycles
}

func bit0E(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.de.e & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitECycles
}
func bit1E(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.de.e & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitECycles
}
func bit2E(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.de.e & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitECycles
}
func bit3E(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.de.e & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitECycles
}
func bit4E(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.de.e & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitECycles
}
func bit5E(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.de.e & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitECycles
}
func bit6E(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.de.e & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitECycles
}
func bit7E(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.de.e & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitECycles
}

func bit0H(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.hl.h & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitHCycles
}
func bit1H(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.hl.h & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitHCycles
}
func bit2H(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.hl.h & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitHCycles
}
func bit3H(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.hl.h & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitHCycles
}
func bit4H(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.hl.h & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitHCycles
}
func bit5H(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.hl.h & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitHCycles
}
func bit6H(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.hl.h & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitHCycles
}
func bit7H(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.hl.h & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitHCycles
}

func bit0L(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.hl.l & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitLCycles
}
func bit1L(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.hl.l & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitLCycles
}
func bit2L(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.hl.l & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitLCycles
}
func bit3L(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.hl.l & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitLCycles
}
func bit4L(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.hl.l & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitLCycles
}
func bit5L(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.hl.l & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitLCycles
}
func bit6L(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.hl.l & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitLCycles
}
func bit7L(cpu *cpu) cycleCount {
	cpu.r.setFlagZ((cpu.r.hl.l & 0x01) == 0x00)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH(true)
	return bitLCycles
}

func bit0MemHl(cpu *cpu) cycleCount {
	// TODO: To implement
	return bitMemHlCycles
}
func bit1MemHl(cpu *cpu) cycleCount {
	// TODO: To implement
	return bitMemHlCycles
}
func bit2MemHl(cpu *cpu) cycleCount {
	// TODO: To implement
	return bitMemHlCycles
}
func bit3MemHl(cpu *cpu) cycleCount {
	// TODO: To implement
	return bitMemHlCycles
}
func bit4MemHl(cpu *cpu) cycleCount {
	// TODO: To implement
	return bitMemHlCycles
}
func bit5MemHl(cpu *cpu) cycleCount {
	// TODO: To implement
	return bitMemHlCycles
}
func bit6MemHl(cpu *cpu) cycleCount {
	// TODO: To implement
	return bitMemHlCycles
}
func bit7MemHl(cpu *cpu) cycleCount {
	// TODO: To implement
	return bitMemHlCycles
}

// 3.3.7.2. SET b,r
// Description:
// 	Set bit b in register r.
// Use with:
// 	b = 0 - 7, r = A,B,C,D,E,H,L,(HL)
// Flags affected:
// 	None.
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		SET 			b,A 			CB C7 		8
// 		SET 			b,B 			CB C0 		8
// 		SET 			b,C 			CB C1 		8
// 		SET 			b,D 			CB C2 		8
// 		SET 			b,E 			CB C3 		8
// 		SET 			b,H 			CB C4 		8
// 		SET 			b,L 			CB C5 		8
// 		SET 			b,(HL) 			CB C6 		16

// 		SET 			b,A 			CB C7 		8
// 		SET 			b,B 			CB C0 		8
// 		SET 			b,C 			CB C1 		8
// 		SET 			b,D 			CB C2 		8
// 		SET 			b,E 			CB C3 		8
// 		SET 			b,H 			CB C4 		8
// 		SET 			b,L 			CB C5 		8
// 		SET 			b,MemHl 			CB C6 		16

const (
	setACycles     = 8
	setBCycles     = 8
	setCCycles     = 8
	setDCycles     = 8
	setECycles     = 8
	setHCycles     = 8
	setLCycles     = 8
	setMemHlCycles = 16
)

func set0A(cpu *cpu) cycleCount {
	// Sets the 0 bit of the register A
	cpu.r.af.a = cpu.r.af.a | (0x01 << 0)
	return setACycles
}
func set0B(cpu *cpu) cycleCount {
	// Sets the 0 bit of the register B
	cpu.r.bc.b = cpu.r.bc.b | (0x01 << 0)
	return setBCycles
}
func set0C(cpu *cpu) cycleCount {
	// Sets the 0 bit of the register C
	cpu.r.bc.c = cpu.r.bc.c | (0x01 << 0)
	return setCCycles
}
func set0D(cpu *cpu) cycleCount {
	// Sets the 0 bit of the register D
	cpu.r.de.d = cpu.r.de.d | (0x01 << 0)
	return setDCycles
}
func set0E(cpu *cpu) cycleCount {
	// Sets the 0 bit of the register E
	cpu.r.de.e = cpu.r.de.e | (0x01 << 0)
	return setECycles
}
func set0H(cpu *cpu) cycleCount {
	// Sets the 0 bit of the register H
	cpu.r.hl.h = cpu.r.hl.h | (0x01 << 0)
	return setHCycles
}
func set0L(cpu *cpu) cycleCount {
	// Sets the 0 bit of the register L
	cpu.r.hl.l = cpu.r.hl.l | (0x01 << 0)
	return setLCycles
}

func set1A(cpu *cpu) cycleCount {
	// Sets the 1 bit of the register A
	cpu.r.af.a = cpu.r.af.a | (0x01 << 1)
	return setACycles
}
func set1B(cpu *cpu) cycleCount {
	// Sets the 1 bit of the register B
	cpu.r.bc.b = cpu.r.bc.b | (0x01 << 1)
	return setBCycles
}
func set1C(cpu *cpu) cycleCount {
	// Sets the 1 bit of the register C
	cpu.r.bc.c = cpu.r.bc.c | (0x01 << 1)
	return setCCycles
}
func set1D(cpu *cpu) cycleCount {
	// Sets the 1 bit of the register D
	cpu.r.de.d = cpu.r.de.d | (0x01 << 1)
	return setDCycles
}
func set1E(cpu *cpu) cycleCount {
	// Sets the 1 bit of the register E
	cpu.r.de.e = cpu.r.de.e | (0x01 << 1)
	return setECycles
}
func set1H(cpu *cpu) cycleCount {
	// Sets the 1 bit of the register H
	cpu.r.hl.h = cpu.r.hl.h | (0x01 << 1)
	return setHCycles
}
func set1L(cpu *cpu) cycleCount {
	// Sets the 1 bit of the register L
	cpu.r.hl.l = cpu.r.hl.l | (0x01 << 1)
	return setLCycles
}

func set2A(cpu *cpu) cycleCount {
	// Sets the 2 bit of the register A
	cpu.r.af.a = cpu.r.af.a | (0x01 << 2)
	return setACycles
}
func set2B(cpu *cpu) cycleCount {
	// Sets the 2 bit of the register B
	cpu.r.bc.b = cpu.r.bc.b | (0x01 << 2)
	return setBCycles
}
func set2C(cpu *cpu) cycleCount {
	// Sets the 2 bit of the register C
	cpu.r.bc.c = cpu.r.bc.c | (0x01 << 2)
	return setCCycles
}
func set2D(cpu *cpu) cycleCount {
	// Sets the 2 bit of the register D
	cpu.r.de.d = cpu.r.de.d | (0x01 << 2)
	return setDCycles
}
func set2E(cpu *cpu) cycleCount {
	// Sets the 2 bit of the register E
	cpu.r.de.e = cpu.r.de.e | (0x01 << 2)
	return setECycles
}
func set2H(cpu *cpu) cycleCount {
	// Sets the 2 bit of the register H
	cpu.r.hl.h = cpu.r.hl.h | (0x01 << 2)
	return setHCycles
}
func set2L(cpu *cpu) cycleCount {
	// Sets the 2 bit of the register L
	cpu.r.hl.l = cpu.r.hl.l | (0x01 << 2)
	return setLCycles
}

func set3A(cpu *cpu) cycleCount {
	// Sets the 3 bit of the register A
	cpu.r.af.a = cpu.r.af.a | (0x01 << 3)
	return setACycles
}
func set3B(cpu *cpu) cycleCount {
	// Sets the 3 bit of the register B
	cpu.r.bc.b = cpu.r.bc.b | (0x01 << 3)
	return setBCycles
}
func set3C(cpu *cpu) cycleCount {
	// Sets the 3 bit of the register C
	cpu.r.bc.c = cpu.r.bc.c | (0x01 << 3)
	return setCCycles
}
func set3D(cpu *cpu) cycleCount {
	// Sets the 3 bit of the register D
	cpu.r.de.d = cpu.r.de.d | (0x01 << 3)
	return setDCycles
}
func set3E(cpu *cpu) cycleCount {
	// Sets the 3 bit of the register E
	cpu.r.de.e = cpu.r.de.e | (0x01 << 3)
	return setECycles
}
func set3H(cpu *cpu) cycleCount {
	// Sets the 3 bit of the register H
	cpu.r.hl.h = cpu.r.hl.h | (0x01 << 3)
	return setHCycles
}
func set3L(cpu *cpu) cycleCount {
	// Sets the 3 bit of the register L
	cpu.r.hl.l = cpu.r.hl.l | (0x01 << 3)
	return setLCycles
}

func set4A(cpu *cpu) cycleCount {
	// Sets the 4 bit of the register A
	cpu.r.af.a = cpu.r.af.a | (0x01 << 4)
	return setACycles
}
func set4B(cpu *cpu) cycleCount {
	// Sets the 4 bit of the register B
	cpu.r.bc.b = cpu.r.bc.b | (0x01 << 4)
	return setBCycles
}
func set4C(cpu *cpu) cycleCount {
	// Sets the 4 bit of the register C
	cpu.r.bc.c = cpu.r.bc.c | (0x01 << 4)
	return setCCycles
}
func set4D(cpu *cpu) cycleCount {
	// Sets the 4 bit of the register D
	cpu.r.de.d = cpu.r.de.d | (0x01 << 4)
	return setDCycles
}
func set4E(cpu *cpu) cycleCount {
	// Sets the 4 bit of the register E
	cpu.r.de.e = cpu.r.de.e | (0x01 << 4)
	return setECycles
}
func set4H(cpu *cpu) cycleCount {
	// Sets the 4 bit of the register H
	cpu.r.hl.h = cpu.r.hl.h | (0x01 << 4)
	return setHCycles
}
func set4L(cpu *cpu) cycleCount {
	// Sets the 4 bit of the register L
	cpu.r.hl.l = cpu.r.hl.l | (0x01 << 4)
	return setLCycles
}

func set5A(cpu *cpu) cycleCount {
	// Sets the 5 bit of the register A
	cpu.r.af.a = cpu.r.af.a | (0x01 << 5)
	return setACycles
}
func set5B(cpu *cpu) cycleCount {
	// Sets the 5 bit of the register B
	cpu.r.bc.b = cpu.r.bc.b | (0x01 << 5)
	return setBCycles
}
func set5C(cpu *cpu) cycleCount {
	// Sets the 5 bit of the register C
	cpu.r.bc.c = cpu.r.bc.c | (0x01 << 5)
	return setCCycles
}
func set5D(cpu *cpu) cycleCount {
	// Sets the 5 bit of the register D
	cpu.r.de.d = cpu.r.de.d | (0x01 << 5)
	return setDCycles
}
func set5E(cpu *cpu) cycleCount {
	// Sets the 5 bit of the register E
	cpu.r.de.e = cpu.r.de.e | (0x01 << 5)
	return setECycles
}
func set5H(cpu *cpu) cycleCount {
	// Sets the 5 bit of the register H
	cpu.r.hl.h = cpu.r.hl.h | (0x01 << 5)
	return setHCycles
}
func set5L(cpu *cpu) cycleCount {
	// Sets the 5 bit of the register L
	cpu.r.hl.l = cpu.r.hl.l | (0x01 << 5)
	return setLCycles
}

func set6A(cpu *cpu) cycleCount {
	// Sets the 6 bit of the register A
	cpu.r.af.a = cpu.r.af.a | (0x01 << 6)
	return setACycles
}
func set6B(cpu *cpu) cycleCount {
	// Sets the 6 bit of the register B
	cpu.r.bc.b = cpu.r.bc.b | (0x01 << 6)
	return setBCycles
}
func set6C(cpu *cpu) cycleCount {
	// Sets the 6 bit of the register C
	cpu.r.bc.c = cpu.r.bc.c | (0x01 << 6)
	return setCCycles
}
func set6D(cpu *cpu) cycleCount {
	// Sets the 6 bit of the register D
	cpu.r.de.d = cpu.r.de.d | (0x01 << 6)
	return setDCycles
}
func set6E(cpu *cpu) cycleCount {
	// Sets the 6 bit of the register E
	cpu.r.de.e = cpu.r.de.e | (0x01 << 6)
	return setECycles
}
func set6H(cpu *cpu) cycleCount {
	// Sets the 6 bit of the register H
	cpu.r.hl.h = cpu.r.hl.h | (0x01 << 6)
	return setHCycles
}
func set6L(cpu *cpu) cycleCount {
	// Sets the 6 bit of the register L
	cpu.r.hl.l = cpu.r.hl.l | (0x01 << 6)
	return setLCycles
}

func set7A(cpu *cpu) cycleCount {
	// Sets the 7 bit of the register A
	cpu.r.af.a = cpu.r.af.a | (0x01 << 7)
	return setACycles
}
func set7B(cpu *cpu) cycleCount {
	// Sets the 7 bit of the register B
	cpu.r.bc.b = cpu.r.bc.b | (0x01 << 7)
	return setBCycles
}
func set7C(cpu *cpu) cycleCount {
	// Sets the 7 bit of the register C
	cpu.r.bc.c = cpu.r.bc.c | (0x01 << 7)
	return setCCycles
}
func set7D(cpu *cpu) cycleCount {
	// Sets the 7 bit of the register D
	cpu.r.de.d = cpu.r.de.d | (0x01 << 7)
	return setDCycles
}
func set7E(cpu *cpu) cycleCount {
	// Sets the 7 bit of the register E
	cpu.r.de.e = cpu.r.de.e | (0x01 << 7)
	return setECycles
}
func set7H(cpu *cpu) cycleCount {
	// Sets the 7 bit of the register H
	cpu.r.hl.h = cpu.r.hl.h | (0x01 << 7)
	return setHCycles
}
func set7L(cpu *cpu) cycleCount {
	// Sets the 7 bit of the register L
	cpu.r.hl.l = cpu.r.hl.l | (0x01 << 7)
	return setLCycles
}

func set0MemHl(cpu *cpu) cycleCount {
	// set the 0 bit of the position of memory pointed by register HL
	// TODO: To implement
	return setMemHlCycles
}
func set1MemHl(cpu *cpu) cycleCount {
	// set the 1 bit of the position of memory pointed by register HL
	// TODO: To implement
	return setMemHlCycles
}
func set2MemHl(cpu *cpu) cycleCount {
	// set the 2 bit of the position of memory pointed by register HL
	// TODO: To implement
	return setMemHlCycles
}
func set3MemHl(cpu *cpu) cycleCount {
	// set the 3 bit of the position of memory pointed by register HL
	// TODO: To implement
	return setMemHlCycles
}
func set4MemHl(cpu *cpu) cycleCount {
	// set the 4 bit of the position of memory pointed by register HL
	// TODO: To implement
	return setMemHlCycles
}
func set5MemHl(cpu *cpu) cycleCount {
	// set the 5 bit of the position of memory pointed by register HL
	// TODO: To implement
	return setMemHlCycles
}
func set6MemHl(cpu *cpu) cycleCount {
	// set the 6 bit of the position of memory pointed by register HL
	// TODO: To implement
	return setMemHlCycles
}
func set7MemHl(cpu *cpu) cycleCount {
	// set the 7 bit of the position of memory pointed by register HL
	// TODO: To implement
	return setMemHlCycles
}

// 3.3.7.3. RES b,r
// Description:
// 	Reset bit b in register r.
// Use with:
// 	b = 0 - 7, r = A,B,C,D,E,H,L,(HL)
// Flags affected:
// 	None.
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		RES 			b,A 			CB 87 		8
// 		RES 			b,B 			CB 80 		8
// 		RES 			b,C 			CB 81 		8
// 		RES 			b,D 			CB 82 		8
// 		RES 			b,E 			CB 83 		8
// 		RES 			b,H 			CB 84 		8
// 		RES 			b,L 			CB 85 		8
// 		RES 			b,(HL) 			CB 86 		16

const (
	resACycles     = 8
	resBCycles     = 8
	resCCycles     = 8
	resDCycles     = 8
	resECycles     = 8
	resHCycles     = 8
	resLCycles     = 8
	resMemHlCycles = 16
)

func res0A(cpu *cpu) cycleCount {
	// Reset the 0 bit of the register A
	cpu.r.af.a = cpu.r.af.a & (0xFF ^ (0x01 << 0))
	return resACycles
}
func res0B(cpu *cpu) cycleCount {
	// Reset the 0 bit of the register B
	cpu.r.bc.b = cpu.r.bc.b & (0xFF ^ (0x01 << 0))
	return resBCycles
}
func res0C(cpu *cpu) cycleCount {
	// Reset the 0 bit of the register C
	cpu.r.bc.c = cpu.r.bc.c & (0xFF ^ (0x01 << 0))
	return resCCycles
}
func res0D(cpu *cpu) cycleCount {
	// Reset the 0 bit of the register D
	cpu.r.de.d = cpu.r.de.d & (0xFF ^ (0x01 << 0))
	return resDCycles
}
func res0E(cpu *cpu) cycleCount {
	// Reset the 0 bit of the register E
	cpu.r.de.e = cpu.r.de.e & (0xFF ^ (0x01 << 0))
	return resECycles
}
func res0H(cpu *cpu) cycleCount {
	// Reset the 0 bit of the register H
	cpu.r.hl.h = cpu.r.hl.h & (0xFF ^ (0x01 << 0))
	return resHCycles
}
func res0L(cpu *cpu) cycleCount {
	// Reset the 0 bit of the register L
	cpu.r.hl.l = cpu.r.hl.l & (0xFF ^ (0x01 << 0))
	return resLCycles
}

func res1A(cpu *cpu) cycleCount {
	// Reset the 1 bit of the register A
	cpu.r.af.a = cpu.r.af.a & (0xFF ^ (0x01 << 1))
	return resACycles
}
func res1B(cpu *cpu) cycleCount {
	// Reset the 1 bit of the register B
	cpu.r.bc.b = cpu.r.bc.b & (0xFF ^ (0x01 << 1))
	return resBCycles
}
func res1C(cpu *cpu) cycleCount {
	// Reset the 1 bit of the register C
	cpu.r.bc.c = cpu.r.bc.c & (0xFF ^ (0x01 << 1))
	return resCCycles
}
func res1D(cpu *cpu) cycleCount {
	// Reset the 1 bit of the register D
	cpu.r.de.d = cpu.r.de.d & (0xFF ^ (0x01 << 1))
	return resDCycles
}
func res1E(cpu *cpu) cycleCount {
	// Reset the 1 bit of the register E
	cpu.r.de.e = cpu.r.de.e & (0xFF ^ (0x01 << 1))
	return resECycles
}
func res1H(cpu *cpu) cycleCount {
	// Reset the 1 bit of the register H
	cpu.r.hl.h = cpu.r.hl.h & (0xFF ^ (0x01 << 1))
	return resHCycles
}
func res1L(cpu *cpu) cycleCount {
	// Reset the 1 bit of the register L
	cpu.r.hl.l = cpu.r.hl.l & (0xFF ^ (0x01 << 1))
	return resLCycles
}

func res2A(cpu *cpu) cycleCount {
	// Reset the 2 bit of the register A
	cpu.r.af.a = cpu.r.af.a & (0xFF ^ (0x01 << 2))
	return resACycles
}
func res2B(cpu *cpu) cycleCount {
	// Reset the 2 bit of the register B
	cpu.r.bc.b = cpu.r.bc.b & (0xFF ^ (0x01 << 2))
	return resBCycles
}
func res2C(cpu *cpu) cycleCount {
	// Reset the 2 bit of the register C
	cpu.r.bc.c = cpu.r.bc.c & (0xFF ^ (0x01 << 2))
	return resCCycles
}
func res2D(cpu *cpu) cycleCount {
	// Reset the 2 bit of the register D
	cpu.r.de.d = cpu.r.de.d & (0xFF ^ (0x01 << 2))
	return resDCycles
}
func res2E(cpu *cpu) cycleCount {
	// Reset the 2 bit of the register E
	cpu.r.de.e = cpu.r.de.e & (0xFF ^ (0x01 << 2))
	return resECycles
}
func res2H(cpu *cpu) cycleCount {
	// Reset the 2 bit of the register H
	cpu.r.hl.h = cpu.r.hl.h & (0xFF ^ (0x01 << 2))
	return resHCycles
}
func res2L(cpu *cpu) cycleCount {
	// Reset the 2 bit of the register L
	cpu.r.hl.l = cpu.r.hl.l & (0xFF ^ (0x01 << 2))
	return resLCycles
}

func res3A(cpu *cpu) cycleCount {
	// Reset the 3 bit of the register A
	cpu.r.af.a = cpu.r.af.a & (0xFF ^ (0x01 << 3))
	return resACycles
}
func res3B(cpu *cpu) cycleCount {
	// Reset the 3 bit of the register B
	cpu.r.bc.b = cpu.r.bc.b & (0xFF ^ (0x01 << 3))
	return resBCycles
}
func res3C(cpu *cpu) cycleCount {
	// Reset the 3 bit of the register C
	cpu.r.bc.c = cpu.r.bc.c & (0xFF ^ (0x01 << 3))
	return resCCycles
}
func res3D(cpu *cpu) cycleCount {
	// Reset the 3 bit of the register D
	cpu.r.de.d = cpu.r.de.d & (0xFF ^ (0x01 << 3))
	return resDCycles
}
func res3E(cpu *cpu) cycleCount {
	// Reset the 3 bit of the register E
	cpu.r.de.e = cpu.r.de.e & (0xFF ^ (0x01 << 3))
	return resECycles
}
func res3H(cpu *cpu) cycleCount {
	// Reset the 3 bit of the register H
	cpu.r.hl.h = cpu.r.hl.h & (0xFF ^ (0x01 << 3))
	return resHCycles
}
func res3L(cpu *cpu) cycleCount {
	// Reset the 3 bit of the register L
	cpu.r.hl.l = cpu.r.hl.l & (0xFF ^ (0x01 << 3))
	return resLCycles
}

func res4A(cpu *cpu) cycleCount {
	// Reset the 4 bit of the register A
	cpu.r.af.a = cpu.r.af.a & (0xFF ^ (0x01 << 4))
	return resACycles
}
func res4B(cpu *cpu) cycleCount {
	// Reset the 4 bit of the register B
	cpu.r.bc.b = cpu.r.bc.b & (0xFF ^ (0x01 << 4))
	return resBCycles
}
func res4C(cpu *cpu) cycleCount {
	// Reset the 4 bit of the register C
	cpu.r.bc.c = cpu.r.bc.c & (0xFF ^ (0x01 << 4))
	return resCCycles
}
func res4D(cpu *cpu) cycleCount {
	// Reset the 4 bit of the register D
	cpu.r.de.d = cpu.r.de.d & (0xFF ^ (0x01 << 4))
	return resDCycles
}
func res4E(cpu *cpu) cycleCount {
	// Reset the 4 bit of the register E
	cpu.r.de.e = cpu.r.de.e & (0xFF ^ (0x01 << 4))
	return resECycles
}
func res4H(cpu *cpu) cycleCount {
	// Reset the 4 bit of the register H
	cpu.r.hl.h = cpu.r.hl.h & (0xFF ^ (0x01 << 4))
	return resHCycles
}
func res4L(cpu *cpu) cycleCount {
	// Reset the 4 bit of the register L
	cpu.r.hl.l = cpu.r.hl.l & (0xFF ^ (0x01 << 4))
	return resLCycles
}

func res5A(cpu *cpu) cycleCount {
	// Reset the 5 bit of the register A
	cpu.r.af.a = cpu.r.af.a & (0xFF ^ (0x01 << 5))
	return resACycles
}
func res5B(cpu *cpu) cycleCount {
	// Reset the 5 bit of the register B
	cpu.r.bc.b = cpu.r.bc.b & (0xFF ^ (0x01 << 5))
	return resBCycles
}
func res5C(cpu *cpu) cycleCount {
	// Reset the 5 bit of the register C
	cpu.r.bc.c = cpu.r.bc.c & (0xFF ^ (0x01 << 5))
	return resCCycles
}
func res5D(cpu *cpu) cycleCount {
	// Reset the 5 bit of the register D
	cpu.r.de.d = cpu.r.de.d & (0xFF ^ (0x01 << 5))
	return resDCycles
}
func res5E(cpu *cpu) cycleCount {
	// Reset the 5 bit of the register E
	cpu.r.de.e = cpu.r.de.e & (0xFF ^ (0x01 << 5))
	return resECycles
}
func res5H(cpu *cpu) cycleCount {
	// Reset the 5 bit of the register H
	cpu.r.hl.h = cpu.r.hl.h & (0xFF ^ (0x01 << 5))
	return resHCycles
}
func res5L(cpu *cpu) cycleCount {
	// Reset the 5 bit of the register L
	cpu.r.hl.l = cpu.r.hl.l & (0xFF ^ (0x01 << 5))
	return resLCycles
}

func res6A(cpu *cpu) cycleCount {
	// Reset the 6 bit of the register A
	cpu.r.af.a = cpu.r.af.a & (0xFF ^ (0x01 << 6))
	return resACycles
}
func res6B(cpu *cpu) cycleCount {
	// Reset the 6 bit of the register B
	cpu.r.bc.b = cpu.r.bc.b & (0xFF ^ (0x01 << 6))
	return resBCycles
}
func res6C(cpu *cpu) cycleCount {
	// Reset the 6 bit of the register C
	cpu.r.bc.c = cpu.r.bc.c & (0xFF ^ (0x01 << 6))
	return resCCycles
}
func res6D(cpu *cpu) cycleCount {
	// Reset the 6 bit of the register D
	cpu.r.de.d = cpu.r.de.d & (0xFF ^ (0x01 << 6))
	return resDCycles
}
func res6E(cpu *cpu) cycleCount {
	// Reset the 6 bit of the register E
	cpu.r.de.e = cpu.r.de.e & (0xFF ^ (0x01 << 6))
	return resECycles
}
func res6H(cpu *cpu) cycleCount {
	// Reset the 6 bit of the register H
	cpu.r.hl.h = cpu.r.hl.h & (0xFF ^ (0x01 << 6))
	return resHCycles
}
func res6L(cpu *cpu) cycleCount {
	// Reset the 6 bit of the register L
	cpu.r.hl.l = cpu.r.hl.l & (0xFF ^ (0x01 << 6))
	return resLCycles
}

func res7A(cpu *cpu) cycleCount {
	// Reset the 7 bit of the register A
	cpu.r.af.a = cpu.r.af.a & (0xFF ^ (0x01 << 7))
	return resACycles
}
func res7B(cpu *cpu) cycleCount {
	// Reset the 7 bit of the register B
	cpu.r.bc.b = cpu.r.bc.b & (0xFF ^ (0x01 << 7))
	return resBCycles
}
func res7C(cpu *cpu) cycleCount {
	// Reset the 7 bit of the register C
	cpu.r.bc.c = cpu.r.bc.c & (0xFF ^ (0x01 << 7))
	return resCCycles
}
func res7D(cpu *cpu) cycleCount {
	// Reset the 7 bit of the register D
	cpu.r.de.d = cpu.r.de.d & (0xFF ^ (0x01 << 7))
	return resDCycles
}
func res7E(cpu *cpu) cycleCount {
	// Reset the 7 bit of the register E
	cpu.r.de.e = cpu.r.de.e & (0xFF ^ (0x01 << 7))
	return resECycles
}
func res7H(cpu *cpu) cycleCount {
	// Reset the 7 bit of the register H
	cpu.r.hl.h = cpu.r.hl.h & (0xFF ^ (0x01 << 7))
	return resHCycles
}
func res7L(cpu *cpu) cycleCount {
	// Reset the 7 bit of the register L
	cpu.r.hl.l = cpu.r.hl.l & (0xFF ^ (0x01 << 7))
	return resLCycles
}

func res0MemHl(cpu *cpu) cycleCount {
	// Reset the 0 bit of the position of memory pointed by the register HL
	// TODO: To implement
	return resMemHlCycles
}
func res1MemHl(cpu *cpu) cycleCount {
	// Reset the 1 bit of the position of memory pointed by the register HL
	// TODO: To implement
	return resMemHlCycles
}
func res2MemHl(cpu *cpu) cycleCount {
	// Reset the 2 bit of the position of memory pointed by the register HL
	// TODO: To implement
	return resMemHlCycles
}
func res3MemHl(cpu *cpu) cycleCount {
	// Reset the 3 bit of the position of memory pointed by the register HL
	// TODO: To implement
	return resMemHlCycles
}
func res4MemHl(cpu *cpu) cycleCount {
	// Reset the 4 bit of the position of memory pointed by the register HL
	// TODO: To implement
	return resMemHlCycles
}
func res5MemHl(cpu *cpu) cycleCount {
	// Reset the 5 bit of the position of memory pointed by the register HL
	// TODO: To implement
	return resMemHlCycles
}
func res6MemHl(cpu *cpu) cycleCount {
	// Reset the 6 bit of the position of memory pointed by the register HL
	// TODO: To implement
	return resMemHlCycles
}
func res7MemHl(cpu *cpu) cycleCount {
	// Reset the 7 bit of the position of memory pointed by the register HL
	// TODO: To implement
	return resMemHlCycles
}

// 3.3.8. Jumps

// 3.3.8.1. JP nn
// Description:
// 	Jump to address nn.
// Use with:
// 	nn = two byte immediate value. (LS byte first.)
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		JP 				nn 				C3 			12
func jp(cpu *cpu) cycleCount {
	// Jump to address nn
	// (nn: parameter from immediate value)
	// TODO: To implement
	return jpCycles
}

// 3.3.8.2. JP cc,nn
// Description:
// 	Jump to address n if following condition is true:
// 	cc = NZ, Jump if Z flag is reset.
// 	cc = Z, Jump if Z flag is set.
// 	cc = NC, Jump if C flag is reset.
// 	cc = C, Jump if C flag is set.
// Use with:
// 	nn = two byte immediate value. (LS byte first.)
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		JP 				NZ,nn 			C2 			12
// 		JP 				Z,nn 			CA 			12
// 		JP 				NC,nn 			D2 			12
// 		JP 				C,nn 			DA 			12

func jpNZ(cpu *cpu) cycleCount {
	// Jump to address nn if the flag Z is reset
	// (nn: parameter from immediate value)
	// TODO: To implement
	return jpNZCycles
}
func jpZ(cpu *cpu) cycleCount {
	// Jump to address nn if the flag Z is set
	// (nn: parameter from immediate value)
	// TODO: To implement
	return jpZCycles
}
func jpNC(cpu *cpu) cycleCount {
	// Jump to address nn if the flag C is reset
	// (nn: parameter from immediate value)
	// TODO: To implement
	return jpNCCycles
}
func jpC(cpu *cpu) cycleCount {
	// Jump to address nn if the flag C is set
	// (nn: parameter from immediate value)
	// TODO: To implement
	return jpCCycles
}

// 3.3.8.3. JP (HL)
// Description:
// 	Jump to address contained in HL.
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		JP 				(HL) 			E9 			4

func jpMemHl(cpu *cpu) cycleCount {
	// Jump to address contained in register HL
	// TODO: To implement
	return jpMemHlCycles
}

// 3.3.8.4. JR n
// Description:
// 	Add n to current address and jump to it.
// Use with:
// 	n = one byte signed immediate value
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		JR 				n 				18 			8
func jr(cpu *cpu) cycleCount {
	// Adds current address + nn, and jumps to it
	// (nn: parameter from immediate value)
	// TODO: To implement
	return jrCycles
}

// 3.3.8.5. JR cc,n
// Description:
// 	If following condition is true then add n to current
// 	address and jump to it:
// Use with:
// 	n = one byte signed immediate value
// 	cc = NZ, Jump if Z flag is reset.
// 	cc = Z, Jump if Z flag is set.
// 	cc = NC, Jump if C flag is reset.
// 	cc = C, Jump if C flag is set.
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		JR 				NZ,* 			20 			8
// 		JR 				Z,* 			28 			8
// 		JR 				NC,* 			30 			8
// 		JR 				C,* 			38 			8
func jrNZ(cpu *cpu) cycleCount {
	// If the flag Z is reset, then
	// adds current address + nn, and jumps to it
	// (nn: parameter from immediate value)
	// TODO: To implement
	return jrNZCycles
}
func jrZ(cpu *cpu) cycleCount {
	// If the flag Z is set, then
	// adds current address + nn, and jumps to it
	// (nn: parameter from immediate value)
	// TODO: To implement
	return jrZCycles
}
func jrNC(cpu *cpu) cycleCount {
	// If the flag C is reset, then
	// adds current address + nn, and jumps to it
	// (nn: parameter from immediate value)
	// TODO: To implement
	return jrNCCycles
}
func jrC(cpu *cpu) cycleCount {
	// If the flag C is set, then
	// adds current address + nn, and jumps to it
	// (nn: parameter from immediate value)
	// TODO: To implement
	return jrCCycles
}

// 3.3.9. Calls

// 3.3.9.1. CALL nn
// Description:
// 	Push address of next instruction onto stack and then
// 	jump to address nn.
// Use with:
// 	nn = two byte immediate value. (LS byte first.)
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		CALL 			nn 				CD 			12

func call(cpu *cpu) cycleCount {
	// Push the address of the next instruction onto stack
	// and jump to address nn
	// (nn: parameter from immediate value)
	// TODO: To implement
	return callCycles
}

// 3.3.9.2. CALL cc,nn
// Description:
// 	Call address n if following condition is true:
// 	cc = NZ, Call if Z flag is reset.
// 	cc = Z, Call if Z flag is set.
// 	cc = NC, Call if C flag is reset.
// 	cc = C, Call if C flag is set.
// Use with:
// 	nn = two byte immediate value. (LS byte first.)
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		CALL 			NZ,nn 			C4 			12
// 		CALL 			Z,nn 			CC 			12
// 		CALL 			NC,nn 			D4 			12
// 		CALL 			C,nn 			DC 			12

func callNZ(cpu *cpu) cycleCount {
	// If Z flag is reset, then
	// push the address of the next instruction onto stack
	// and jump to address nn
	// (nn: parameter from immediate value)
	// TODO: To implement
	return callNZCycles
}
func callZ(cpu *cpu) cycleCount {
	// If Z flag is set, then
	// push the address of the next instruction onto stack
	// and jump to address nn
	// (nn: parameter from immediate value)
	// TODO: To implement
	return callZCycles
}
func callNC(cpu *cpu) cycleCount {
	// If C flag is reset, then
	// push the address of the next instruction onto stack
	// and jump to address nn
	// (nn: parameter from immediate value)
	// TODO: To implement
	return callNCCycles
}
func callC(cpu *cpu) cycleCount {
	// If C flag is set, then
	// push the address of the next instruction onto stack
	// and jump to address nn
	// (nn: parameter from immediate value)
	// TODO: To implement
	return callCCycles
}

//3.3.10. Restarts

// 3.3.10.1. RST n
// Description:
// 	Push present address onto stack.
// 	Jump to address $0000 + n.
// Use with:
// 	n = $00,$08,$10,$18,$20,$28,$30,$38
// Opcodes:
// 		Instruction 	Parameters 		Opcode 		Cycles
// 		RST 			00H 			C7 			32
// 		RST 			08H 			CF 			32
// 		RST 			10H 			D7 			32
// 		RST 			18H 			DF 			32
// 		RST 			20H 			E7 			32
// 		RST 			28H 			EF 			32
// 		RST 			30H 			F7 			32
// 		RST 			38H 			FF 			32

func rst00H(cpu *cpu) cycleCount {
	// Push current address to stack
	// and jumps to 0x00
	// TODO: To implement
	return rst00HCycles
}
func rst08H(cpu *cpu) cycleCount {
	// Push current address to stack
	// and jumps to 0x08
	// TODO: To implement
	return rst08HCycles
}
func rst10H(cpu *cpu) cycleCount {
	// Push current address to stack
	// and jumps to 0x10
	// TODO: To implement
	return rst10HCycles
}
func rst18H(cpu *cpu) cycleCount {
	// Push current address to stack
	// and jumps to 0x18
	// TODO: To implement
	return rst18HCycles
}
func rst20H(cpu *cpu) cycleCount {
	// Push current address to stack
	// and jumps to 0x20
	// TODO: To implement
	return rst20HCycles
}
func rst28H(cpu *cpu) cycleCount {
	// Push current address to stack
	// and jumps to 0x28
	// TODO: To implement
	return rst28HCycles
}
func rst30H(cpu *cpu) cycleCount {
	// Push current address to stack
	// and jumps to 0x30
	// TODO: To implement
	return rst30HCycles
}
func rst38H(cpu *cpu) cycleCount {
	// Push current address to stack
	// and jumps to 0x38
	// TODO: To implement
	return rst38HCycles
}
