// Package timer provides the necessary structure to manage the
// synchronization between the other components
package timer

type timer struct {
}

func NewTimer() *timer {
	t := new(timer)
	return t
}
