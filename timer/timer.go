// Package timer provides the necessary structure to manage the
// synchronization between the other components
package timer

import (
	"github.com/lbarrios/yesSGMB/logger"
	"github.com/lbarrios/yesSGMB/types"
	"sync"
	"time"
)

const (
	DIV_ADDRESS  = types.Word(0xFF04)
	TIMA_ADDRESS = types.Word(0xFF05)
	TMA_ADDRESS  = types.Word(0xFF06)
	TAC_ADDRESS  = types.Word(0xFF07)
)

const (
	FREQ_4096   = 4096
	FREQ_16384  = 16384
	FREQ_65536  = 65536
	FREQ_262144 = 262144
)

type timer struct {
	log         logger.Logger
	div         *byte
	tima        *byte
	tma         *byte
	tac         *byte
}

func NewTimer(l *logger.Logger) *timer {
	t := new(timer)
	t.log = *l
	t.log.SetPrefix("\033[0;32mTIMER: ")
	return t
}

func (t *timer) Reset() {
	t.log.Println("Timer reset triggered.")
}

func (t *timer) MapByte(logical_address types.Address, physical_address *byte) {
	addr := logical_address.AsWord()
	switch {
	case addr == DIV_ADDRESS:
		t.div = physical_address
	case addr == TIMA_ADDRESS:
		t.tima = physical_address
	case addr == TMA_ADDRESS:
		t.tma = physical_address
	case addr == TAC_ADDRESS:
		t.tac = physical_address
	default:
		t.log.Fatalf("Trying to map unexpected address: 0x%.4x", addr)
	}
}

func (t *timer) step() {
	*t.tima -= *t.tma
	// t clock main = 0;
	if (*t.tima == 0) {
		*t.tima = *t.tma;
		// interrupt flag |= 4;
	}
}

func (t *timer) Run(wg *sync.WaitGroup) {
	t.log.Println("Timer started.")
	for {
		*t.div += 32
		t.log.Printf("tic, div=%d", *t.div)
		if *t.div == 0x00 {
			t.log.Println("div=0x00, stopping Timer.")
			break
		}
		time.Sleep(1000 * time.Millisecond)
	}
	wg.Done()
}
