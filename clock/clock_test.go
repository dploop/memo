package clock_test

import (
	"time"

	"github.com/dploop/memo/clock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Clock", func() {
	Describe("Real Clock", func() {
		It("should be close to the current time", func() {
			realClock := clock.NewRealClock()
			now := realClock.Now()
			current := time.Now()
			Expect(current.Sub(now)).To(
				BeNumerically("~", 0, time.Second),
			)
		})
	})
	Describe("Fake Clock", func() {
		It("should advance to the future time", func() {
			fakeClock := clock.NewFakeClock()
			now := fakeClock.Now()
			d := 7 * time.Minute
			fakeClock.Advance(d)
			future := fakeClock.Now()
			Expect(future.Sub(now)).To(Equal(d))
		})
	})
})
