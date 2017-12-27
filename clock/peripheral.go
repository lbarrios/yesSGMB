package clock

import "sync"

type Peripheral interface {
	ConnectClock(*sync.WaitGroup, Clock) chan uint64
	GetName() string
}
