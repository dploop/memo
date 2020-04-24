package clock

import (
	"sync"
	"time"
)

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func NewRealClock() RealClock {
	return RealClock{}
}

func (rc RealClock) Now() time.Time {
	return time.Now()
}

type FakeClock struct {
	mutex sync.RWMutex
	clock time.Time
}

func NewFakeClock() *FakeClock {
	return &FakeClock{clock: time.Date(1984, time.April, 4, 0, 0, 0, 0, time.UTC)}
}

func (fc *FakeClock) Now() time.Time {
	fc.mutex.RLock()
	defer fc.mutex.RUnlock()
	return fc.clock
}

func (fc *FakeClock) Advance(d time.Duration) {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()
	fc.clock = fc.clock.Add(d)
}
