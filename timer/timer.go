// Package timer provides the necessary structure to manage the
// synchronization between the other components
package timer

import (
	"github.com/lbarrios/yesSGMB/logger"
	"github.com/lbarrios/yesSGMB/types"
	"sync"
	"github.com/lbarrios/yesSGMB/clock"
)

const (
	DIV_ADDRESS  = types.Word(0xFF04)
	TIMA_ADDRESS = types.Word(0xFF05)
	TMA_ADDRESS  = types.Word(0xFF06)
	TAC_ADDRESS  = types.Word(0xFF07)
)

const (
	FREQ_4096     = 4096
	FREQ_16384    = 16384
	FREQ_65536    = 65536
	FREQ_262144   = 262144
	CYCLES_4096   = 1024
	CYCLES_16384  = 256
	CYCLES_65536  = 64
	CYCLES_262144 = 16
)

type timer struct {
	log   logger.Logger
	clock clock.ClockCounter
	div   *byte
	tima  *byte
	tma   *byte
	tac   *byte
}

func NewTimer(l *logger.Logger) *timer {
	t := new(timer)
	t.log = *l
	t.log.SetPrefix("\033[0;32mTIMER: ")
	return t
}

func (t *timer) ConnectClock(clockWg *sync.WaitGroup, clock clock.Clock) chan uint64 {
	t.clock.Init(clockWg, make(chan uint64), clock)
	return t.clock.Channel
}

func (t *timer) GetName() string {
	return "timer"
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
	if *t.tima == 0 {
		*t.tima = *t.tma
		// interrupt flag |= 4;
	}
}

func (t *timer) Run(wg *sync.WaitGroup) {
	t.log.Println("Timer started.")
	for {
		t.clock.WaitNextCycle()

		if *t.tac == 0xff {
			t.log.Println("tac=0xff, stopping timer.")
			t.clock.Disconnect(t)
			break
		}

		t.clock.Cycles += CYCLES_4096
		*t.div += 1
		if *t.div == 0x00 {
			t.log.Println("tic, div=0x00.")
		}
	}
	wg.Done()
}
