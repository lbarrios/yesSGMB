package timer

import "sync"

type Peripheral interface {
	ConnectClock(*sync.WaitGroup) chan uint64
}
