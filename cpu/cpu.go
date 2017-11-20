// package cpu implements the CPU and Registers type
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
	xZeroThreeCycles  = 0  // 0x03
	xZeroFourCycles   = 0  // 0x04
	xZeroFiveCycles   = 0  // 0x05
	ldBNCycles        = 8  // 0x06
	xZeroSevenCycles  = 0  // 0x07
	ldMemNnSpCycles   = 20 // 0x08
	xZeroNineCycles   = 0  // 0x09
	ldAMemBcCycles    = 8  // 0x0A
	xZeroBCycles      = 0  // 0x0B
	xZeroCCycles      = 0  // 0x0C
	xZeroDCycles      = 0  // 0x0D
	ldCNCycles        = 8  // 0x0E
	xZeroFCycles      = 0  // 0x0F
	xOneZeroCycles    = 0  // 0x10
	ldDeNnCycles      = 12 // 0x11
	xOneTwoCycles     = 0  // 0x12
	xOneThreeCycles   = 0  // 0x13
	xOneFourCycles    = 0  // 0x14
	xOneFiveCycles    = 0  // 0x15
	ldDNCycles        = 8  // 0x16
	xOneSevenCycles   = 0  // 0x17
	xOneEightCycles   = 0  // 0x18
	xOneNineCycles    = 0  // 0x19
	ldAMemDeCycles    = 8  // 0x1A
	xOneBCycles       = 0  // 0x1B
	xOneCCycles       = 0  // 0x1C
	xOneDCycles       = 0  // 0x1D
	ldENCycles        = 8  // 0x1E
	xOneFCycles       = 0  // 0x1F
	xTwoZeroCycles    = 0  // 0x20
	ldHlNnCycles      = 12 // 0x21
	ldiMemHLACycles   = 8  // 0x22
	xTwoThreeCycles   = 0  // 0x23
	xTwoFourCycles    = 0  // 0x24
	xTwoFiveCycles    = 0  // 0x25
	ldHNCycles        = 8  // 0x26
	xTwoSevenCycles   = 0  // 0x27
	xTwoEightCycles   = 0  // 0x28
	xTwoNineCycles    = 0  // 0x29
	ldiAMemHLCycles   = 8  // 0x2A
	xTwoBCycles       = 0  // 0x2B
	xTwoCCycles       = 0  // 0x2C
	xTwoDCycles       = 0  // 0x2D
	ldLNCycles        = 8  // 0x2E
	xTwoFCycles       = 0  // 0x2F
	xThreeZeroCycles  = 0  // 0x30
	ldSpNnCycles      = 12 // 0x31
	lddMemHLACycles   = 8  // 0x32
	xThreeThreeCycles = 0  // 0x33
	xThreeFourCycles  = 0  // 0x34
	xThreeFiveCycles  = 0  // 0x35
	ldMemHlNCycles    = 12 // 0x36
	xThreeSevenCycles = 0  // 0x37
	xThreeEightCycles = 0  // 0x38
	xThreeNineCycles  = 0  // 0x39
	lddAMemHLCycles   = 8  // 0x3A
	xThreeBCycles     = 0  // 0x3B
	xThreeCCycles     = 0  // 0x3C
	xThreeDCycles     = 0  // 0x3D
	ldANCycles        = 0  // 0x3E
	xThreeFCycles     = 0  // 0x3F
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
	xSevenSixCycles   = 0  // 0x76
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
	xNineEightCycles  = 0  // 0x98
	xNineNineCycles   = 0  // 0x99
	xNineACycles      = 0  // 0x9A
	xNineBCycles      = 0  // 0x9B
	xNineCCycles      = 0  // 0x9C
	xNineDCycles      = 0  // 0x9D
	xNineECycles      = 0  // 0x9E
	xNineFCycles      = 0  // 0x9F
	xAZeroCycles      = 0  // 0xA0
	xAOneCycles       = 0  // 0xA1
	xATwoCycles       = 0  // 0xA2
	xAThreeCycles     = 0  // 0xA3
	xAFourCycles      = 0  // 0xA4
	xAFiveCycles      = 0  // 0xA5
	xASixCycles       = 0  // 0xA6
	xASevenCycles     = 0  // 0xA7
	xAEightCycles     = 0  // 0xA8
	xANineCycles      = 0  // 0xA9
	xAACycles         = 0  // 0xAA
	xABCycles         = 0  // 0xAB
	xACCycles         = 0  // 0xAC
	xADCycles         = 0  // 0xAD
	xAECycles         = 0  // 0xAE
	xAFCycles         = 0  // 0xAF
	xBZeroCycles      = 0  // 0xB0
	xBOneCycles       = 0  // 0xB1
	xBTwoCycles       = 0  // 0xB2
	xBThreeCycles     = 0  // 0xB3
	xBFourCycles      = 0  // 0xB4
	xBFiveCycles      = 0  // 0xB5
	xBSixCycles       = 0  // 0xB6
	xBSevenCycles     = 0  // 0xB7
	xBEightCycles     = 0  // 0xB8
	xBNineCycles      = 0  // 0xB9
	xBACycles         = 0  // 0xBA
	xBBCycles         = 0  // 0xBB
	xBCCycles         = 0  // 0xBC
	xBDCycles         = 0  // 0xBD
	xBECycles         = 0  // 0xBE
	xBFCycles         = 0  // 0xBF
	xCZeroCycles      = 0  // 0xC0
	popBcCycles       = 12 // 0xC1
	xCTwoCycles       = 0  // 0xC2
	xCThreeCycles     = 0  // 0xC3
	xCFourCycles      = 0  // 0xC4
	pushBcCycles      = 16 // 0xC5
	addANnCycles      = 8  // 0xC6
	xCSevenCycles     = 0  // 0xC7
	xCEightCycles     = 0  // 0xC8
	xCNineCycles      = 0  // 0xC9
	xCACycles         = 0  // 0xCA
	xCBCycles         = 0  // 0xCB
	xCCCycles         = 0  // 0xCC
	xCDCycles         = 0  // 0xCD
	xCECycles         = 0  // 0xCE
	xCFCycles         = 0  // 0xCF
	xDZeroCycles      = 0  // 0xD0
	popDeCycles       = 12 // 0xD1
	xDTwoCycles       = 0  // 0xD2
	xDThreeCycles     = 0  // 0xD3
	xDFourCycles      = 0  // 0xD4
	pushDeCycles      = 16 // 0xD5
	subANnCycles      = 8  // 0xD6
	xDSevenCycles     = 0  // 0xD7
	xDEightCycles     = 0  // 0xD8
	xDNineCycles      = 0  // 0xD9
	xDACycles         = 0  // 0xDA
	xDBCycles         = 0  // 0xDB
	xDCCycles         = 0  // 0xDC
	xDDCycles         = 0  // 0xDD
	xDECycles         = 0  // 0xDE
	xDFCycles         = 0  // 0xDF
	ldStackNACycles   = 12 // 0xE0
	popHlCycles       = 12 // 0xE1
	ldStackCACycles   = 8  // 0xE2
	xEThreeCycles     = 0  // 0xE3
	xEFourCycles      = 0  // 0xE4
	pushHlCycles      = 16 // 0xE5
	xESixCycles       = 0  // 0xE6
	xESevenCycles     = 0  // 0xE7
	xEEightCycles     = 0  // 0xE8
	xENineCycles      = 0  // 0xE9
	xEACycles         = 0  // 0xEA
	xEBCycles         = 0  // 0xEB
	xECCycles         = 0  // 0xEC
	xEDCycles         = 0  // 0xED
	xEECycles         = 0  // 0xEE
	xEFCycles         = 0  // 0xEF
	ldAStackNCycles   = 12 // 0xF0
	popAfCycles       = 12 // 0xF1
	ldAStackCCycles   = 8  // 0xF2
	xFThreeCycles     = 0  // 0xF3
	xFFourCycles      = 0  // 0xF4
	pushAfCycles      = 16 // 0xF5
	xFSixCycles       = 0  // 0xF6
	xFSevenCycles     = 0  // 0xF7
	ldHlSpNCycles     = 12 // 0xF8
	ldSpHlCycles      = 0  // 0xF9
	ldAMemNnCycles    = 16 // 0xFA
	xFBCycles         = 0  // 0xFB
	xFCCycles         = 0  // 0xFC
	xFDCycles         = 0  // 0xFD
	xFECycles         = 0  // 0xFE
	xFFCycles         = 0  // 0xFF
)

var op = [0x100] instructions{
	nop,       //0x00
	ldBcNn,    //0x01
	TODO,      //0x02
	TODO,      //0x03
	TODO,      //0x04
	TODO,      //0x05
	ldBN,      //0x06
	TODO,      //0x07
	ldMemNnSp, //0x08
	TODO,      //0x09
	ldAMemBc,  //0x0A
	TODO,      //0x0B
	TODO,      //0x0C
	TODO,      //0x0D
	ldCN,      //0x0E
	TODO,      //0x0F
	TODO,      //0x10
	ldDeNn,    //0x11
	TODO,      //0x12
	TODO,      //0x13
	TODO,      //0x14
	TODO,      //0x15
	ldDN,      //0x16
	TODO,      //0x17
	TODO,      //0x18
	TODO,      //0x19
	ldAMemDe,  //0x1A
	TODO,      //0x1B
	TODO,      //0x1C
	TODO,      //0x1D
	ldEN,      //0x1E
	TODO,      //0x1F
	TODO,      //0x20
	ldHlNn,    //0x21
	ldiMemHLA, //0x22
	TODO,      //0x23
	TODO,      //0x24
	TODO,      //0x25
	ldHN,      //0x26
	TODO,      //0x27
	TODO,      //0x28
	TODO,      //0x29
	ldiAMemHL, //0x2A
	TODO,      //0x2B
	TODO,      //0x2C
	TODO,      //0x2D
	ldLN,      //0x2E
	TODO,      //0x2F
	TODO,      //0x30
	ldSpNn,    //0x31
	lddMemHLA, //0x32
	TODO,      //0x33
	TODO,      //0x34
	TODO,      //0x35
	ldMemHlN,  //0x36
	TODO,      //0x37
	TODO,      //0x38
	TODO,      //0x39
	lddAMemHL, //0x3A
	TODO,      //0x3B
	TODO,      //0x3C
	TODO,      //0x3D
	ldAN,      //0x3E
	TODO,      //0x3F
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
	TODO,      //0x76
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
	TODO,      //0x98
	TODO,      //0x99
	TODO,      //0x9A
	TODO,      //0x9B
	TODO,      //0x9C
	TODO,      //0x9D
	TODO,      //0x9E
	TODO,      //0x9F
	TODO,      //0xA0
	TODO,      //0xA1
	TODO,      //0xA2
	TODO,      //0xA3
	TODO,      //0xA4
	TODO,      //0xA5
	TODO,      //0xA6
	TODO,      //0xA7
	TODO,      //0xA8
	TODO,      //0xA9
	TODO,      //0xAA
	TODO,      //0xAB
	TODO,      //0xAC
	TODO,      //0xAD
	TODO,      //0xAE
	TODO,      //0xAF
	TODO,      //0xB0
	TODO,      //0xB1
	TODO,      //0xB2
	TODO,      //0xB3
	TODO,      //0xB4
	TODO,      //0xB5
	TODO,      //0xB6
	TODO,      //0xB7
	TODO,      //0xB8
	TODO,      //0xB9
	TODO,      //0xBA
	TODO,      //0xBB
	TODO,      //0xBC
	TODO,      //0xBD
	TODO,      //0xBE
	TODO,      //0xBF
	TODO,      //0xC0
	popBc,     //0xC1
	TODO,      //0xC2
	TODO,      //0xC3
	TODO,      //0xC4
	pushBc,    //0xC5
	addANn,    //0xC6
	TODO,      //0xC7
	TODO,      //0xC8
	TODO,      //0xC9
	TODO,      //0xCA
	TODO,      //0xCB
	TODO,      //0xCC
	TODO,      //0xCD
	TODO,      //0xCE
	TODO,      //0xCF
	TODO,      //0xD0
	popDe,     //0xD1
	TODO,      //0xD2
	TODO,      //0xD3
	TODO,      //0xD4
	pushDe,    //0xD5
	subANn,    //0xD6
	TODO,      //0xD7
	TODO,      //0xD8
	TODO,      //0xD9
	TODO,      //0xDA
	TODO,      //0xDB
	TODO,      //0xDC
	TODO,      //0xDD
	TODO,      //0xDE
	TODO,      //0xDF
	ldStackNA, //0xE0
	popHl,     //0xE1
	ldStackCA, //0xE2
	TODO,      //0xE3
	TODO,      //0xE4
	pushHl,    //0xE5
	TODO,      //0xE6
	TODO,      //0xE7
	TODO,      //0xE8
	TODO,      //0xE9
	TODO,      //0xEA
	TODO,      //0xEB
	TODO,      //0xEC
	TODO,      //0xED
	TODO,      //0xEE
	TODO,      //0xEF
	ldAStackN, //0xF0
	popAf,     //0xF1
	ldAStackC, //0xF2
	TODO,      //0xF3
	TODO,      //0xF4
	pushAf,    //0xF5
	TODO,      //0xF6
	TODO,      //0xF7
	ldHlSpN,   //0xF8
	ldSpHl,    //0xF9
	ldAMemNn,  //0xFA
	TODO,      //0xFB
	TODO,      //0xFC
	TODO,      //0xFD
	TODO,      //0xFE
	TODO,      //0xFF
}

func TODO(cpu *cpu) cycleCount {
	// This function is not defined!
	return 0
}

func nop(cpu *cpu) cycleCount {
	// Does nothing
	return nopCycles
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

func lddAMemHL(cpu *cpu) cycleCount {
	// Put the value from the position of memory in HL into the register A.
	// Then, decrement HL.
	// TODO: to implement
	return lddAMemHLCycles
}

// 3.3.1.10. LD (HLD), A
// Description: Same as: LDD (HL),A

// 3.3.1.11. LD (HLD), A
// Description: Same as: LDD (HL),A

// 3.3.1.12. LDD (HL), A
// Description:
// 		Put A into memory address HL. Decrement HL.
// Same as: LD (HL),A - DEC HL

func lddMemHLA(cpu *cpu) cycleCount {
	// Put value of the register A into the memory address pointed by HL.
	// Then, decrement HL.
	// TODO: to implement
	return lddMemHLACycles
}

// 3.3.1.13. LD A,(HLI)
// Description: Same as: LDI A,(HL)

// 3.3.1.14. LD A,(HL+)
// Description: Same as: LDI A,(HL)

// 3.3.1.15. LDI A,(HL)
// Description:
//		Put value at address HL into A. Increment HL.
// Same as: LD A,(HL) - INC HL

func ldiAMemHL(cpu *cpu) cycleCount {
	// Put the value from the position of memory pointed by HL into the register A.
	// Then, increment  HL.
	// TODO: to implement
	return ldiAMemHLCycles
}

// 3.3.1.16. LD (HLI),A
// Description: Same as: LDI (HL),A

// 3.3.1.17. LD (HL+),A
// Description: Same as: LDI (HL),A

// 3.3.1.18. LDI (HL),A
// Description:
//		Put A into memory address HL. Increment HL.
// Same as: LD (HL),A - INC HL

func ldiMemHLA(cpu *cpu) cycleCount {
	// Put value of the register A into the memory address pointed by HL.
	// Then, increment HL.
	// TODO: to implement
	return ldiMemHLACycles
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
	cpu.r.setFlagH((cpu.r.af.a & (1 << 4)) == 1)
	cpu.r.setFlagC((cpu.r.af.a < oldA))
	return addAACycles
}

func addAB(cpu *cpu) cycleCount {
	// Add the value of register B into register A
	var oldA = cpu.r.af.a
	cpu.r.af.a += cpu.r.bc.b
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH((cpu.r.af.a & (1 << 4)) == 1)
	cpu.r.setFlagC((cpu.r.af.a < oldA))
	return addABCycles
}

func addAC(cpu *cpu) cycleCount {
	// Add the value of register C into register A
	var oldA = cpu.r.af.a
	cpu.r.af.a += cpu.r.bc.c
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH((cpu.r.af.a & (1 << 4)) == 1)
	cpu.r.setFlagC((cpu.r.af.a < oldA))
	return addACCycles
}

func addAD(cpu *cpu) cycleCount {
	// Add the value of register D into register A
	var oldA = cpu.r.af.a
	cpu.r.af.a += cpu.r.de.d
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH((cpu.r.af.a & (1 << 4)) == 1)
	cpu.r.setFlagC((cpu.r.af.a < oldA))
	return addADCycles
}

func addAE(cpu *cpu) cycleCount {
	// Add the value of register E into register A
	var oldA = cpu.r.af.a
	cpu.r.af.a += cpu.r.de.e
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH((cpu.r.af.a & (1 << 4)) == 1)
	cpu.r.setFlagC((cpu.r.af.a < oldA))
	return addAECycles
}

func addAH(cpu *cpu) cycleCount {
	// Add the value of register H into register A
	var oldA = cpu.r.af.a
	cpu.r.af.a += cpu.r.hl.h
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH((cpu.r.af.a & (1 << 4)) == 1)
	cpu.r.setFlagC((cpu.r.af.a < oldA))
	return addAHCycles
}

func addAL(cpu *cpu) cycleCount {
	// Add the value of register L into register A
	var oldA = cpu.r.af.a
	cpu.r.af.a += cpu.r.hl.l
	cpu.r.setFlagZ(cpu.r.af.a == 0)
	cpu.r.setFlagN(false)
	cpu.r.setFlagH((cpu.r.af.a & (1 << 4)) == 1)
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
