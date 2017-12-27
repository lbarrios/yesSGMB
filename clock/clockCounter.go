package clock

import "sync"

type ClockCounter struct {
	Wg          *sync.WaitGroup
	Channel     chan uint64
	ClockCycles uint64
	Cycles      uint64
	Clock       Clock
}

func (c *ClockCounter) Init(wg *sync.WaitGroup, channel chan uint64, clock Clock) {
	c.Wg = wg
	c.Wg.Add(1)
	c.Channel = channel
	c.Clock = clock
}

func (c *ClockCounter) WaitNextCycle() {
	for {
		if c.Cycles <= c.ClockCycles {
			return
		}
		c.Wg.Done()
		c.ClockCycles = <-c.Channel
	}
}

func (c *ClockCounter) Disconnect(p Peripheral) {
	c.Clock.DisconnectPeripheral(p)
	c.Wg.Done()
}
