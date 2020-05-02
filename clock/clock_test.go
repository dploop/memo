package clock_test

import (
	"testing"
	"time"

	"github.com/dploop/memo/clock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Clock", func() {
	Describe("Real Clock", func() {
		It("should be a monotonic clock", func() {
			realClock := clock.NewRealClock()
			u := realClock.Now()
			v := realClock.Now()
			Expect(v - u).To(BeNumerically(">=", 0))
		})
	})
	Describe("Fake Clock", func() {
		It("should advance to the future", func() {
			fakeClock := clock.NewFakeClock()
			u := fakeClock.Now()
			d := 7 * time.Minute
			fakeClock.Advance(d)
			v := fakeClock.Now()
			Expect(v - u).To(BeNumerically("==", d))
		})
	})
})

func BenchmarkTimeClock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = time.Now()
	}
}

func BenchmarkRealClock(b *testing.B) {
	realClock := clock.NewRealClock()
	for i := 0; i < b.N; i++ {
		_ = realClock.Now()
	}
}

func BenchmarkFakeClock(b *testing.B) {
	fakeClock := clock.NewFakeClock()
	for i := 0; i < b.N; i++ {
		fakeClock.Advance(time.Nanosecond)
		_ = fakeClock.Now()
	}
}
