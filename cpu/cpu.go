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

type Registers struct {
	af struct {
		a byte  // high
		f Flags // low - flags - not accessible by programmer
	}
	bc struct {
		b byte // high
		c byte // low
	}
	de struct {
		d byte // high
		e byte // low
	}
	hl struct {
		h byte // high
		l byte // low
	}
	sp int // stack pointer
	pc int // program counter
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
	nopSize         = 4  // 0x00
	xZeroOneSize    = 0  // 0x01
	xZeroTwoSize    = 0  // 0x02
	xZeroThreeSize  = 0  // 0x03
	xZeroFourSize   = 0  // 0x04
	xZeroFiveSize   = 0  // 0x05
	ld_b_nSize      = 8  // 0x06
	xZeroSevenSize  = 0  // 0x07
	xZeroEightSize  = 0  // 0x08
	xZeroNineSize   = 0  // 0x09
	xZeroASize      = 0  // 0x0A
	xZeroBSize      = 0  // 0x0B
	xZeroCSize      = 0  // 0x0C
	xZeroDSize      = 0  // 0x0D
	ld_c_nSize      = 8  // 0x0E
	xZeroFSize      = 0  // 0x0F
	xOneZeroSize    = 0  // 0x10
	xOneOneSize     = 0  // 0x11
	xOneTwoSize     = 0  // 0x12
	xOneThreeSize   = 0  // 0x13
	xOneFourSize    = 0  // 0x14
	xOneFiveSize    = 0  // 0x15
	ld_d_nSize      = 8  // 0x16
	xOneSevenSize   = 0  // 0x17
	xOneEightSize   = 0  // 0x18
	xOneNineSize    = 0  // 0x19
	xOneASize       = 0  // 0x1A
	xOneBSize       = 0  // 0x1B
	xOneCSize       = 0  // 0x1C
	xOneDSize       = 0  // 0x1D
	ld_e_nSize      = 8  // 0x1E
	xOneFSize       = 0  // 0x1F
	xTwoZeroSize    = 0  // 0x20
	xTwoOneSize     = 0  // 0x21
	xTwoTwoSize     = 0  // 0x22
	xTwoThreeSize   = 0  // 0x23
	xTwoFourSize    = 0  // 0x24
	xTwoFiveSize    = 0  // 0x25
	ld_h_nSize      = 8  // 0x26
	xTwoSevenSize   = 0  // 0x27
	xTwoEightSize   = 0  // 0x28
	xTwoNineSize    = 0  // 0x29
	xTwoASize       = 0  // 0x2A
	xTwoBSize       = 0  // 0x2B
	xTwoCSize       = 0  // 0x2C
	xTwoDSize       = 0  // 0x2D
	ld_l_nSize      = 8  // 0x2E
	xTwoFSize       = 0  // 0x2F
	xThreeZeroSize  = 0  // 0x30
	xThreeOneSize   = 0  // 0x31
	xThreeTwoSize   = 0  // 0x32
	xThreeThreeSize = 0  // 0x33
	xThreeFourSize  = 0  // 0x34
	xThreeFiveSize  = 0  // 0x35
	ld_hl_nSize     = 12 // 0x36
	xThreeSevenSize = 0  // 0x37
	xThreeEightSize = 0  // 0x38
	xThreeNineSize  = 0  // 0x39
	xThreeASize     = 0  // 0x3A
	xThreeBSize     = 0  // 0x3B
	xThreeCSize     = 0  // 0x3C
	xThreeDSize     = 0  // 0x3D
	xThreeESize     = 0  // 0x3E
	xThreeFSize     = 0  // 0x3F
	ld_b_bSize      = 4  // 0x40
	ld_b_cSize      = 4  // 0x41
	ld_b_dSize      = 4  // 0x42
	ld_b_eSize      = 4  // 0x43
	ld_b_hSize      = 4  // 0x44
	ld_b_lSize      = 4  // 0x45
	ld_b_hlSize     = 8  // 0x46
	xFourSevenSize  = 0  // 0x47
	ld_c_bSize      = 4  // 0x48
	ld_c_cSize      = 4  // 0x49
	ld_c_dSize      = 4  // 0x4A
	ld_c_eSize      = 4  // 0x4B
	ld_c_hSize      = 4  // 0x4C
	ld_c_lSize      = 4  // 0x4D
	ld_c_hlSize     = 8  // 0x4E
	xFourFSize      = 0  // 0x4F
	ld_d_bSize      = 4  // 0x50
	ld_d_cSize      = 4  // 0x51
	ld_d_dSize      = 4  // 0x52
	ld_d_eSize      = 4  // 0x53
	ld_d_hSize      = 4  // 0x54
	ld_d_lSize      = 4  // 0x55
	ld_d_hlSize     = 8  // 0x56
	xFiveSevenSize  = 0  // 0x57
	ld_e_bSize      = 4  // 0x58
	ld_e_cSize      = 4  // 0x59
	ld_e_dSize      = 4  // 0x5A
	ld_e_eSize      = 4  // 0x5B
	ld_e_hSize      = 4  // 0x5C
	ld_e_lSize      = 4  // 0x5D
	ld_e_hlSize     = 8  // 0x5E
	xFiveFSize      = 0  // 0x5F
	ld_h_bSize      = 4  // 0x60
	ld_h_cSize      = 4  // 0x61
	ld_h_dSize      = 4  // 0x62
	ld_h_eSize      = 4  // 0x63
	ld_h_hSize      = 4  // 0x64
	ld_h_lSize      = 4  // 0x65
	ld_h_hlSize     = 8  // 0x66
	xSixSevenSize   = 0  // 0x67
	ld_l_bSize      = 4  // 0x68
	ld_l_cSize      = 4  // 0x69
	ld_l_dSize      = 4  // 0x6A
	ld_l_eSize      = 4  // 0x6B
	ld_l_hSize      = 4  // 0x6C
	ld_l_lSize      = 4  // 0x6D
	ld_l_hlSize     = 8  // 0x6E
	xSixFSize       = 0  // 0x6F
	ld_hl_bSize     = 8  // 0x70
	ld_hl_cSize     = 8  // 0x71
	ld_hl_dSize     = 8  // 0x72
	ld_hl_eSize     = 8  // 0x73
	ld_hl_hSize     = 8  // 0x74
	ld_hl_lSize     = 8  // 0x75
	xSevenSixSize   = 0  // 0x76
	xSevenSevenSize = 0  // 0x77
	ld_a_bSize      = 4  // 0x78
	ld_a_cSize      = 4  // 0x79
	ld_a_dSize      = 4  // 0x7A
	ld_a_eSize      = 4  // 0x7B
	ld_a_hSize      = 4  // 0x7C
	ld_a_lSize      = 4  // 0x7D
	ld_a_hlSize     = 8  // 0x7E
	ld_a_aSize      = 4  // 0x7F
	xEightZeroSize  = 0  // 0x80
	xEightOneSize   = 0  // 0x81
	xEightTwoSize   = 0  // 0x82
	xEightThreeSize = 0  // 0x83
	xEightFourSize  = 0  // 0x84
	xEightFiveSize  = 0  // 0x85
	xEightSixSize   = 0  // 0x86
	xEightSevenSize = 0  // 0x87
	xEightEightSize = 0  // 0x88
	xEightNineSize  = 0  // 0x89
	xEightASize     = 0  // 0x8A
	xEightBSize     = 0  // 0x8B
	xEightCSize     = 0  // 0x8C
	xEightDSize     = 0  // 0x8D
	xEightESize     = 0  // 0x8E
	xEightFSize     = 0  // 0x8F
	xNineZeroSize   = 0  // 0x90
	xNineOneSize    = 0  // 0x91
	xNineTwoSize    = 0  // 0x92
	xNineThreeSize  = 0  // 0x93
	xNineFourSize   = 0  // 0x94
	xNineFiveSize   = 0  // 0x95
	xNineSixSize    = 0  // 0x96
	xNineSevenSize  = 0  // 0x97
	xNineEightSize  = 0  // 0x98
	xNineNineSize   = 0  // 0x99
	xNineASize      = 0  // 0x9A
	xNineBSize      = 0  // 0x9B
	xNineCSize      = 0  // 0x9C
	xNineDSize      = 0  // 0x9D
	xNineESize      = 0  // 0x9E
	xNineFSize      = 0  // 0x9F
	xAZeroSize      = 0  // 0xA0
	xAOneSize       = 0  // 0xA1
	xATwoSize       = 0  // 0xA2
	xAThreeSize     = 0  // 0xA3
	xAFourSize      = 0  // 0xA4
	xAFiveSize      = 0  // 0xA5
	xASixSize       = 0  // 0xA6
	xASevenSize     = 0  // 0xA7
	xAEightSize     = 0  // 0xA8
	xANineSize      = 0  // 0xA9
	xAASize         = 0  // 0xAA
	xABSize         = 0  // 0xAB
	xACSize         = 0  // 0xAC
	xADSize         = 0  // 0xAD
	xAESize         = 0  // 0xAE
	xAFSize         = 0  // 0xAF
	xBZeroSize      = 0  // 0xB0
	xBOneSize       = 0  // 0xB1
	xBTwoSize       = 0  // 0xB2
	xBThreeSize     = 0  // 0xB3
	xBFourSize      = 0  // 0xB4
	xBFiveSize      = 0  // 0xB5
	xBSixSize       = 0  // 0xB6
	xBSevenSize     = 0  // 0xB7
	xBEightSize     = 0  // 0xB8
	xBNineSize      = 0  // 0xB9
	xBASize         = 0  // 0xBA
	xBBSize         = 0  // 0xBB
	xBCSize         = 0  // 0xBC
	xBDSize         = 0  // 0xBD
	xBESize         = 0  // 0xBE
	xBFSize         = 0  // 0xBF
	xCZeroSize      = 0  // 0xC0
	xCOneSize       = 0  // 0xC1
	xCTwoSize       = 0  // 0xC2
	xCThreeSize     = 0  // 0xC3
	xCFourSize      = 0  // 0xC4
	xCFiveSize      = 0  // 0xC5
	xCSixSize       = 0  // 0xC6
	xCSevenSize     = 0  // 0xC7
	xCEightSize     = 0  // 0xC8
	xCNineSize      = 0  // 0xC9
	xCASize         = 0  // 0xCA
	xCBSize         = 0  // 0xCB
	xCCSize         = 0  // 0xCC
	xCDSize         = 0  // 0xCD
	xCESize         = 0  // 0xCE
	xCFSize         = 0  // 0xCF
	xDZeroSize      = 0  // 0xD0
	xDOneSize       = 0  // 0xD1
	xDTwoSize       = 0  // 0xD2
	xDThreeSize     = 0  // 0xD3
	xDFourSize      = 0  // 0xD4
	xDFiveSize      = 0  // 0xD5
	xDSixSize       = 0  // 0xD6
	xDSevenSize     = 0  // 0xD7
	xDEightSize     = 0  // 0xD8
	xDNineSize      = 0  // 0xD9
	xDASize         = 0  // 0xDA
	xDBSize         = 0  // 0xDB
	xDCSize         = 0  // 0xDC
	xDDSize         = 0  // 0xDD
	xDESize         = 0  // 0xDE
	xDFSize         = 0  // 0xDF
	xEZeroSize      = 0  // 0xE0
	xEOneSize       = 0  // 0xE1
	xETwoSize       = 0  // 0xE2
	xEThreeSize     = 0  // 0xE3
	xEFourSize      = 0  // 0xE4
	xEFiveSize      = 0  // 0xE5
	xESixSize       = 0  // 0xE6
	xESevenSize     = 0  // 0xE7
	xEEightSize     = 0  // 0xE8
	xENineSize      = 0  // 0xE9
	xEASize         = 0  // 0xEA
	xEBSize         = 0  // 0xEB
	xECSize         = 0  // 0xEC
	xEDSize         = 0  // 0xED
	xEESize         = 0  // 0xEE
	xEFSize         = 0  // 0xEF
	xFZeroSize      = 0  // 0xF0
	xFOneSize       = 0  // 0xF1
	xFTwoSize       = 0  // 0xF2
	xFThreeSize     = 0  // 0xF3
	xFFourSize      = 0  // 0xF4
	xFFiveSize      = 0  // 0xF5
	xFSixSize       = 0  // 0xF6
	xFSevenSize     = 0  // 0xF7
	xFEightSize     = 0  // 0xF8
	xFNineSize      = 0  // 0xF9
	xFASize         = 0  // 0xFA
	xFBSize         = 0  // 0xFB
	xFCSize         = 0  // 0xFC
	xFDSize         = 0  // 0xFD
	xFESize         = 0  // 0xFE
	xFFSize         = 0  // 0xFF
)

var op = [0x100] instructions{
	nop,     //0x00
	TODO,    //0x01
	TODO,    //0x02
	TODO,    //0x03
	TODO,    //0x04
	TODO,    //0x05
	ld_b_n,  //0x06
	TODO,    //0x07
	TODO,    //0x08
	TODO,    //0x09
	TODO,    //0x0A
	TODO,    //0x0B
	TODO,    //0x0C
	TODO,    //0x0D
	ld_c_n,  //0x0E
	TODO,    //0x0F
	TODO,    //0x10
	TODO,    //0x11
	TODO,    //0x12
	TODO,    //0x13
	TODO,    //0x14
	TODO,    //0x15
	ld_d_n,  //0x16
	TODO,    //0x17
	TODO,    //0x18
	TODO,    //0x19
	TODO,    //0x1A
	TODO,    //0x1B
	TODO,    //0x1C
	TODO,    //0x1D
	ld_e_n,  //0x1E
	TODO,    //0x1F
	TODO,    //0x20
	TODO,    //0x21
	TODO,    //0x22
	TODO,    //0x23
	TODO,    //0x24
	TODO,    //0x25
	ld_h_n,  //0x26
	TODO,    //0x27
	TODO,    //0x28
	TODO,    //0x29
	TODO,    //0x2A
	TODO,    //0x2B
	TODO,    //0x2C
	TODO,    //0x2D
	ld_l_n,  //0x2E
	TODO,    //0x2F
	TODO,    //0x30
	TODO,    //0x31
	TODO,    //0x32
	TODO,    //0x33
	TODO,    //0x34
	TODO,    //0x35
	ld_hl_n, //0x36
	TODO,    //0x37
	TODO,    //0x38
	TODO,    //0x39
	TODO,    //0x3A
	TODO,    //0x3B
	TODO,    //0x3C
	TODO,    //0x3D
	TODO,    //0x3E
	TODO,    //0x3F
	ld_b_b,  //0x40
	ld_b_c,  //0x41
	ld_b_d,  //0x42
	ld_b_e,  //0x43
	ld_b_h,  //0x44
	ld_b_l,  //0x45
	ld_b_hl, //0x46
	TODO,    //0x47
	ld_c_b,  //0x48
	ld_c_c,  //0x49
	ld_c_d,  //0x4A
	ld_c_e,  //0x4B
	ld_c_h,  //0x4C
	ld_c_l,  //0x4D
	ld_c_hl, //0x4E
	TODO,    //0x4F
	ld_d_b,  //0x50
	ld_d_c,  //0x51
	ld_d_d,  //0x52
	ld_d_e,  //0x53
	ld_d_h,  //0x54
	ld_d_l,  //0x55
	ld_d_hl, //0x56
	TODO,    //0x57
	ld_e_b,  //0x58
	ld_e_c,  //0x59
	ld_e_d,  //0x5A
	ld_e_e,  //0x5B
	ld_e_h,  //0x5C
	ld_e_l,  //0x5D
	ld_e_hl, //0x5E
	TODO,    //0x5F
	ld_h_b,  //0x60
	ld_h_c,  //0x61
	ld_h_d,  //0x62
	ld_h_e,  //0x63
	ld_h_h,  //0x64
	ld_h_l,  //0x65
	ld_h_hl, //0x66
	TODO,    //0x67
	ld_l_b,  //0x68
	ld_l_c,  //0x69
	ld_l_d,  //0x6A
	ld_l_e,  //0x6B
	ld_l_h,  //0x6C
	ld_l_l,  //0x6D
	ld_l_hl, //0x6E
	TODO,    //0x6F
	ld_hl_b, //0x70
	ld_hl_c, //0x71
	ld_hl_d, //0x72
	ld_hl_e, //0x73
	ld_hl_l, //0x74
	ld_hl_h, //0x75
	TODO,    //0x76
	TODO,    //0x77
	ld_a_b,  //0x78
	ld_a_c,  //0x79
	ld_a_d,  //0x7A
	ld_a_e,  //0x7B
	ld_a_h,  //0x7C
	ld_a_l,  //0x7D
	ld_a_hl, //0x7E
	ld_a_a,  //0x7F
	TODO,    //0x80
	TODO,    //0x81
	TODO,    //0x82
	TODO,    //0x83
	TODO,    //0x84
	TODO,    //0x85
	TODO,    //0x86
	TODO,    //0x87
	TODO,    //0x88
	TODO,    //0x89
	TODO,    //0x8A
	TODO,    //0x8B
	TODO,    //0x8C
	TODO,    //0x8D
	TODO,    //0x8E
	TODO,    //0x8F
	TODO,    //0x90
	TODO,    //0x91
	TODO,    //0x92
	TODO,    //0x93
	TODO,    //0x94
	TODO,    //0x95
	TODO,    //0x96
	TODO,    //0x97
	TODO,    //0x98
	TODO,    //0x99
	TODO,    //0x9A
	TODO,    //0x9B
	TODO,    //0x9C
	TODO,    //0x9D
	TODO,    //0x9E
	TODO,    //0x9F
	TODO,    //0xA0
	TODO,    //0xA1
	TODO,    //0xA2
	TODO,    //0xA3
	TODO,    //0xA4
	TODO,    //0xA5
	TODO,    //0xA6
	TODO,    //0xA7
	TODO,    //0xA8
	TODO,    //0xA9
	TODO,    //0xAA
	TODO,    //0xAB
	TODO,    //0xAC
	TODO,    //0xAD
	TODO,    //0xAE
	TODO,    //0xAF
	TODO,    //0xB0
	TODO,    //0xB1
	TODO,    //0xB2
	TODO,    //0xB3
	TODO,    //0xB4
	TODO,    //0xB5
	TODO,    //0xB6
	TODO,    //0xB7
	TODO,    //0xB8
	TODO,    //0xB9
	TODO,    //0xBA
	TODO,    //0xBB
	TODO,    //0xBC
	TODO,    //0xBD
	TODO,    //0xBE
	TODO,    //0xBF
	TODO,    //0xC0
	TODO,    //0xC1
	TODO,    //0xC2
	TODO,    //0xC3
	TODO,    //0xC4
	TODO,    //0xC5
	TODO,    //0xC6
	TODO,    //0xC7
	TODO,    //0xC8
	TODO,    //0xC9
	TODO,    //0xCA
	TODO,    //0xCB
	TODO,    //0xCC
	TODO,    //0xCD
	TODO,    //0xCE
	TODO,    //0xCF
	TODO,    //0xD0
	TODO,    //0xD1
	TODO,    //0xD2
	TODO,    //0xD3
	TODO,    //0xD4
	TODO,    //0xD5
	TODO,    //0xD6
	TODO,    //0xD7
	TODO,    //0xD8
	TODO,    //0xD9
	TODO,    //0xDA
	TODO,    //0xDB
	TODO,    //0xDC
	TODO,    //0xDD
	TODO,    //0xDE
	TODO,    //0xDF
	TODO,    //0xE0
	TODO,    //0xE1
	TODO,    //0xE2
	TODO,    //0xE3
	TODO,    //0xE4
	TODO,    //0xE5
	TODO,    //0xE6
	TODO,    //0xE7
	TODO,    //0xE8
	TODO,    //0xE9
	TODO,    //0xEA
	TODO,    //0xEB
	TODO,    //0xEC
	TODO,    //0xED
	TODO,    //0xEE
	TODO,    //0xEF
	TODO,    //0xF0
	TODO,    //0xF1
	TODO,    //0xF2
	TODO,    //0xF3
	TODO,    //0xF4
	TODO,    //0xF5
	TODO,    //0xF6
	TODO,    //0xF7
	TODO,    //0xF8
	TODO,    //0xF9
	TODO,    //0xFA
	TODO,    //0xFB
	TODO,    //0xFC
	TODO,    //0xFD
	TODO,    //0xFE
	TODO,    //0xFF
}

func TODO(cpu *cpu) cycleCount {
	// This function is not defined!
	return 0
}

func nop(cpu *cpu) cycleCount {
	// Does nothing
	return nopSize
}

// 3.3. Instructions
// The GameBoy CPU is based on a subset of the Z80 micro-
// processor. A summary of these commands is given below.
// If 'Flags affected' is not given for a command then
// none are affected.

// 3.3.1. 8-Bit Loads

// 3.3.1.1. LD nn,n
// Description:
// 		Put value nn into n.
// Use with:
// 		nn = B,C,D,E,H,L,BC,DE,HL,SP
// 		n = 8 bit immediate value

func ld_b_n(cpu *cpu) cycleCount {
	// Put value of register B into the parameter address
	// TODO: to implement
	return ld_b_nSize
}

func ld_c_n(cpu *cpu) cycleCount {
	// Put value of register C into the parameter address
	// TODO: to implement
	return ld_c_nSize
}

func ld_d_n(cpu *cpu) cycleCount {
	// Put value of register D into the parameter address
	// TODO: to implement
	return ld_d_nSize
}

func ld_e_n(cpu *cpu) cycleCount {
	// Put value of register E into the parameter address
	// TODO: to implement
	return ld_e_nSize
}

func ld_h_n(cpu *cpu) cycleCount {
	// Put value of register H into the parameter address
	// TODO: to implement
	return ld_h_nSize
}

func ld_l_n(cpu *cpu) cycleCount {
	// Put value of register L into the parameter address
	// TODO: to implement
	return ld_l_nSize
}

// 3.3.1.2. LD r1,r2
// Description:
// 		Put value r2 into r1.
// Use with:
// 		r1 = A,B,C,D,E,H,L,(HL)
//		r2 = A,B,C,D,E,H,L,(HL)

func ld_a_a(cpu *cpu) cycleCount {
	// Put value of register A into register A
	cpu.r.af.a = cpu.r.af.a
	return ld_a_aSize
}

func ld_a_b(cpu *cpu) cycleCount {
	// Put value of register B into register A
	cpu.r.af.a = cpu.r.bc.b
	return ld_a_bSize
}

func ld_a_c(cpu *cpu) cycleCount {
	// Put value of register C into register A
	cpu.r.af.a = cpu.r.bc.c
	return ld_a_cSize
}

func ld_a_d(cpu *cpu) cycleCount {
	// Put value of register D into register A
	cpu.r.af.a = cpu.r.de.d
	return ld_a_dSize
}

func ld_a_e(cpu *cpu) cycleCount {
	// Put value of register E into register A
	cpu.r.af.a = cpu.r.de.e
	return ld_a_eSize
}

func ld_a_h(cpu *cpu) cycleCount {
	// Put value of register H into register A
	cpu.r.af.a = cpu.r.hl.h
	return ld_a_hSize
}

func ld_a_l(cpu *cpu) cycleCount {
	// Put value of register L into register A
	cpu.r.af.a = cpu.r.hl.l
	return ld_a_lSize
}

func ld_a_hl(cpu *cpu) cycleCount {
	// Put value of the position of memory indicated by register HL into register A
	// TODO: how to implement this?
	//cpu.r.af.a = cpu.r.hl
	return ld_a_hlSize
}

func ld_b_b(cpu *cpu) cycleCount {
	// Put value of register B into register B
	cpu.r.bc.b = cpu.r.bc.b
	return ld_b_bSize
}

func ld_b_c(cpu *cpu) cycleCount {
	// Put value of register C into register B
	cpu.r.bc.b = cpu.r.bc.c
	return ld_b_cSize
}

func ld_b_d(cpu *cpu) cycleCount {
	// Put value of register D into register B
	cpu.r.bc.b = cpu.r.de.d
	return ld_b_dSize
}

func ld_b_e(cpu *cpu) cycleCount {
	// Put value of register E into register B
	cpu.r.bc.b = cpu.r.de.e
	return ld_b_eSize
}

func ld_b_h(cpu *cpu) cycleCount {
	// Put value of register H into register B
	cpu.r.bc.b = cpu.r.hl.h
	return ld_b_hSize
}

func ld_b_l(cpu *cpu) cycleCount {
	// Put value of register L into register B
	cpu.r.bc.b = cpu.r.hl.l
	return ld_b_lSize
}

func ld_b_hl(cpu *cpu) cycleCount {
	// Put value of the position of memory indicated by register HL into register B
	// TODO: how to implement this?
	return ld_b_hlSize
}

func ld_c_b(cpu *cpu) cycleCount {
	// Put value of register B into register C
	cpu.r.bc.c = cpu.r.bc.b
	return ld_c_bSize
}
func ld_c_c(cpu *cpu) cycleCount {
	// Put value of register C into register C
	cpu.r.bc.c = cpu.r.bc.c
	return ld_c_cSize
}
func ld_c_d(cpu *cpu) cycleCount {
	// Put value of register D into register C
	cpu.r.bc.c = cpu.r.de.d
	return ld_c_dSize
}
func ld_c_e(cpu *cpu) cycleCount {
	// Put value of register E into register C
	cpu.r.bc.c = cpu.r.de.e
	return ld_c_eSize
}
func ld_c_h(cpu *cpu) cycleCount {
	// Put value of register H into register C
	cpu.r.bc.c = cpu.r.hl.h
	return ld_c_hSize
}
func ld_c_l(cpu *cpu) cycleCount {
	// Put value of register L into register C
	cpu.r.bc.c = cpu.r.hl.l
	return ld_c_lSize
}
func ld_c_hl(cpu *cpu) cycleCount {
	// Put value of the position of memory indicated by register HL into register C
	// TODO: how to implement this?
	return ld_c_hlSize
}
func ld_d_b(cpu *cpu) cycleCount {
	// Put value of register B into register D
	cpu.r.de.d = cpu.r.bc.b
	return ld_d_bSize
}
func ld_d_c(cpu *cpu) cycleCount {
	// Put value of register C into register D
	cpu.r.de.d = cpu.r.bc.c
	return ld_d_cSize
}
func ld_d_d(cpu *cpu) cycleCount {
	// Put value of register D into register D
	cpu.r.de.d = cpu.r.de.d
	return ld_d_dSize
}
func ld_d_e(cpu *cpu) cycleCount {
	// Put value of register E into register D
	cpu.r.de.d = cpu.r.de.e
	return ld_d_eSize
}
func ld_d_h(cpu *cpu) cycleCount {
	// Put value of register H into register D
	cpu.r.de.d = cpu.r.hl.h
	return ld_d_hSize
}
func ld_d_l(cpu *cpu) cycleCount {
	// Put value of register L into register D
	cpu.r.de.d = cpu.r.hl.l
	return ld_d_lSize
}
func ld_d_hl(cpu *cpu) cycleCount {
	// Put value of the position of memory indicated by register HL into register D
	// TODO: how to implement this?
	return ld_d_hlSize
}
func ld_e_b(cpu *cpu) cycleCount {
	// Put value of register B into register E
	// TODO: to implement
	return ld_e_bSize
}
func ld_e_c(cpu *cpu) cycleCount {
	// Put value of register C into register E
	// TODO: to implement
	return ld_e_cSize
}
func ld_e_d(cpu *cpu) cycleCount {
	// Put value of register D into register E
	// TODO: to implement
	return ld_e_dSize
}
func ld_e_e(cpu *cpu) cycleCount {
	// Put value of register E into register E
	// TODO: to implement
	return ld_e_eSize
}
func ld_e_h(cpu *cpu) cycleCount {
	// Put value of register H into register E
	// TODO: to implement
	return ld_e_hSize
}
func ld_e_l(cpu *cpu) cycleCount {
	// Put value of register L into register E
	// TODO: to implement
	return ld_e_lSize
}
func ld_e_hl(cpu *cpu) cycleCount {
	// Put value of the position of memory indicated by register HL into register E
	// TODO: to implement
	return ld_e_hlSize
}
func ld_h_b(cpu *cpu) cycleCount {
	// Put value of register B into register H
	// TODO: to implement
	return ld_h_bSize
}
func ld_h_c(cpu *cpu) cycleCount {
	// Put value of register C into register H
	// TODO: to implement
	return ld_h_cSize
}
func ld_h_d(cpu *cpu) cycleCount {
	// Put value of register D into register H
	// TODO: to implement
	return ld_h_dSize
}
func ld_h_e(cpu *cpu) cycleCount {
	// Put value of register H into register H
	// TODO: to implement
	return ld_h_eSize
}
func ld_h_h(cpu *cpu) cycleCount {
	// Put value of register H into register H
	// TODO: to implement
	return ld_h_hSize
}
func ld_h_l(cpu *cpu) cycleCount {
	// Put value of register L into register H
	// TODO: to implement
	return ld_h_lSize
}
func ld_h_hl(cpu *cpu) cycleCount {
	// Put value of the position of memory indicated by register HL into register H
	// TODO: to implement
	return ld_h_hlSize
}
func ld_l_b(cpu *cpu) cycleCount {
	// Put value of register B into register L
	// TODO: to implement
	return ld_l_bSize
}
func ld_l_c(cpu *cpu) cycleCount {
	// Put value of register C into register L
	// TODO: to implement
	return ld_l_cSize
}
func ld_l_d(cpu *cpu) cycleCount {
	// Put value of register D into register L
	// TODO: to implement
	return ld_l_dSize
}
func ld_l_e(cpu *cpu) cycleCount {
	// Put value of register E into register L
	// TODO: to implement
	return ld_l_eSize
}
func ld_l_h(cpu *cpu) cycleCount {
	// Put value of register H into register L
	// TODO: to implement
	return ld_l_hSize
}
func ld_l_l(cpu *cpu) cycleCount {
	// Put value of register L into register L
	// TODO: to implement
	return ld_l_lSize
}
func ld_l_hl(cpu *cpu) cycleCount {
	// Put value of the position of memory indicated by register HL into register L
	// TODO: to implement
	return ld_l_hlSize
}
func ld_hl_b(cpu *cpu) cycleCount {
	// Put value of register B into the position of memory indicated by register HL
	// TODO: to implement
	return ld_hl_bSize
}
func ld_hl_c(cpu *cpu) cycleCount {
	// Put value of register C into the position of memory indicated by register HL
	// TODO: to implement
	return ld_hl_cSize
}
func ld_hl_d(cpu *cpu) cycleCount {
	// Put value of register D into the position of memory indicated by register HL
	// TODO: to implement
	return ld_hl_dSize
}
func ld_hl_e(cpu *cpu) cycleCount {
	// Put value of register E into the position of memory indicated by register HL
	// TODO: to implement
	return ld_hl_eSize
}
func ld_hl_l(cpu *cpu) cycleCount {
	// Put value of register L into the position of memory indicated by register HL
	// TODO: to implement
	return ld_hl_lSize
}
func ld_hl_h(cpu *cpu) cycleCount {
	// Put value of register H into the position of memory indicated by register HL
	// TODO: to implement
	return ld_hl_hSize
}

func ld_hl_n(cpu *cpu) cycleCount {
	// Put value of register ??? n ??? into the position of memory indicated by register HL
	// TODO: to implement
	// TODO: check what this OP does, as it
	return ld_hl_nSize
}
