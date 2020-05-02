package memo_test

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dploop/memo"
	"github.com/dploop/memo/clock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Memo", func() {
	positive := func(v memo.Value, err error, target interface{}) {
		Expect(v).To(Equal(target))
		Expect(err).NotTo(HaveOccurred())
	}
	negative := func(v memo.Value, err error, target interface{}) {
		Expect(v).To(BeNil())
		Expect(err).To(MatchError(target))
	}

	Describe("Simple Memo", func() {
		It("can get nothing when empty", func() {
			m := memo.NewMemo()
			v, err := m.Get("k")
			negative(v, err, memo.ErrNotFound)
		})
		It("can get something when set", func() {
			m := memo.NewMemo()
			m.Set("k", "v")
			v, err := m.Get("k")
			positive(v, err, "v")
		})
	})

	vLoader := func(memo.Key) (memo.Value, error) {
		return "v", nil
	}
	uLoader := func(memo.Key) (memo.Value, error) {
		return "u", nil
	}
	Describe("Loader Memo", func() {
		It("can get something with specified loader", func() {
			m := memo.NewMemo()
			v, err := m.Get("k", memo.WithLoader(vLoader))
			positive(v, err, "v")
		})
		It("can get something with default loader", func() {
			m := memo.NewMemo(memo.WithLoader(vLoader))
			v, err := m.Get("k")
			positive(v, err, "v")
		})
		It("specified loader override default loader", func() {
			m := memo.NewMemo(memo.WithLoader(vLoader))
			v, err := m.Get("k", memo.WithLoader(uLoader))
			positive(v, err, "u")
		})
		It("loader will only take effect when needed", func() {
			m := memo.NewMemo()
			v, err := m.Get("k", memo.WithLoader(vLoader))
			positive(v, err, "v")
			v, err = m.Get("k", memo.WithLoader(uLoader))
			positive(v, err, "v")
		})
	})

	Describe("Expiration Memo", func() {
		var fakeClock *clock.FakeClock
		BeforeEach(func() {
			fakeClock = clock.NewFakeClock()
		})
		It("can get something when not expired", func() {
			m := memo.NewMemo(memo.WithClock(fakeClock))
			m.Set("k", "v", memo.WithExpiration(time.Minute))
			fakeClock.Advance(time.Minute - time.Second)
			v, err := m.Get("k")
			positive(v, err, "v")
		})
		It("expire something through specified expiration", func() {
			m := memo.NewMemo(memo.WithClock(fakeClock))
			m.Set("k", "v", memo.WithExpiration(time.Minute))
			fakeClock.Advance(time.Minute)
			v, err := m.Get("k")
			negative(v, err, memo.ErrNotFound)
		})
		It("expire something through default expiration", func() {
			m := memo.NewMemo(memo.WithClock(fakeClock), memo.WithExpiration(time.Minute))
			m.Set("k", "v")
			fakeClock.Advance(time.Minute)
			v, err := m.Get("k")
			negative(v, err, memo.ErrNotFound)
		})
		It("specified expiration will override default expiration", func() {
			m := memo.NewMemo(memo.WithClock(fakeClock), memo.WithExpiration(time.Hour))
			m.Set("k", "v", memo.WithExpiration(time.Minute))
			fakeClock.Advance(time.Minute)
			v, err := m.Get("k")
			negative(v, err, memo.ErrNotFound)
		})
	})

	Describe("Complicated Memo", func() {
		var fakeClock *clock.FakeClock
		BeforeEach(func() {
			fakeClock = clock.NewFakeClock()
		})
		It("can get nothing when loaded but expired", func() {
			m := memo.NewMemo(memo.WithClock(fakeClock))
			v, err := m.Get("k", memo.WithLoader(vLoader), memo.WithExpiration(time.Minute))
			positive(v, err, "v")
			fakeClock.Advance(time.Minute)
			v, err = m.Get("k")
			negative(v, err, memo.ErrNotFound)
		})
		It("access a key concurrently will only load once", func() {
			var counter int32
			sLoader := func(memo.Key) (memo.Value, error) {
				time.Sleep(time.Second)
				atomic.AddInt32(&counter, 1)
				return "s", nil
			}
			m := memo.NewMemo(memo.WithClock(fakeClock), memo.WithLoader(sLoader))
			start := time.Now()
			var wg sync.WaitGroup
			for i := 0; i < 10000; i++ {
				k := i % 100
				wg.Add(1)
				go func() {
					defer wg.Done()
					_, _ = m.Get(k)
				}()
			}
			wg.Wait()
			Expect(counter).To(BeEquivalentTo(100))
			Expect(time.Since(start)).To(
				BeNumerically("~", time.Second, time.Second/10),
			)
		})
	})
})

func BenchmarkMemo_Get(b *testing.B) {
	m := memo.NewMemo()
	m.Set("k", "v")
	for i := 0; i < b.N; i++ {
		_, _ = m.Get("k")
	}
}

func BenchmarkMemo_Set(b *testing.B) {
	m := memo.NewMemo()
	for i := 0; i < b.N; i++ {
		m.Set("k", "v")
	}
	_, _ = m.Get("k")
}

func BenchmarkMemo_GetConcurrent(b *testing.B) {
	m := memo.NewMemo()
	m.Set("k", "v")
	var wg sync.WaitGroup
	concurrency := runtime.NumCPU()
	each := b.N / concurrency
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < each; j++ {
				_, _ = m.Get("k")
			}
		}()
	}
	wg.Wait()
}

func BenchmarkMemo_SetConcurrent(b *testing.B) {
	m := memo.NewMemo()
	var wg sync.WaitGroup
	concurrency := runtime.NumCPU()
	each := b.N / concurrency
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < each; j++ {
				m.Set("k", "v")
			}
		}()
	}
	wg.Wait()
}
