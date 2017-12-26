package timer

import (
	"github.com/lbarrios/yesSGMB/logger"
	"sync"
)

const (
	CLOCK_FREQ = 0x400000
)

type Clock struct {
	Log         *logger.Logger
	t           uint64
	peripherals []chan uint64
	wg          sync.WaitGroup
}

func (c *Clock) step() {
	c.t++
}

func (c *Clock) ConnectPeripheral(p Peripheral) {
	c.peripherals = append(c.peripherals, p.ConnectClock(&c.wg))
}

func (c *Clock) Run(wg *sync.WaitGroup) {
	c.Log.Println("Clock started.")
	for {
		c.Log.Printf("c.t = %d", c.t)
		c.t += 1
		for _, p := range c.peripherals {
			c.wg.Add(1)
			p <- c.t
		}
		c.wg.Wait()
	}
}
