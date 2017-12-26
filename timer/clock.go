package timer

import (
	"github.com/lbarrios/yesSGMB/logger"
	"sync"
)

const (
	CLOCK_FREQ = 0x400000
)

type clock struct {
	log         logger.Logger
	t           uint64
	mutex       sync.Mutex
	peripherals map[string]chan uint64
	wg          sync.WaitGroup
}

type Clock interface {
	DisconnectPeripheral(peripheral Peripheral)
}

func NewClock(l *logger.Logger) *clock {
	c := new(clock)
	c.log = *l
	c.peripherals = make(map[string]chan uint64)
	c.log.SetPrefix("\033[0;35mCPU: ")
	return c
}

func (c *clock) step() {
	//c.log.Printf("c.t = %d", c.t)
	c.t += 4
}

func (c *clock) ConnectPeripheral(p Peripheral) {
	c.peripherals[p.GetName()] = p.ConnectClock(&c.wg, c)
}

// Lazy disconnect
func (c *clock) DisconnectPeripheral(p Peripheral) {
	c.mutex.Lock()
	delete(c.peripherals, p.GetName())
	c.mutex.Unlock()
}

func (c *clock) Run(wg *sync.WaitGroup) {
	c.log.Println("Clock started.")
	for {
		c.step()
		for _, p := range c.peripherals {
			c.wg.Add(1)
			p <- c.t
		}
		c.wg.Wait()
		if len(c.peripherals) == 0 {
			break
		}
	}
	wg.Done()
}
