package memo

import (
	"sync"
	"time"
)

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func NewRealClock() Clock {
	return RealClock{}
}

func (rc RealClock) Now() time.Time {
	return time.Now()
}

type FakeClock interface {
	Clock
	Advance(d time.Duration)
}

func NewFakeClock() FakeClock {
	return &fakeClock{
		now: time.Date(1984, time.April, 4, 0, 0, 0, 0, time.UTC),
	}
}

type fakeClock struct {
	sync.RWMutex
	now time.Time
}

func (fc *fakeClock) Now() time.Time {
	fc.RLock()
	defer fc.RUnlock()
	return fc.now
}

func (fc *fakeClock) Advance(d time.Duration) {
	fc.Lock()
	defer fc.Unlock()
	fc.now = fc.now.Add(d)
}
