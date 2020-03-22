package memo_test

import (
	"testing"
	"time"

	"github.com/dploop/memo"
)

func TestRealClock(t *testing.T) {
	realClock := memo.NewRealClock()
	now := realClock.Now().UnixNano()
	if now <= 0 {
		t.Errorf("now(%v) is expected to be positive", now)
	}
}

func TestFakeClock(t *testing.T) {
	fakeClock := memo.NewFakeClock()
	u := fakeClock.Now()
	fakeClock.Advance(5 * time.Minute)
	v := fakeClock.Now()
	d := v.Sub(u)
	if d != 5*time.Minute {
		t.Errorf("d(%v) is expected to be 5m0s", d)
	}

}
