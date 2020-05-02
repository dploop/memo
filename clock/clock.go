package clock

import (
	"sync/atomic"
	"time"

	// For go:linkname mark.
	_ "unsafe"
)

type Clock interface {
	Now() int64
}

type RealClock struct{}

func NewRealClock() RealClock {
	return RealClock{}
}

func (rc RealClock) Now() int64 {
	return nanotime()
}

//go:linkname nanotime runtime.nanotime
func nanotime() int64

type FakeClock struct {
	nanotime int64
}

func NewFakeClock() *FakeClock {
	return &FakeClock{}
}

func (fc *FakeClock) Now() int64 {
	return atomic.LoadInt64(&fc.nanotime)
}

func (fc *FakeClock) Advance(d time.Duration) {
	atomic.AddInt64(&fc.nanotime, int64(d))
}
